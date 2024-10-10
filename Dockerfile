# source image
FROM golang:latest AS builder

LABEL author="Trial <10223062kong_liangliang@cn.tre-inc.com>" describe="BeanQ Monitoring UI"
LABEL describe="BeanQ UI"

## create work folder for docker
RUN mkdir -p /www/webapp
## set work folder
WORKDIR /www/webapp

COPY . /www/webapp

RUN CGO_ENABLED=0 GOOS=linux go build -o beanqui .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /www/webapp/env.json .
COPY --from=builder /www/webapp/beanqui .

EXPOSE 9090

RUN chmod +x beanqui
ENTRYPOINT ["./beanqui","-port",":9090"]