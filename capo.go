package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
)

const (
	commands_file = "commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Starting capo...")

	var commands []string

	if len(os.Args) == 2 {
		if os.Args[1] == "--foreman" {
			commands = readCommands("Procfile", true)
		} else {
			commands = readCommands(os.Args[1], false)
		}
	} else {
		commands = readCommands(commands_file, false)
	}

	errorChan := make(chan error)
	processList := make(chan *os.Process, len(commands))
	done := make(chan bool, len(commands))
	sig := make(chan os.Signal)

	doneCount := 0
	allDone := len(commands)

	signal.Notify(sig, os.Interrupt, os.Kill)

	for _, command := range commands {
		go runCommandString(command, processList, errorChan, done)
	}

	for doneCount != allDone {
		select {
		case <-done:
			doneCount++
		case <-sig:
			fmt.Println("Terminating from signal")
			stopProcesses(processList)
		case <-errorChan:
			stopProcesses(processList)
		}
	}
}

func readCommands(file string, foreman_mode bool) []string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Could not read file", file)
		os.Exit(1)
	}

	commands := make([]string, 0)

	for _, s := range strings.Split(string(content), "\n") {
		if s != "" {
			if foreman_mode {
				c := strings.SplitN(s, ":", 2)[1]
				commands = append(commands, strings.Trim(c, " "))
			} else {
				commands = append(commands, s)
			}
		}
	}

	return commands
}

func runCommandString(command string, processList chan *os.Process, errorChan chan error, done chan bool) {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()

	processList <- cmd.Process

	err := cmd.Wait()

	if err == nil {
		done <- true
	} else {
		errorChan <- err
	}

}

// There is a potential race condition where a process may be added after we're finished killing the existing ones.
func stopProcesses(processList chan *os.Process) {
	n := len(processList)
	for i := 0; i < n; i++ {
		process := <-processList
		// NOTE: process may already be done. Ideally it should skip those
		process.Kill()
	}
	os.Exit(1)
}
