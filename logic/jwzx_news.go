// @author cold bin
// @date 2022/9/15

package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"jwzx-mail/conf"
	"jwzx-mail/model"
	"jwzx-mail/util"
	"log"
	"net/http"
	"strings"
)

const (
	// NewsListUrl 需要传入查询的偏移量和分页值，首页是1，默认每页大小为20
	NewsListUrl = "http://jwzx.cqupt.edu.cn/data/json_files.php?pageNo=%d&pageSize=%d&searchKey="

	// NewsContentUrl 需要拼接教务在线消息栏的消息id
	NewsContentUrl = "http://jwzx.cqupt.edu.cn/fileShowContent.php?id=%d"

	// NewsFile 需要拼接教务在线消息栏的消息id
	NewsFile = "http://jwzx.cqupt.edu.cn/fileDownLoadAttach.php?id=%s"
)

var ErrNotTheLatest = errors.New("[jwzx-mail]: not the latest news")

// GetNewNews
// TODO
// 	若文件更新就爬取所有更新的文件，返回所有更新的文件：怎么找到所有新的文件呢？
// 	先请求分页查询出前20条数据（一般新通知不会这么多），然后再依次取最新的那一批文件id，之后，再批量请求获取每个文件的具体内容
// 	然后返回
func GetNewNews() (contents []model.NewsContent, err error) {
	// 获取第一条信息
	var first model.ListInside
	first, err = GetJwzxFirstNews()
	if err != nil {
		return contents, err
	}

	// 若没有更新，就终止服务，直接返回
	if !conf.AConf.IsLatestFile(first.FileId) {
		return contents, ErrNotTheLatest
	}

	var preFileId = conf.AConf.LatestFileId
	// 更新最新文件的id

	if err = conf.AConf.UpdateFileId(first.FileId); err != nil {
		return contents, err
	}
	// 先获取最新的所有消息
	var newsList model.ListInsides
	newsList, err = GetJwzxNewsList(preFileId, conf.AConf.LatestFileId)
	if err != nil {
		return contents, err
	}

	// 取出最新消息对应的所有内容（一些头信息+网页主体+附件）
	var content model.NewsContent
	contents = make([]model.NewsContent, 0, 20)
	// 根据文件id，批量请求提取出数据，附件的数据还没有添加上
	for _, v := range newsList {
		content, err = GetJwzxContent(v.FileId)
		if err != nil {
			return contents, err
		}
		contents = append(contents, content)
	}

	return
}

// GetJwzxFirstNews 获取jwzx第一条信息
func GetJwzxFirstNews() (first model.ListInside, err error) {
	// 先获取第一条记录id
	url := fmt.Sprintf(NewsListUrl, 1, 1)

	var response *http.Response
	response, err = util.Get(url)
	if err != nil {
		return first, err
	}
	defer response.Body.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return first, err
	}

	var s model.ListOutside
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return first, err
	}

	first = s.Data[0]

	log.Println("logic GetNewNews: 解析出第一条消息 -> fileId:", first.FileId)
	return
}

// GetJwzxNewsList 把所有最新的文件返回
func GetJwzxNewsList(preFileId, LatestFileId int) (res model.ListInsides, err error) {
	// 先取20条
	url := fmt.Sprintf(NewsListUrl, 1, 20)

	var response *http.Response
	response, err = util.Get(url)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	var s model.ListOutside
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return res, err
	}

	// 筛选
	res = make(model.ListInsides, 0, len(s.Data))
	for _, v := range s.Data {
		if v.FileId > preFileId && v.FileId <= LatestFileId {
			res = append(res, v)
		} else {
			break
		}
	}
	// 反序
	res.Reverse()
	return
}

// GetJwzxContent 根据jwzx消息栏里的fileId来获取对应的网页内容及其附件，因为一般都不大，所以直接放到内存里
func GetJwzxContent(fileId int) (c model.NewsContent, err error) {
	// 发送请求
	response, err := http.Get(fmt.Sprintf(NewsContentUrl, fileId))
	if err != nil {
		log.Println("FindArticle: http.Get error:", err)
		return model.NewsContent{}, err
	}

	// 生成goquery的document结构体
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("FindArticle: goquery.NewDocumentFromReader error:", err)
		return model.NewsContent{}, err
	}

	// 寻找文章主体所在div
	mainPanel := document.Find("div#mainPanel")

	// 文章标题
	var article model.NewsContent
	article.Title = mainPanel.Find("h3").Text()

	// 获取时间、发布人、阅读量
	var dateAuthCountHtml string
	dateAuthCountHtml, err = mainPanel.Find("p").First().Html()
	log.Println("dateAuthCountHtml ->", dateAuthCountHtml)
	if err != nil {
		log.Println("FindArticle: dateAuthCountHtml find error:", err)
		return
	}

	// 主体内容
	var html string
	html, err = mainPanel.Find("div").First().Html()
	if err != nil {
		log.Println("FindArticle: html error:", err)
		return
	}

	html = dateAuthCountHtml + html
	html = strings.ReplaceAll(html, "\"", "'")
	html = strings.ReplaceAll(html, "\n", "")
	article.Content = html

	// 获取附件
	fileList := mainPanel.Find("ul")
	// 由于下面有固定的ul，故当ul个数 > 1时，才会有附件
	if len(fileList.Nodes) > 1 {
		fileList.First().Find("li").Each(func(i int, s *goquery.Selection) {
			s = s.Find("a")
			href, _ := s.Attr("href")
			index := strings.LastIndex(href, "=")

			// 获取具体附件的内容和返回类型
			var file model.Attachment
			file, err = GetJwzxAttachmentFile(href[index+1:])
			if err != nil {
				log.Println("FindArticle: GetJwzxAttachmentFile error:", err)
				return
			}

			article.Files = append(article.Files, model.Attachment{
				Id:          href[index+1:],
				Name:        s.Text(),
				InputStream: file.InputStream,
				Header:      file.Header,
			})
		})
	}

	c = article
	return
}

// GetJwzxAttachmentFile 根据附件的fileId，获取对应的附件
func GetJwzxAttachmentFile(fileId string) (file model.Attachment, err error) {
	file = model.Attachment{
		InputStream: make([]byte, 0, 100000),
		Header:      make(map[string]string),
	}
	fileUrl := fmt.Sprintf(NewsFile, fileId)

	resp, err := util.Get(fileUrl)
	if err != nil {
		log.Println("Get file error:", err)
		return model.Attachment{}, err
	}
	body := resp.Body
	defer body.Close()

	// 设置头部
	file.Header["Content-Disposition"] = resp.Header.Get("Content-Disposition")
	file.Header["Content-Length"] = resp.Header.Get("Content-Length")
	file.Header["Content-Type"] = resp.Header.Get("Content-Type")

	file.InputStream, err = ioutil.ReadAll(body)
	if err != nil {
		log.Println("Read file error:", err)
		return model.Attachment{}, err
	}

	return
}
