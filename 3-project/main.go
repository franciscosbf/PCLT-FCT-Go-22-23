package main

import (
	"cpl_go_proj22/builder"
	"cpl_go_proj22/parser"
	"flag"
	"fmt"
	"os"
	"time"
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
	watch := flag.Bool("w", false, "set flag to true to maintain watch over dependencies.")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: project [-w] <depfile>")
		os.Exit(0)
	}
	fileName := args[0]

	dFile, err := parser.ParseFile(fileName)
	if err != nil {
		_ = fmt.Errorf(err.Error())
		os.Exit(1)
	}

	ch := builder.MakeController(dFile)
	if *watch {
		fmt.Println("Starting build in watch mode.")
		for {
			oneShot(ch)
			time.Sleep(time.Second)
		}
	} else {
		oneShot(ch)
	}

}
