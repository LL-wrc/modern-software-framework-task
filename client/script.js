const display = document.getElementById('display');
let currentExpression = '';

function appendCharacter(char) {
    if (display.innerText === '0' && char !== '.') {
        display.innerText = '';
    }
    if (char === '%' && currentExpression === '') return; 
    const lastChar = currentExpression.slice(-1);
    const operators = ['+', '-', '*', '/', '%'];   //todo 先示 % 为求余，后续会改为百分号，以符合计算器规则
    if (operators.includes(lastChar) && operators.includes(char)) {
        if (char === '-' && (currentExpression.length === 0 || operators.includes(currentExpression.slice(-2, -1)))) {
            // 处理+/-的情况
            currentExpression += char;
            display.innerText += char;
        } else if (operators.includes(lastChar) && char !== '-'){
             // 多个连续运算符取最后一个输入
            currentExpression = currentExpression.slice(0, -1) + char;
            display.innerText = display.innerText.slice(0, -1) + char;
        } else if (!operators.includes(lastChar)){
            currentExpression += char;
            display.innerText += char;
        }
    } else {
        currentExpression += char;
        display.innerText += char;
    }
}

function clearDisplay() {
    currentExpression = '';
    display.innerText = '0';
}

function deleteLast() {
    currentExpression = currentExpression.slice(0, -1);
    display.innerText = display.innerText.slice(0, -1);
    if (display.innerText === '') {
        display.innerText = '0';
    }
}

async function calculateResult() {
    if (currentExpression === '' || ['+', '-', '*', '/', '%', '.'].includes(currentExpression.slice(-1))) {
        display.innerText = 'Error';
        currentExpression = '';
        return;
    }
    try {
        const response = await fetch('/api/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ expression: currentExpression }),
        });
        const data = await response.json();
        if (data.error) {
            display.innerText = `Error: ${data.error}`;
            currentExpression = '';
        } else {
            display.innerText = data.result;
            currentExpression = String(data.result); // 结果可以参加下一次运算
        }
    } catch (error) {
        display.innerText = 'Error';
        currentExpression = '';
        console.error('Calculation error:', error);
    }
}

async function calculateUnaryOperation(operation) {
    if (currentExpression === '' || ['+', '-', '*', '/', '%', '.'].includes(currentExpression.slice(-1))) {
        display.innerText = 'Error';
        currentExpression = '';
        return;
    }
    let expressionToCalculate = '';
    switch (operation) {
        case 'sq':
            // 修改平方的计算方式，使用 (expression)*(expression) 替代 (expression)^2
            // 确保后端能够正确处理乘法
            expressionToCalculate = `(${currentExpression})*(${currentExpression})`;
            break;
        case 'sqrt':
            expressionToCalculate = `sqrt(${currentExpression})`;
            break;
        case 'inv':
            expressionToCalculate = `1/(${currentExpression})`;
            break;
        default:
            display.innerText = 'Error';
            currentExpression = '';
            return;
    }

    try {
        const response = await fetch('/api/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ expression: expressionToCalculate }),
        });
        const data = await response.json();
        if (data.error) {
            display.innerText = `Error: ${data.error}`;
            currentExpression = '';
        } else {
            display.innerText = data.result;
            currentExpression = String(data.result);
        }
    } catch (error) {
        display.innerText = 'Error';
        currentExpression = '';
        console.error('Unary operation error:', error);
    }
}

function clearEntry() {
    // CE (Clear Entry) - 清除运算符或运算符后整串数字，暂时和计算器CE功能有差异 
    if (display.innerText !== '0' && display.innerText !== 'Error') {
        const parts = currentExpression.match(/(\d+\.?\d*|[^\d.]+)/g) || [];
        if (parts.length > 0) {
            parts.pop(); 
            currentExpression = parts.join('');
            display.innerText = currentExpression || '0';
        } else {
            clearDisplay();
        }
        if (currentExpression === '') {
             display.innerText = '0';
        }
    } else {
        clearDisplay();
    }
}


function toggleSign() {
    if (currentExpression === '' || currentExpression === '0' || display.innerText === 'Error') return;

    const match = currentExpression.match(/([\+\-\*\/\%]?\d+\.?\d*)$/);
    if (match) {
        let lastNumberStr = match[0];
        let prefix = currentExpression.substring(0, currentExpression.length - lastNumberStr.length);
        
        if (lastNumberStr.startsWith('-')) {
            lastNumberStr = lastNumberStr.substring(1);
        } else if (lastNumberStr.startsWith('+')) {
            lastNumberStr = '-' + lastNumberStr.substring(1);
        } else if (lastNumberStr.startsWith('*') || lastNumberStr.startsWith('/') || lastNumberStr.startsWith('%')) {
            lastNumberStr = lastNumberStr[0] + '(-' + lastNumberStr.substring(1) + ')';
        } else {
            lastNumberStr = '-' + lastNumberStr;
        }
        currentExpression = prefix + lastNumberStr;
        display.innerText = currentExpression;
    } else if (!isNaN(parseFloat(currentExpression))) {
        if (currentExpression.startsWith('-')) {
            currentExpression = currentExpression.substring(1);
        } else {
            currentExpression = '-' + currentExpression;
        }
        display.innerText = currentExpression;
    }
}
