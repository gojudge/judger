# judger [![Build Status](https://drone.io/github.com/duguying/judger/status.png)](https://drone.io/github.com/duguying/judger/latest) #

the judger server of online judge system

# Install from Docker

```shell
docker pull duguying/judger
```

```shell
mkdir /var/goj/judger
docker run -d -p 1004:1004 -p 1005:1005 -v /var/goj/judger:/data duguying/judger
```

# License #

MIT License
