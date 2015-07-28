package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/metal3d/vymad/freemind"
	"github.com/metal3d/vymad/vym"
	"github.com/metal3d/vymad/xmind"
)

// TPL is the main file content template.
const (
	TPL = `% {{ .Title }}

{{ .Content }}

`
)

var (
	VERSION = "master"
)

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

	if file == "" {
		fmt.Println("You must provide a vym file")
		os.Exit(1)
	}

	if file[len(file)-3:] == ".mm" { // Freemind
		freemind.Open(file, TPL)
	} else if file[len(file)-4:] == ".vym" { // Vym
		vym.Open(file, TPL)
	} else if file[len(file)-6:] == ".xmind" { //xmind
		xmind.Open(file, TPL)
	}

}
