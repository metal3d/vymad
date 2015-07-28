package xmind

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"
)

var PARTS = make([]string, 0)

type Xmap struct {
	XMLName xml.Name `xml:"xmap-content"`
	Sheet   Sheet    `xml:"sheet"`
}

type Sheet struct {
	XMLName xml.Name `xml:"sheet"`
	Topic   Topic    `xml:"topic"`
}

type Topic struct {
	XMLName xml.Name `xml:"topic"`
	Title   string   `xml:"title"`

	Content string `xml:"notes>plain"`

	Children []Topic `xml:"children>topics>topic"`
}

func buildXMLTree(rc io.ReadCloser) *Xmap {

	defer rc.Close()

	c, err := ioutil.ReadAll(rc)
	if err != nil {
		panic(err)
	}

	x := Xmap{}
	xml.Unmarshal(c, &x)
	return &x

}

func parse(t Topic, level int) {
	PARTS = append(PARTS, fmt.Sprintf("%s %s", strings.Repeat("#", level+1), t.Title))
	PARTS = append(PARTS, t.Content)

	for _, ch := range t.Children {
		parse(ch, level+1)
	}
}

func Open(file, tpl string) {
	r, err := zip.OpenReader(file)
	if err != nil {
		panic(err)
	}

	defer r.Close()

	title := ""
	for _, f := range r.File {
		if f.Name == "content.xml" {
			rc, err := f.Open()
			if err != nil {
				panic(err)
			}

			x := buildXMLTree(rc)
			title = x.Sheet.Topic.Title
			for _, t := range x.Sheet.Topic.Children {
				parse(t, 0)
			}
			break
		}
	}

	t, _ := template.New("doc").Parse(tpl)
	var b []byte
	buff := bytes.NewBuffer(b)
	t.Execute(buff, map[string]string{
		"Title":   title,
		"Content": strings.Join(PARTS, "\n\n"),
	})

	fmt.Println(buff)

}
