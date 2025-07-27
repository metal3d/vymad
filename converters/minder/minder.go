// Package minder provides functionality to parse and convert XML mind maps into a structured format.
package minder

import (
	"encoding/xml"
	"fmt"
	"html"
	"html/template"
	"os"
	"strings"
)

var parts = []string{}

type Text struct {
	Data string `xml:"data,attr"`
}

type Node struct {
	ID    string `xml:"id,attr"`
	Title Text   `xml:"nodename>text"`
	Note  string `xml:"nodenote"`
	Nodes []Node `xml:"nodes>node"`
}

type Minder struct {
	Nodes []Node `xml:"nodes>node"`
}

func doPatrs(n *Node, level int) {
	// decore the node title that may have html encoded chars like &#xe9;
	n.Title.Data = html.UnescapeString(n.Title.Data)
	n.Title.Data = strings.ReplaceAll(n.Title.Data, "\n", " ")
	n.Title.Data = strings.ReplaceAll(n.Title.Data, "\r", " ")
	n.Title.Data = strings.ReplaceAll(n.Title.Data, "\t", " ")
	n.Title.Data = strings.ReplaceAll(n.Title.Data, "  ", " ")
	n.Title.Data = strings.TrimSpace(n.Title.Data)

	for _, node := range n.Nodes {
		parts = append(parts, fmt.Sprintf("%s %s", strings.Repeat("#", level), node.Title.Data))
		if node.Note != "" {
			parts = append(parts, node.Note)
		}
		doPatrs(&node, level+1)
	}
}

func ExecuteTpl(filename, tpl string, richtext bool) error {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	minder := Minder{}
	err = xml.NewDecoder(xmlFile).Decode(&minder)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	title := minder.Nodes[0].Title.Data
	doPatrs(&minder.Nodes[0], 1)

	t := template.New("minder")
	t.Funcs(template.FuncMap{
		"html": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	t.Parse(tpl)

	t.Execute(os.Stdout, map[string]string{
		"Title":   title,
		"Content": strings.Join(parts, "\n\n"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
