package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/linterpreteur/brunchfeed/cmd"
)

func main() {
	actionMap := map[string]func(){
		"fetch": func() {
			id := flag.String("id", "", "A unique brunch account identifier")
			src := flag.String("src", "", "Path where data are saved")
			flag.Parse()

			brunchfeed.Fetch(brunchfeed.FetchParams{
				Id:  *id,
				Src: *src,
			})
		},
		"build": func() {
			template := flag.String("template", "", "Post output template string")
			src := flag.String("src", "", "Path which raw data are read from")
			dest := flag.String("dest", "", "Path which post data are saved into")
			flag.Parse()

			brunchfeed.Build(brunchfeed.BuildParams{
				Template: *template,
				Src:      *src,
				Dest:     *dest,
			})
		},
	}

	exit := func(message string) {
		const newline = "\n"
		var buffer bytes.Buffer
		buffer.WriteString(newline)
		buffer.WriteString(message)
		buffer.WriteString(newline)
		buffer.WriteString("Available commands are:")
		buffer.WriteString(newline)
		for command := range actionMap {
			buffer.WriteString(newline)
			buffer.WriteString("\tgo run main.go ")
			buffer.WriteString(command)
		}
		buffer.WriteString(newline)
		fmt.Println(buffer.String())
		os.Exit(1)
	}

	args := os.Args
	if len(args) < 2 {
		exit("Command is not specified.")
	}

	command := args[1]
	action, ok := actionMap[command]
	if !ok {
		exit("Unavailable command '" + command + "'.")
	}

	os.Args = args[1:] // Enables GNU style flag parsing
	action()
}
