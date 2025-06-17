2025年春重庆大学现代软件架构小组作业
实现了一个简易的网页计算器，前端使用`HTML + CSS + Javascript`，后端使用`Golang`，前后端间通过`websocket`通信。

一、源码目录结构

go-calculator
├─ backend                  后端
│    ├─ go.mod              
│    ├─ go.sum              
│    ├─ main.go             主程序入口
│    └─ packages            计算器模块
│           ├─ calc         
│           ├─ launch       
│           └─ stack        
├─ frontend                 前端
│    ├─ .eslintrc.json      
│    ├─ index.html          
│    ├─ script.js           
│    └─ style.css          
└─ run.bat                 