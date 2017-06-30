FROM google/golang
MAINTAINER duguying2008@gmail.com

RUN apt-get update -y
RUN apt-get upgrade -y
RUN go env

ADD . /go/src/github.com/gojudge/judger

# set the working directory and add current stuff
WORKDIR /go/src/github.com/gojudge/judger
RUN export GOBIN=$GOPATH/bin
RUN go get
RUN go build

EXPOSE 1004 1005
ENTRYPOINT []
CMD ["./judger","-c=/data/config_docker.ini","-mode=docker"]
