FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app

MAINTAINER wenzhenglin(http://g.haodai.net/wenzhenglin/wechat-receiver.git)

COPY wechat-receiver /app

CMD /app/wechat-receiver
ENTRYPOINT ["./wechat-receiver"]

# EXPOSE 8080