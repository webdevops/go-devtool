package command

import (
	"os"
	"os/signal"
	"github.com/webdevops/go-shell"
	"fmt"
)

func NewSigIntHandler(callback func()) func() {
	isSigintExit := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		shell.Panic = false
		isSigintExit = true
		fmt.Println("Starting termination as requested by user...")
	}()

	return func () {
		callback()
		if (isSigintExit) {
			fmt.Println("Terminated by user (SIGINT)")
		}
	}
}
