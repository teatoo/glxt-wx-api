@echo off
:RESTART
tasklist | find /C "glxt-wx-api_v2.exe" > temp.txt
set /p num=<temp.txt
del /F temp.txt
echo It has %num% program runing! 
if "%num%" == "0" start "�̹̽���ϵͳ΢�ŷ������" /d "." glxt-wx-api_v2.exe 
ping -n 10 127.0.0.1 > temp.txt
del /F temp.txt
goto RESTART
