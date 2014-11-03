:: set devcpp binary path as global environment path

@set PATH=%PATH%;C:\Dev-Cpp\bin

cd %2
g++ %1 1> BUILD.LOG 2>&1
echo %ERRORLEVEL% > BUILDRESULT