# -w http://wechat-notify.devops.haodai.net
#   -r wenzhenglin \
echo compiling....
go build
./wechat-receiver \
  -w http://wechat-notify.haodai.net \
  -agentid 1000005 \
  -secret 3Kds9ib-5JwY7-DrlxGIBq7XOjYDf846W3_Tda2sLe0 \
  -agentidDev 1000006 \
  -secretDev zfLrCAxVHJxXhKk636krUtRwg0wxDO3EAUpihGw19AA \
  -party 2 \
  -partyDev 10 \
  -releaseurl http://localhost:8089/api/wechat
  # -r wenzhenglin
  # -party 10 