// @author cold bin
// @date 2022/9/15

package main

import (
	"fmt"
	"github.com/jordan-wright/email"
	"jwzx-mail/conf"
	"jwzx-mail/logic"
	"jwzx-mail/mail"
	"jwzx-mail/model"
	"jwzx-mail/util"
	"log"
	"runtime"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	news, err := logic.GetNewNews()
	if err != nil {
		if err == logic.ErrNotTheLatest {
			log.Println("当前无新通知")
			return
		}
		log.Println("err: ", err)
		return
	}

	log.Println("一共有", len(news), "条新消息")
	if len(news) == 0 {
		return
	}

	// 然后把邮件发出去
	client := mail.NewClient()
	defer client.EPool.Close()

	chans := make([]chan struct{}, len(news))

	for i := 0; i < len(chans); i++ {
		chans[i] = make(chan struct{})
	}

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
		//wg.Add(1)
		if i == 0 {
			go func(v model.NewsContent, client mail.Client, i int) {
				fmt.Println("i=", i)
				if err := constructAndSend(v, client); err != nil {
					log.Println("有一封邮件发送出错：", err)
				}

				//wg.Done()
				chans[i] <- struct{}{}
				close(chans[i])
			}(v, *client, i)
		} else if i == len(news)-1 {
			go func(v model.NewsContent, client mail.Client, i int) {
				fmt.Println("i=", i)
				<-chans[i-1]
				if err := constructAndSend(v, client); err != nil {
					log.Println("有一封邮件发送出错：", err)
				}

				//wg.Done()
				chans[i] <- struct{}{}
				close(chans[i])
			}(v, *client, i)
		} else {
			go func(v model.NewsContent, client mail.Client, i int) {
				log.Println("i=", i)
				<-chans[i-1]

				if err := constructAndSend(v, client); err != nil {
					log.Println("有一封邮件发送出错：", err)
				}

				//wg.Done()
				chans[i] <- struct{}{}
				close(chans[i])

			}(v, *client, i)
		}
	}

	//wg.Wait()
	log.Println("阻塞....")
	<-chans[len(news)-1]
	log.Println("阻塞结束....")

	log.Println("爬取成功")
	log.Println("num goroutine: ", runtime.NumGoroutine())
	//os.Exit(0)
	return
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
