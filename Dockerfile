#源镜像
FROM golang:latest

LABEL author="Trial <10223062kong_liangliang@cn.tre-inc.com>" describe="BEANQ monitor ui"
LABEL describe="test image"

## 在docker的根目录下创建相应的使用目录
RUN mkdir -p /www/webapp
## 设置工作目录
WORKDIR /www/webapp
## 把当前（宿主机上）目录下的文件都复制到docker上刚创建的目录下
COPY . /www/webapp

#go构建可执行文件
RUN go build main.go
#暴露端口
EXPOSE 8080

RUN chmod +x main
ENTRYPOINT ["./main"]