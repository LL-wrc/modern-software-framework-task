# 2025年春重庆大学现代软件架构小组作业

这是一个简单的计算器应用，包含以下三个主要部分：

1.  **客户端 (Client)**: 使用 HTML, CSS, 和 JavaScript 实现的用户界面,类似计算器界面，用户可以在此输入数学和运算符，并查看计算结果。
2.  **服务端/中间层 (Server/Middle Layer)**: 使用 Python (Flask) 实现的中间层。它负责接收来自客户端的请求，并将计算任务转发给后端服务，然后将后端返回的结果再传回给客户端。同时，它也负责提供客户端的静态文件。
3.  **后端 (Backend)**: 使用 Go 实现的计算引擎。它接收服务端转发过来的数学表达式，执行实际的计算，并将结果返回给服务端。

## 目录结构

```
d:\cal_3/
├── backend/            # Go 后端服务
│   └── main.go
├── client/             # JavaScript 客户端
│   ├── background.jpg
│   ├── index.html
│   ├── script.js
│   └── style.css
├── Middle Layer/             # Python 服务端
│   └── Middle Layer.py
├── API_Doc.md          # 接口文档
└── README.md           # 项目说明
```

## 运行指南

1.  **启动服务端 (Python)**:
    ```bash
    cd Middle Layer
    # 确保已安装 Flask (pip install Flask requests)
    python Middle Layer.py
    # 服务端将在 http://localhost:8080 启动
    ```
    
2.  **启动后端服务 (Go)**:
    ```bash
    cd backend
    go run main.go
    # 后端服务将在 http://localhost:8000 启动
    ```
    
3.  **访问客户端**:
    在浏览器中打开 `http://localhost:8080/` 即可访问计算器界面。

## 注意事项

*   确保各服务端口未被占用，或根据实际情况修改各服务中的端口配置。
*   后端 Go 服务目前仅支持基础的四则运算和一元运算，运算优先级为 +/- > 乘除 > 加减 > 一元运算
