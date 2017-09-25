package main

import (
	"os"
	"log"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/webdevops/go-shell"
	"./logger"
	"./command"
)

const (
	Name    = "godevtool"
	Author  = "webdevops.io"
	Version = "0.1.0"
)

var (
	Logger *logger.SyncLogger
	argparser *flags.Parser
	args []string
)

var opts struct {
	Verbose            []bool   `short:"v"  long:"verbose"      description:"verbose mode"`
	ShowVersion        bool     `short:"V"  long:"version"      description:"show version and exit"`
	ShowOnlyVersion    bool     `           long:"dumpversion"  description:"show only version number and exit"`
	ShowHelp           bool     `           long:"help"         description:"show this help message"`
}
func createArgparser() {
	var err error
	shell.VerboseFunc = func(c *shell.Command) {
		fmt.Println(c.ToString())
	}

	argparser = flags.NewParser(&opts, flags.Default)
	argparser.AddCommand("mysql:backup", "MySQl dump schema", "Backup MySQL schema to file", &command.MysqlBackup{});
	argparser.AddCommand("mysql:restore", "MySQl restore schema", "Restore MySQL schema from file", &command.MysqlRestore{});
	argparser.AddCommand("mysql:debug", "MySQl debug", "Show MySQL query log", &command.MysqlDebug{});
	argparser.AddCommand("mysql:slowlog", "MySQl slow query log", "Show MySQL slow query log", &command.MysqlSlowLog{});
	argparser.AddCommand("mysql:convert", "MySQl convert schema charset/collation", "Convert a schema to a charset and collation", &command.MysqlConvert{});

	args, err = argparser.Parse()

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	// --dumpversion
	if opts.ShowOnlyVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	// --version
	if opts.ShowVersion {
		fmt.Println(fmt.Sprintf("%s version %s", Name, Version))
		fmt.Println(fmt.Sprintf("Copyright (C) 2017 %s", Author))
		os.Exit(0)
	}
}

func main() {
	createArgparser()

	switch {
	case len(opts.Verbose) >= 2:
		shell.Trace = true
		shell.TracePrefix = "[CMD] "
		Logger = logger.GetInstance(argparser.Command.Name, log.Ldate|log.Ltime|log.Lshortfile)
		fallthrough
	case len(opts.Verbose) >= 1:
		logger.Verbose = true
		shell.VerboseFunc = func(c *shell.Command) {
			Logger.Command(c.ToString())
		}
		fallthrough
	default:
		if Logger == nil {
			Logger = logger.GetInstance(argparser.Command.Name, 0)
		}
	}

	os.Exit(0)
}
