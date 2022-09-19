// @author cold bin
// @date 2022/9/15

package conf

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

var AConf AppConf

type AppConf struct {
	HttpTimeOut         time.Duration `json:"http_time_out,omitempty"`  // http客户端请求超时时间，为零时没有超时
	LatestFileId        int           `json:"latest_file_id,omitempty"` // 消息栏最新fileId
	SendMailPoolGet     int           `json:"send_mail_pool_get"`
	SendMailToQqs       []string      `json:"send_mail_to_qqs,omitempty"`       // 邮件将要发送的邮箱列表
	QQMailAuthCode      string        `json:"qq_mail_auth_code,omitempty"`      // qq邮箱的授权码
	SendMailServerQQ    string        `json:"send_mail_server_qq,omitempty"`    // 当前qq邮箱作为服务端
	SendMailServerHost  string        `json:"send_mail_server_host"`            // 主机地址
	SendMailServerAddr  string        `json:"send_mail_server_addr,omitempty"`  // 发送邮件的服务地址
	SendMailFrom        string        `json:"send_mail_from,omitempty"`         // from
	SendMailTitlePrefix string        `json:"send_mail_title_prefix,omitempty"` // 邮件标题前缀
}

func (a *AppConf) IsLatestFile(fileId int) bool {
	return a.LatestFileId < fileId
}

var rwMutex sync.RWMutex

func (a *AppConf) UpdateFileId(fileId int) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	// 持久化到配置文件 /home/coldbin/jwzx-mail/conf/conf.json
	f1, err := os.OpenFile("conf/conf.json", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f1.Close()

	AConf.LatestFileId = fileId

	bytes, err := json.Marshal(AConf)
	if err != nil {
		return err
	}

	if _, err = f1.Write(bytes); err != nil {
		return err
	}

	return nil
}

// 读取配置文件
func init() {
	f, err := os.OpenFile("conf/conf.json", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
		panic(err)
		return
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&AConf); err != nil {
		log.Println("json 文件解析失败: ", err)
		panic(err)
		return
	}

	log.Println("qq服务端qq的授权码: ", AConf.QQMailAuthCode)
}
