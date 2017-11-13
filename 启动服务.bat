@echo off
:RESTART
tasklist | find /C "glxt-wx-api_v2.exe" > temp.txt
set /p num=<temp.txt
del /F temp.txt
echo It has %num% program runing! 
if "%num%" == "0" start "继教管理系统微信服务程序" /d "." glxt-wx-api_v2.exe 
ping -n 10 127.0.0.1 > temp.txt
del /F temp.txt
goto RESTART
