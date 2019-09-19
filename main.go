package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var history []string

func main() {
	reader := bufio.NewReader(os.Stdin)
	var pwd string
	var currentUser *user.User
	var username string
	var userHome string
	var hostname string
	var input string
	var err error

	for {

		if currentUser, err = user.Current(); err == nil {
			username = currentUser.Name
			userHome = currentUser.HomeDir

		} else {
			username = "Donkey"
		}

		if hostname, err = os.Hostname(); err != nil {
			hostname = "localhost"
		}

		if pwd, err = os.Getwd(); err != nil {
			pwd = ""
		}

		c := color.New(color.BgCyan, color.FgBlack, color.Bold)
		c.Printf("%s@%s=>%s$", username, hostname, pwd)

		input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		if err = runCmd(input, userHome); err != nil {
			fmt.Fprint(os.Stderr, err, "\n")

		}

	}
}

func runCmd(input string, userHome string) (err error) {

	input = strings.TrimSuffix(input, "\n")

	history = append(history, input)
	args := strings.Split(input, " ")
	histroyR, _ := regexp.Compile("^![0-9]+")
	//histroyUP, _ := regexp.Compile("|")
	//histroyDOWN, _ := regexp.Compile("^[[B]")

	//fmt.Print(args[0])
	//fmt.Print(r.MatchString(args[0]))

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return os.Chdir(userHome)
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(1)
	case "":
		return
	case "history":
		if len(args) < 2 {
			for i, his := range history {
				fmt.Printf("%d %s \n", i+1, his)
			}
		} else if len(args) == 2 && args[1] == "-c" {
			history = nil
			fmt.Fprintf(os.Stdout, "History cleared \n")
		} else {
			return errors.New("sub-command not supported")
		}
		return

	}
	if histroyR.MatchString(args[0]) {
		num, _ := strconv.Atoi(strings.TrimPrefix(args[0], "!"))
		fmt.Fprintln(os.Stdout, history[num-1])
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()

}
