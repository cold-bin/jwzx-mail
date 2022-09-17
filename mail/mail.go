// @author cold bin
// @date 2022/9/15

// Package mail
// 邮件服务
package mail

import (
	"bytes"
	"errors"
	"github.com/jordan-wright/email"
	"jwzx-mail/conf"
	"jwzx-mail/model"
	"log"
	"net/smtp"
)

type Client struct {
	EPool *email.Pool
}

// NewClient email连接池
func NewClient() *Client {
	c := conf.AConf

	var p, err = email.NewPool(
		c.SendMailServerAddr,
		c.SendMailPoolGet,
		smtp.PlainAuth("", c.SendMailServerQQ, c.QQMailAuthCode, c.SendMailServerHost),
	)
	if err != nil {
		log.Fatal("failed to create pool:", err)
	}
	return &Client{p}
}

// SendHtmlWithAttachments 使用前，先将html和附件填充好
func (c *Client) SendHtmlWithAttachments(e *email.Email) error {
	log.Println("	正在发送邮件...")
	return c.EPool.Send(e, -1)
}

var ErrEmptyMailHeader = errors.New("empty `from` or `tos`")

func PutHeader(from string, tos []string, e *email.Email) error {
	if from == "" || len(tos) == 0 {
		return ErrEmptyMailHeader
	}
	e.From = from
	e.To = tos

	return nil
}

func PutTitle(subject string, e *email.Email) {
	e.Subject = conf.AConf.SendMailTitlePrefix + subject
}

// PutHtml 将html放到邮件正文
func PutHtml(html []byte, e *email.Email) {
	e.HTML = html
}

func PutTos(tos []string, e *email.Email) {
	e.To = tos
}

// PutAttachments 将附件内容放到邮件里
func PutAttachments(files []model.Attachment, e *email.Email) error {
	if len(e.Attachments) != 0 {
		return nil
	}
	log.Println("	有", len(files), "个附件")
	for _, f := range files {
		log.Println("		当前附件: ", f.Name)
		attach, err := e.Attach(bytes.NewReader(f.InputStream), f.Name, f.Header["Content-Type"])
		if err != nil {
			return err
		}
		e.Attachments = append(e.Attachments, attach)
	}

	return nil
}
