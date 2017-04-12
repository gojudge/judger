# judger 

| windows | linux-travis |
|---------|--------------|
| [![Build status](https://ci.appveyor.com/api/projects/status/4n48mdcqoss6nhsm?svg=true)](https://ci.appveyor.com/project/duguying/judger-5bddq) | [![Build Status](https://travis-ci.org/gojudge/judger.svg?branch=master)](https://travis-ci.org/gojudge/judger) |

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
