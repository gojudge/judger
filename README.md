# judger 

| windows | linux-travis |
|---------|--------------|
| [![appveyor](https://ci.appveyor.com/api/projects/status/s20r7g9jfgxngiik?svg=true)](https://ci.appveyor.com/project/duguying/judger) | [![travis](https://api.travis-ci.org/duguying/judger.png)](https://travis-ci.org/duguying/judger) |

the judger server of online judge system

## Build

```shell
go get
go build
```

## Install from Docker

```shell
docker pull duguying/judger
```

```shell
mkdir /var/goj/judger
docker run -d -p 1004:1004 -p 1005:1005 -v /var/goj/judger:/data duguying/judger
```

## License #

MIT License
