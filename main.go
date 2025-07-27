package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/metal3d/vymad/converters"
	"github.com/metal3d/vymad/converters/freemind"
	"github.com/metal3d/vymad/converters/minder"
	"github.com/metal3d/vymad/converters/vym"
	"github.com/metal3d/vymad/converters/xmind"
)

// TPL is the main file content template.
const (
	TPL = `% {{ .Title | html }}

{{ .Content | html }}

`
)

var (
	VERSION  = "master"
	RICHTEXT = false
)

func main() {
	v := flag.Bool("version", false, "print version")
	flag.BoolVar(&RICHTEXT, "richtext", RICHTEXT, "Try to parse richtext (for vym and xmind, automatic for freemind)")

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

	if file == "" {
		fmt.Println("You must provide a file")
		os.Exit(1)
	}

	ext := filepath.Ext(file)

	var converter converters.Converter
	switch ext {
	case ".mm":
		converter = freemind.ExecuteTpl
	case ".vym":
		converter = vym.ExecuteTpl
	case ".xmind":
		converter = xmind.ExecuteTpl
	case ".minder":
		converter = minder.ExecuteTpl
	default:
		fmt.Println("Unknown file extension:", ext)
		return
	}

	if err := converter(file, TPL, RICHTEXT); err != nil {
		fmt.Println("Error executing template:", err)
		os.Exit(1)
	}
}
