::沙箱运行调用脚本
::%4 is the directory of executable
::%1 is executable relative
::%2 is -t=time
::%3 is -m=mem

::@ echo off 
::chcp 65001
cd %4
"D:\GOPATH\src\github.com\duguying\judger\sandbox\c\build\executer.exe" %1 %2 %3 > RUN.LOG
