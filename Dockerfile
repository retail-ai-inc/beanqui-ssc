# source image
FROM golang:latest

LABEL author="Trial <10223062kong_liangliang@cn.tre-inc.com>" describe="BEANQ Monitor UI"
LABEL describe="test image"

## create work folder for docker
RUN mkdir -p /www/webapp
## set work folder
WORKDIR /www/webapp

COPY . /www/webapp

RUN go build main.go

EXPOSE 9090

RUN chmod +x main
ENTRYPOINT ["./main","-port",":9090"]