package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"text/template"
	"time"
)

var cmdTemplate *template.Template
var concurrentCount int
var showHelp bool

func init() {
	var err error
	flag.IntVar(&concurrentCount, "c", 1, "number of processes to run concurrently")
	flag.BoolVar(&showHelp, "-h", false, "print usage")

	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln("Need to specify command template as first argument")
	}

	if showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	cmdTemplateSource := args[0]
	cmdTemplate, err = template.New("cmd").Parse(cmdTemplateSource)

	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	input := os.Stdin
	runner := NewParallelRunner(concurrentCount)

	go processInput(input, runner)
	go runner.Start()
	runner.Wait()
}

func processInput(input io.Reader, runner *ParallelRunner) {
	decoder := json.NewDecoder(input)
	// var data map[string]interface{}
	var data interface{}
	for {
		err := decoder.Decode(&data)
		if err == io.EOF {
			runner.End()
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
			// case []interface{}:
		case map[string]interface{}:
			w := bytes.NewBuffer(buf)
			cmdTemplate.Execute(w, data)
			cmd := w.String()
			runner.Run(cmd)
		}
	}
}

type CmdRunner interface {
	Run(cmd string)
	End()
}

type ParallelRunner struct {
	// a command to run
	cmdChan chan interface{}
	// signal eof
	// eofChan  chan error
	runningCounts chan int
	doneChan      chan bool
}

func NewParallelRunner(concunrrency int) *ParallelRunner {
	runner := &ParallelRunner{
		cmdChan:       make(chan interface{}),
		doneChan:      make(chan bool),
		runningCounts: make(chan int, concunrrency),
	}
	return runner
}

func (r *ParallelRunner) Run(cmd string) {
	r.cmdChan <- cmd
}

func (r *ParallelRunner) End() {
	r.cmdChan <- io.EOF
}

func (r *ParallelRunner) Wait() {
	<-r.doneChan
}

func (r *ParallelRunner) Start() {
	var group sync.WaitGroup
loop:
	for {
		cmd := <-r.cmdChan
		switch cmd := cmd.(type) {
		case string:
			group.Add(1)
			r.runningCounts <- 1 // block if currently running at paralleism capacity
			go func() {
				log.Printf("run cmd: %v", cmd)
				runCommand(cmd)
				<-r.runningCounts
				group.Done()
			}()
			// stagger the next command a bit
			time.Sleep(20 * time.Millisecond)
		case error: // should be EOF
			break loop
		}
	}

	// wait for all running commands to finish
	group.Wait()

	r.doneChan <- true
}

func runCommand(cmdString string) {
	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Stdout = os.Stdout
	cmd.Run()
	// return
}
