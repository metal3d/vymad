// Package vym provides functions to parse Vym mind maps and convert them to markdown.
package vym

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/metal3d/vymad/pandoc"
)

// PARTS contains document lines.
var (
	PARTS    = make([]string, 0)
	RICHTEXT = false
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
	c, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	x := VymMap{}
	if err := xml.Unmarshal(c, &x); err != nil {
		panic(err)
	}

	return &x
}

// Parse branches.
func doParts(b *Branch, level int) {
	heading := b.Heading
	PARTS = append(PARTS, fmt.Sprintf("%s %s", strings.Repeat("#", level+1), heading))

	if RICHTEXT {
		var (
			o   string
			err error
		)
		if o, err = pandoc.Launch(b.VymNote, "markdown"); err != nil {
			panic(err)
		}
		PARTS = append(PARTS, o)

	} else {
		PARTS = append(PARTS, b.VymNote)
	}
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

func ExecuteTpl(filename, tpl string, richtext bool) error {
	RICHTEXT = richtext

	// try to open file
	r, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name[len(f.Name)-4:] != ".xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		x := xmlBuildStruct(rc)
		parseVymTree(x, tpl)
		break
	}
	return nil
}
