package command

import (
	"bufio"
	"os"
	"errors"
	"github.com/webdevops/go-stubfilegenerator"
	"strings"
)

type FileStubs struct {
	Positional struct {
		SourceFile string `description:"Source file"`
	} `positional-args:"true"`
	RootPath    string  `long:"path"                  description:"Prefix path for stub files"`
	SourceStdin bool    `long:"stdin"                 description:"Use stdin as file list source"`
	Force       bool    `short:"f"  long:"force"      description:"Overwrite existing files"`
}

func (conf *FileStubs) Execute(args []string) error {
	var err error
	var fileSource *os.File

	if conf.Positional.SourceFile == "" && !conf.SourceStdin {
		return errors.New("[ERROR] No source defined, either specify file as argument or use --stdin")
	}

	rootPath, err := conf.GetRootPath()
	if err != nil {
		return err
	}

	defer NewSigIntHandler(func() {})

	stubgen := stubfilegenerator.NewStubGenerator()
	err = stubgen.SetRootPath(rootPath)
	if err != nil {
		return err
	}

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
		// get relative path
		relPath := strings.TrimSpace(scanner.Text())
		if relPath == "" {
			continue
		}

		err := stubgen.Generate(relPath)
		if err == nil {
			Logger.Printlnf(" >> %s", relPath)
		} else {
			Logger.Printlnf(" [!!!] %s: %v", relPath, err)
		}
	}

	return nil
}

func (conf *FileStubs) GetRootPath() (string, error) {
	if conf.RootPath == "" {
		// use current working dir as root path
		rootPath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		conf.RootPath = rootPath
	}

	return conf.RootPath, nil
}
