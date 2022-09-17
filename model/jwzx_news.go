// @author cold bin
// @date 2022/9/15

package model

import "sort"

type ListOutside struct {
	TotalPage int          `json:"totalPage"`
	Data      []ListInside `json:"data"`
}

type ListInside struct {
	TotalCount int    `json:"totalCount"`
	FileId     int    `json:"fileId"`
	Title      string `json:"title"`
	PubTime    string `json:"pubTime"`
	ReadCount  int    `json:"readCount"`
	Days       int    `json:"days"` // 通知发表天数，0表示当天发布，1表示昨天发布，以此类推
}

type ListInsides []ListInside

func (l ListInsides) Len() int {
	return len(l)
}

func (l ListInsides) Less(i, j int) bool {
	return l[i].FileId < l[j].FileId
}

func (l ListInsides) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ListInsides) Reverse() {
	sort.Sort(l)
}

type NewsContent struct {
	//Id        int          `json:"id"`
	Title string `json:"title"`
	//Date      string       `json:"date"`
	//ReadCount string       `json:"read_count"`
	//Author    string       `json:"author"`
	Content string       `json:"content"`
	Files   []Attachment `json:"files"`
}

type Attachment struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	InputStream []byte            `json:"input_stream"`
	Header      map[string]string `json:"header"`
}
