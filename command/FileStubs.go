package command

import (
	"fmt"
	"bufio"
	"os"
	"errors"
	"path/filepath"
	"github.com/webdevops/go-stubfilegenerator"
)

type FileStubs struct {
	Positional struct {
		SourceFile string `description:"Source file"`
	} `positional-args:"true"`
	RootPath    string  `long:"path"                  description:"Prefix path for stub files"  default:"./"`
	SourceStdin bool    `long:"stdin"                 description:"Use stdin as file list source"`
	Force       bool    `short:"f"  long:"force"      description:"Overwrite existing files"`
}

func (conf *FileStubs) Execute(args []string) error {
	var fileSource *os.File

	if conf.Positional.SourceFile == "" && !conf.SourceStdin {
		return errors.New("[ERROR] No source defined, either specify file as argument or use --stdin")
	}

	defer NewSigIntHandler(func() {});

	stubgen := stubfilegenerator.NewStubGenerator()

	if conf.Force {
		// use --force
		stubgen.Overwrite = true
	}

	if conf.SourceStdin {
		// use --stdin
		fileSource = os.Stdin
	} else if conf.Positional.SourceFile != "" {
		// use --source
		f, err := os.Open(conf.Positional.SourceFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		fileSource = f
	}
	
	scanner := bufio.NewScanner(fileSource)
	for scanner.Scan() {
		relPath := scanner.Text()
		path := scanner.Text()

		if conf.RootPath != "" {
			path = filepath.Join(conf.RootPath, path)
		}

		fmt.Println(fmt.Sprintf(" >> %s", path))
		stubgen.TemplateVariables["PATH"] = relPath
		stubgen.Generate(path)
	}

	return nil
}
