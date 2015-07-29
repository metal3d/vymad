package xmind

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"github.com/metal3d/vymad/pandoc"
)

var (
	PARTS    = make([]string, 0)
	STYLES   = make(map[string][]string)
	RICHTEXT = false
)

type XMapStyle struct {
	XMLName xml.Name `xml:"xmap-styles"`
	Styles  []Style  `xml:"styles>style"`
}

type Style struct {
	XMLName    xml.Name   `xml:"style"`
	Id         string     `xml:"id,attr"`
	Properties []Property `xml:"text-properties"`
}

type Property struct {
	XMLName    xml.Name `xml:"text-properties"`
	Weight     string   `xml:"font-weight,attr"`
	Decoration string   `xml:"text-decoration,attr"`
	Style      string   `xml:"font-style,attr"`
}

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
	Html    string `xml:"notes>html"`

	Children []Topic `xml:"children>topics>topic"`
}

func buildXMLTree(rc io.ReadCloser) *Xmap {

	defer rc.Close()

	c, err := ioutil.ReadAll(rc)

	// replace html content to chardata
	re := regexp.MustCompile(`<notes><html>(.*?)</html>`)
	c = re.ReplaceAll(c, []byte("<notes><html><![CDATA[$1]]></html>"))

	if err != nil {
		panic(err)
	}

	x := Xmap{}
	xml.Unmarshal(c, &x)
	return &x

}

func build(ff *zip.File) string {

	rc, err := ff.Open()
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	x := buildXMLTree(rc)
	title := x.Sheet.Topic.Title
	for _, t := range x.Sheet.Topic.Children {
		parse(t, 0)
	}
	return title
}

func parseStyles(f *zip.File) {

	s, err := f.Open()
	if err != nil {
		panic(err)
	}

	x := XMapStyle{}
	content, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}
	xml.Unmarshal(content, &x)

	for _, s := range x.Styles {

		mdmark := []string{""}
		for _, p := range s.Properties {
			if p.Decoration == "underline" {
				mdmark = append(mdmark, "u")
			}
			if p.Style == "italic" {
				mdmark = append(mdmark, "em")
			}
			if p.Weight == "bold" {
				mdmark = append(mdmark, "strong")
			}

			STYLES[s.Id] = append(STYLES[s.Id], mdmark...)
		}

	}
	defer s.Close()
}

func replaceStyles(content string) string {

	content = strings.Replace(content, "<span\n", "<span", -1)
	for id, s := range STYLES {
		found := true
		re := regexp.MustCompile(`(.*)<span.+?style-id="` + id + `".*?>(.+?)</span>`)

		for found {
			brep := ""
			erep := ""
			for _, r := range s {
				if len(r) == 0 {
					continue
				}
				brep = brep + "<" + r + ">"
				erep = "</" + r + ">" + erep
			}

			matches := re.FindAllStringSubmatch(content, -1)
			found = len(matches) > 0
			content = re.ReplaceAllString(content, "$1"+brep+"$2"+erep)
		}
	}

	return content
}

func parse(t Topic, level int) {
	PARTS = append(PARTS, fmt.Sprintf("%s %s", strings.Repeat("#", level+1), t.Title))
	if RICHTEXT {
		var (
			o   string
			err error
		)
		if o, err = pandoc.Launch(t.Html, "html"); err != nil {
			panic(err)
		}
		o = replaceStyles(o)
		if o, err = pandoc.Launch(o, "markdown"); err != nil {
			panic(err)
		}

		PARTS = append(PARTS, o)
	} else {
		PARTS = append(PARTS, t.Content)
	}

	for _, ch := range t.Children {
		parse(ch, level+1)
	}
}

func ExecuteTpl(file, tpl string, richtext bool) {
	RICHTEXT = richtext
	r, err := zip.OpenReader(file)
	if err != nil {
		panic(err)
	}

	defer r.Close()

	ff := new(zip.File)
	for _, f := range r.File {
		switch f.Name {
		case "content.xml":
			ff = f
		case "styles.xml":
			parseStyles(f)
		}
	}

	title := build(ff)

	t, _ := template.New("doc").Parse(tpl)
	var b []byte
	buff := bytes.NewBuffer(b)
	t.Execute(buff, map[string]string{
		"Title":   title,
		"Content": strings.Join(PARTS, "\n\n"),
	})

	fmt.Println(buff)

}
