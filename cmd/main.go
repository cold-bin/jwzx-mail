// @author cold bin
// @date 2022/9/15

package main

import (
	"github.com/jordan-wright/email"
	"jwzx-mail/conf"
	"jwzx-mail/logic"
	"jwzx-mail/mail"
	"jwzx-mail/model"
	"jwzx-mail/util"
	"log"
	"sync"
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

	chans := make([]chan struct{}, len(news))

	for i := 0; i < len(chans); i++ {
		chans[i] = make(chan struct{})
	}

	defer func() {
		for _, c := range chans {
			close(c)
		}
	}()

	var wg sync.WaitGroup
	for i, v := range news {
		//TODO: 原版非并发，顺序执行
		//log.Println("爬取消息 -> ", v.Title)
		//e := email.NewEmail()
		//if err = mail.PutHeader(conf.AConf.SendMailFrom, conf.AConf.SendMailToQqs, e); err != nil {
		//	 return err
		//}
		//mail.PutTitle(v.Title, e)
		//mail.PutHtml(util.QuickS2B(v.Content), e)
		//mail.PutTos(conf.AConf.SendMailToQqs, e)
		//if err = mail.PutAttachments(v.Files, e); err != nil {
		//	 return err
		//}
		//if err = client.SendHtmlWithAttachments(e); err != nil {
		//	 return err
		//}
		wg.Add(1)

		if i == 0 {
			go func(v model.NewsContent, client mail.Client, i int) {
				err := constructAndSend(v, client)
				if err != nil {
					log.Println("有一封邮件发送出错：", err)
				}
				wg.Done()
				chans[i] <- struct{}{}

			}(v, *client, i)
		} else if i == len(news)-1 {
			go func(v model.NewsContent, client mail.Client, i int) {
				<-chans[i-1]
				err := constructAndSend(v, client)
				if err != nil {
					log.Println("有一封邮件发送出错：", err)
				}
				wg.Done()
				chans[i] <- struct{}{}

			}(v, *client, i)
		} else {
			go func(v model.NewsContent, client mail.Client, i int) {
				<-chans[i-1]
				err := constructAndSend(v, client)
				if err != nil {
					log.Println("有一封邮件发送出错：", err)
				}
				wg.Done()
				chans[i] <- struct{}{}
			}(v, *client, i)
		}
	}

	wg.Wait()
	return nil
}

func constructAndSend(v model.NewsContent, client mail.Client) (err error) {
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
	err = client.SendHtmlWithAttachments(e)

	return err
}
