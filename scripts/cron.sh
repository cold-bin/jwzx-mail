# 这里是定时调用 work.sh !注意：晚上jwzx服务器关闭

# 每天 8:00-22:30 之间，每隔30分钟爬取jwzx
*/30 8-22 * * * . /etc/profile;/bin/sh /home/cold-bin/jwzx-mail/scriptcs/work.sh
