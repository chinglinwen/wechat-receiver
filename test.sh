# -w http://wechat-notify.devops.haodai.net
#   -r wenzhenglin \
go build
./wechat-receiver \
  -w http://localhost:8001 \
  -agentid 1000005 \
  -secret 3Kds9ib-5JwY7-DrlxGIBq7XOjYDf846W3_Tda2sLe0 \
  -party 10