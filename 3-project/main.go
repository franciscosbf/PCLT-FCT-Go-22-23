package main

import (
	"cpl_go_proj22/builder"
	"cpl_go_proj22/parser"
	"cpl_go_proj22/utils"
	"flag"
	"fmt"
	"log"
	"os"
)

func oneShot(c chan *builder.Msg) {
	m := <-c
	if m.Type == builder.BuildSuccess {
		fmt.Println("Build was a success.")
	} else {
		fmt.Printf("Something went wrong with the build: %v\n", m.Err)
	}
}

func main() {
	path := flag.String("d", "", "Files location, (current directory by default)")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: project [-d] <location>")
		os.Exit(0)
	}
	fileName := args[0]

	dFile, err := parser.ParseFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	var scan *utils.FileScan
	if scan, err = utils.NewFileScan(*path); err != nil {
		log.Fatal(err.Error())
	}

	ch := builder.MakeController(dFile, scan)
	oneShot(ch)
}
