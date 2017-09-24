package command

import (
	"os"
	"os/signal"
	"github.com/webdevops/go-shell"
	"fmt"
)

func NewSigIntHandler(callback func()) func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		shell.Panic = false
		fmt.Println("Starting termination as requested by user...")
	}()

	return func () {
		callback()
		fmt.Println("Terminated by user (SIGINT)")
	}
}
