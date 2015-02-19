# judger 

| linux-drone | windows | linux-travis |
|-------------|---------|--------------|
| [![drone](https://drone.io/github.com/duguying/judger/status.png)](https://drone.io/github.com/duguying/judger/latest) | [![appveyor](https://ci.appveyor.com/api/projects/status/s20r7g9jfgxngiik?svg=true)](https://ci.appveyor.com/project/duguying/judger) | [![travis](https://api.travis-ci.org/duguying/judger.png)](https://travis-ci.org/duguying/judger) |

the judger server of online judge system

## Build in Linux

```shell
go get
go build
cd sandbox/c/build
cmake ..
make
```

## Build in Windows

```shell
go get
go build
cd sandbox/c/build
cmake -G"NMake Makefiles" ..
nmake
```

## Install from Docker

```shell
docker pull duguying/judger
```

```shell
mkdir /var/goj/judger
docker run -d -p 1004:1004 -p 1005:1005 -v /var/goj/judger:/data duguying/judger
```

## Executer

The executers written in C. Linux Version is a simple sandbox which could intercept dangerous syscalls, the Windows Version does not support syscall interception. So, the Linux Version judger is suggested. If you need a highly security judger, I suggest you deploy it with docker.

## Net

The judger support two kinds of network transmission protocol, TCP and HTTP, the data format are both json-based. Default port TCP:1004 and HTTP:1005.

## License #

MIT License
