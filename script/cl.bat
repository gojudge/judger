:: CL编译器调用
:: 使用方法：
::     cl.bat <filename>
:: 如果编译失败的话ERRORLEVEL文件中的值大于0，成功则等于0

::@ echo off 
chcp 65001
call "C:\Program Files (x86)\Microsoft Visual Studio 10.0\VC\vcvarsall.bat" x86
cd %2
cl %1 > BUILD.LOG

echo %ERRORLEVEL% > BUILDRESULT