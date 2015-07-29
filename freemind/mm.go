package freemind

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"github.com/metal3d/vymad/pandoc"
)

var PARTS = make([]string, 0)

type Map struct {
	Node Node
}

type Node struct {
	XMLName xml.Name `xml:"node"`
	Text    string   `xml:"TEXT,attr"`
	Content string   `xml:"richcontent"`
	Nodes   []Node   `xml:"node"`
}

func write(n *Node, level int) {

	PARTS = append(PARTS, fmt.Sprintf("%s %s", strings.Repeat("#", level+1), n.Text))
	if n.Content != "" {
		o, err := pandoc.Launch(n.Content, "markdown")
		if err != nil {
			panic(err)
		}
		PARTS = append(PARTS, o)
	}

	if len(n.Nodes) > 0 {
		for _, n := range n.Nodes {
			write(&n, level+1)
		}
	}

}

func ExecuteTpl(filename, tpl string, richtext bool) {

	content, _ := ioutil.ReadFile(filename)

	re := regexp.MustCompile(`<richcontent (.*)>`)

	s := re.ReplaceAllString(string(content), "<richcontent><![CDATA[")
	s = strings.Replace(s, "</richcontent>", "]]></richcontent>", -1)
	content = []byte(s)

	n := Map{}
	if err := xml.Unmarshal(content, &n); err != nil {
		fmt.Println(err)
	}

	title := n.Node.Text

	for _, n := range n.Node.Nodes {
		write(&n, 0)
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
