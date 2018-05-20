
功能：
服务端提供一个http服务，客户端通过定时向服务端发送信息来确认并记录客户端的IP地址来达到动态DNS的目的。
此工具通过调用dnspod API来达到实现DDNS功能，并利用telegram bot来达到提供IP变动时的提醒功能。

此工具如何鉴定客户端发送信息的从正确的客户端发出的？
1.服务端和客户端提前设定好的一段字符A.
2.时间戳数值整除一个数值B得到结果C.这样在B秒内,服务端和客户端得到的结果C是一致的.
3.生成token,token=sha256(A+C).token每隔B秒会发生变动.
4.通过鉴定token来确定信息是否来自可信客户端.
5.使用https通讯防止中间人对信息的截获及篡改.

Install:`go install homeip`

`homeip`为服务端可执行程序
`client/myip.sh`为客户端脚本

Usage:`homeip -h`
```
Usage of homeip:
  -c string
        json config file (default "config.json")
  -v    version
```

config.json:
```
{
  "addr": "127.0.0.1:9090",  // http server addres
  "step": 20,  // how much time change the token
  "random": "change it",  // your key
  // For dnspod API https://www.dnspod.cn/docs/records.html#dns
  "dnspodtoken": "change it",
  "format": "json",
  "domain_id": "change it",
  "record_id": "change it",
  "sub_domain": "change it",
  "record_line": "默认",
  // For telegram bot
  "telegram_token": "change it",
  "notification_id":"change it"
}
```

NOTE:
you can `GET /myip` to view your IP with no Control access on this server.if you don't want any body know you home ip,your can run this server under nginx to [Restricting Access with HTTP Basic Authentication](https://docs.nginx.com/nginx/admin-guide/security-controls/configuring-http-basic-authentication/)