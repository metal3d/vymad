package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// TPL is the main file content template.
const (
	TPL = `% {{ .Title }}

{{ .Content }}

`
)

var (
	// PARTS contains document lines.
	PARTS = make([]string, 0)

	VERSION = "master"
)

type VymMap struct {
	MapCenter MapCenter `xml:"mapcenter"`
}

type MapCenter struct {
	XMLName xml.Name `xml:"mapcenter"`
	Heading string   `xml:"heading"`
	Branch  []Branch `xml:"branch"`
}

type Branch struct {
	XMLName xml.Name `xml:"branch"`
	Heading string   `xml:"heading"`
	VymNote string   `xml:"vymnote"`

	Branches []Branch `xml:"branch"`
}

// build XML tree and return *VymMap.
func xmlBuildStruct(r io.ReadCloser) *VymMap {

	defer r.Close()
	c, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	x := VymMap{}
	xml.Unmarshal(c, &x)
	return &x

}

// Parse branches.
func doParts(b *Branch, level int) {

	heading := b.Heading
	PARTS = append(PARTS, fmt.Sprintf("%s %s", strings.Repeat("#", level+1), heading))
	PARTS = append(PARTS, b.VymNote)

	// parse children branches
	if len(b.Branches) > 0 {
		for _, b := range b.Branches {
			doParts(&b, level+1)
		}
	}

}

// parse the VymMap tree to build markdown.
func parse(v *VymMap) {

	// parse each mapcenter branches
	for _, b := range v.MapCenter.Branch {
		doParts(&b, 0)
	}

	t, _ := template.New("doc").Parse(TPL)
	var b []byte
	buff := bytes.NewBuffer(b)

	err := t.Execute(buff, map[string]string{
		"Title":   v.MapCenter.Heading,
		"Content": strings.Join(PARTS, "\n\n"),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buff)

}

func main() {

	v := flag.Bool("version", false, "print version")
	flag.Usage = func() {
		fmt.Println("Usage of " + os.Args[0])
		fmt.Println(os.Args[0] + " [options] file")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *v {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	file := flag.Arg(0)
	fmt.Println(file)
	os.Exit(0)

	if file == "" {
		fmt.Println("You must provide a vym file")
		os.Exit(1)
	}

	// try to open file
	r, err := zip.OpenReader(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name[len(f.Name)-4:] != ".xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			panic(err)
		}
		x := xmlBuildStruct(rc)
		parse(x)
		break
	}
}
