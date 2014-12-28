::沙箱运行调用脚本
::%1 is executable relative
::%2 is -t=time
::%3 is -m=mem

::@ echo off 
::chcp 65001
"D:\GOPATH\src\github.com\duguying\judger\sandbox\c\build\executer.exe" %1 %2 %3
