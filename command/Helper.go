package command

import (
	"os"
	"os/signal"
	"github.com/webdevops/go-shell"
	"fmt"
	"path/filepath"
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

func GetCompressionByFilename(file string) string {
	compression := "plain"
	fileext := filepath.Ext(file)

	switch fileext {
	case ".gz":
		fallthrough
	case ".gzip":
		compression = "gzip"

	case ".bz":
		fallthrough
	case ".bz2":
		fallthrough
	case ".bzip2":
		compression = "bzip2"

	case ".xz":
		fallthrough
	case ".lz":
		fallthrough
	case ".lzma":
		compression = "xz"
	}

	return compression
}
