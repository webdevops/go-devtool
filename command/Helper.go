package command

import (
	"os"
	"os/signal"
	"fmt"
	"time"
	"path/filepath"
	"math/rand"
	"bufio"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
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

func GetDockerEnvList(connection commandbuilder.Connection, containerId string) map[string]string {
	ret := map[string]string{}

	cmd := shell.Cmd(connection.CommandBuilder("docker", "inspect", "-f", "{{range .Config.Env}}{{println .}}{{end}}", containerId)...)
	envList := cmd.Run().Stdout.String()

	scanner := bufio.NewScanner(strings.NewReader(envList))
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2)

		if len(split) == 2 {
			varName, varValue := split[0], split[1]

			ret[varName] = varValue
		}
	}

	return ret
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
