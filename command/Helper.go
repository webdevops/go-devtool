package command

import (
	"os"
	"os/signal"
	"fmt"
	"time"
	"path/filepath"
	"math/rand"
)

func NewSigIntHandler(callback func()) func() {
	isSigintExit := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func randomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
