package vym

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

// PARTS contains document lines.
var PARTS = make([]string, 0)

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
func parseVymTree(v *VymMap, tpl string) {

	// parse each mapcenter branches
	for _, b := range v.MapCenter.Branch {
		doParts(&b, 0)
	}

	t, _ := template.New("doc").Parse(tpl)
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

func Open(filename string, tpl string) {
	// try to open file
	r, err := zip.OpenReader(filename)
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
		parseVymTree(x, tpl)
		break
	}

}
