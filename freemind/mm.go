package freemind

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
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
		cmd := exec.Command("/usr/bin/pandoc", "-f", "html", "-t", "markdown")

		stdin, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := cmd.Start(); err != nil {
			panic(err)
		}

		stdin.Write([]byte(n.Content))
		stdin.Close()

		if err != nil {
			fmt.Println(err)
		}
		o, err := ioutil.ReadAll(stdout)
		if err != nil {
			panic(err)
		}
		PARTS = append(PARTS, string(o))
	}

	if len(n.Nodes) > 0 {
		for _, n := range n.Nodes {
			write(&n, level+1)
		}
	}

}

func Open(filename, tpl string) {

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
