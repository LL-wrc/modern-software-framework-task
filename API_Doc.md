# 计算器项目三端接口文档

本文档详细定义了计算器项目中 **前端 (Client)**、**Python中间层 (Middle Layer)** 和 **Go后端 (Backend)** 之间的接口交互。

---

## 1. 架构概览

本项目采用经典的三层架构：

1.  **前端 (Client)**: HTML/CSS/JS 实现的用户界面，运行在浏览器中。
2.  **Python中间层 (Middle Layer)**: 基于 Flask 的 Web 服务器，负责托管前端静态文件，并作为API网关将计算请求代理到后端。
3.  **Go后端 (Backend)**: 计算服务，负责执行实际的数学表达式计算。

---

## 2. 前端 (Client) <-> Python中间层 (Middle Layer) 接口

这是用户前端与Python服务器之间的接口。

### 2.1. 静态文件服务

*   **接口说明**: 提供计算器应用所需的HTML页面及相关的CSS和JS文件。
*   **URI**:
    *   `http://localhost:8080/`: 获取主页面 `index.html`
    *   `http://localhost:8080/style.css`: 获取CSS样式文件
    *   `http://localhost:8080/script.js`: 获取JavaScript逻辑文件
*   **HTTP Method**: `GET`
*   **Request**: 无特定参数
*   **Response**: 对应文件的内容 (`text/html`, `text/css`, `application/javascript`)

### 2.2. 计算API

*   **接口说明**: 前端将用户输入的表达式发送到此接口进行计算。所有的一元运算（如开方、百分比）在前端被格式化为标准表达式字符串后发送。
*   **URI**: `http://localhost:8080/api/calculate`
*   **HTTP Method**: `POST`
*   **Request**:
    *   **Header**: `Content-Type: application/json`
    *   **Body (JSON)**:
        | 字段名     | 类型   | 是否必须 | 描述与示例                                                                                             |
        |------------|--------|----------|--------------------------------------------------------------------------------------------------------|
        | expression | string | Required | 用户输入的、或由JS格式化后的数学表达式。<br>**示例**: `"5*2"`, `"sqrt(9)"`, `"(10+5)%"` |
*   **Response**:
    *   **成功 (JSON)**:
        | 字段名 | 类型   | 描述     |
        |--------|--------|----------|
        | result | number | 计算结果 |
    *   **失败 (JSON)**:
        | 字段名 | 类型   | 描述         |
        |--------|--------|--------------|
        | error  | string | 错误信息描述 |

---

## 3. Python中间层 (Middle Layer) <-> Go后端 (Backend) 接口

这是Python服务器与Go计算服务之间的内部接口。

### 3.1. 内部计算API

*   **接口说明**: Python中间层将从前端接收到的计算请求原封不动地转发给Go后端服务。
*   **URI**: `http://localhost:8000/calculate` (此地址在`Middle Layer.py`中配置)
*   **HTTP Method**: `POST`
*   **Request**:
    *   **Header**: `Content-Type: application/json`
    *   **Body (JSON)**:
        | 字段名     | 类型   | 是否必须 | 描述                 |
        |------------|--------|----------|----------------------|
        | expression | string | Required | 从前端转发的数学表达式 |
*   **Response**:
    *   **成功 (JSON)**:
        | 字段名 | 类型    | 描述     |
        |--------|---------|----------|
        | result | float64 | 计算结果 |
    *   **失败 (JSON)**:
        | 字段名 | 类型   | 描述                                         |
        |--------|--------|----------------------------------------------|
        | error  | string | 错误信息 (例如: `"invalid expression: ..."`) |
