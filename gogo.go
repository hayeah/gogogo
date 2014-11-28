package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"
)

var cmdTemplate *template.Template

func init() {
	var err error
	if len(os.Args) < 2 {
		log.Fatalln("Need to specify command template as first argument")
	}

	cmdTemplateSource := os.Args[1]
	cmdTemplate, err = template.New("cmd").Parse(cmdTemplateSource)

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	input := os.Stdin
	decoder := json.NewDecoder(input)
	// var data map[string]interface{}
	var data interface{}
	for {
		err := decoder.Decode(&data)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
			continue
		}

		var buf []byte

		switch data := data.(type) {
		default:
			// do nothing
		case map[string]interface{}:
			w := bytes.NewBuffer(buf)
			cmdTemplate.Execute(w, data)
			cmd := w.String()
			log.Printf("run cmd: %v", cmd)
			runCommand(cmd)
		}
	}
}

func runCommand(cmdString string) (err error) {
	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return
}
