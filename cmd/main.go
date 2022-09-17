// @author cold bin
// @date 2022/9/15

package main

import (
	"github.com/jordan-wright/email"
	"jwzx-mail/conf"
	"jwzx-mail/logic"
	"jwzx-mail/mail"
	"jwzx-mail/util"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	if err := do(); err != nil {
		if err == logic.ErrNotTheLatest {
			log.Println("当前无新通知")
			return
		}
		log.Println("err: ", err)
		return
	}
	log.Println("爬取成功")
}

func do() error {
	news, err := logic.GetNewNews()
	if err != nil {
		return err
	}
	log.Println("一共有", len(news), "条新消息")
	if len(news) == 0 {
		return nil
	}
	// 然后把邮件发出去
	client := mail.NewClient()
	defer client.EPool.Close()

	for _, v := range news {
		log.Println("爬取消息 -> ", v.Title)

		e := email.NewEmail()
		if err = mail.PutHeader(conf.AConf.SendMailFrom, conf.AConf.SendMailToQqs, e); err != nil {
			return err
		}
		mail.PutTitle(v.Title, e)
		mail.PutHtml(util.QuickS2B(v.Content), e)
		mail.PutTos(conf.AConf.SendMailToQqs, e)
		if err = mail.PutAttachments(v.Files, e); err != nil {
			return err
		}

		if err = client.SendHtmlWithAttachments(e); err != nil {
			return err
		}
	}

	return nil
}
