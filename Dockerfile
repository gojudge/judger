FROM google/golang
MAINTAINER duguying2008@gmail.com

RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get install -y gcc
RUN apt-get install -y g++
RUN apt-get install -y cmake

ADD . /gopath/src/github.com/duguying/judger

# set the working directory and add current stuff
WORKDIR /gopath/src/github.com/duguying/judger/sandbox/c/build
RUN cmake ..
RUN make

WORKDIR /gopath/src/github.com/duguying/judger
RUN git checkout master
RUN go get
RUN go build
RUN mkdir build

EXPOSE 1004 1005
ENTRYPOINT []
CMD ["./judger","-c=/data/config_docker.ini","-mode=docker"]