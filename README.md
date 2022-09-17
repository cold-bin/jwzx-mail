## `jwzx-mail`
简易的爬虫工具，并将爬取的数据主动推送到QQ邮箱
### 设计思路
本工具分为两个部分：
- mail service: 邮件服务。该服务主动推送教务在线最新消息的邮件，支持正文和附件下载。主动推送，采用轮询方式刷新数据，发现拿到了新数据，就推送。（轮询采用linux定时任务脚本实现，定时请求爬虫api）
- jwzx scrapy service: 爬虫服务。该服务需要爬取教务在线的公告消息栏，并且最好记住一个爬取最新消息的一个时间及文件id。这样下次轮询爬取的时候，
如果发现教务在线第一条最新消息时间及id没有发生变化，那么就不需要触发邮件发送