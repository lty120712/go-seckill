@echo off
:: 切换到项目根目录
cd /d %~dp0\..

:: 启动 fresh，监视文件变动并自动重启服务
fresh
