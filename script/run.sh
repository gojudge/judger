# 沙箱运行调用脚本
# %1 is executable relative
# %2 is -t=time
# %3 is -m=mem

/root/gopath/src/github.com/duguying/judger/sandbox/c/build/executer %1 %2 %3 -c=/root/gopath/src/github.com/duguying/judger/sandbox/c/build/executer.json
