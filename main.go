package main

import (
	"os"
	"log"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"github.com/webdevops/go-shell"
	"./logger"
	"./command"
	"os/signal"
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
	Verbose  []bool   `short:"v"  long:"verbose"      description:"verbose mode"`
}

func createArgparser() {
	var err error

	argparser = flags.NewParser(&opts, flags.Default)
	argparser.CommandHandler = func(command flags.Commander, args []string) error {
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

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			<-c

			// disable panic on SIGINT/SIGTERM
			shell.Panic = false
		}()

		return command.Execute(args)
	}
	
	argparser.AddCommand("version", "MySQL dump schema", "Backup MySQL schema to file", &command.Version{Name:Name, Version:Version, Author:Author})
	
	argparser.AddCommand("mysql:backup", "MySQL dump schema", "Backup MySQL schema to file", &command.MysqlBackup{})
	argparser.AddCommand("mysql:restore", "MySQL restore schema", "Restore MySQL schema from file", &command.MysqlRestore{})
	argparser.AddCommand("mysql:debug", "MySQL debug", "Show MySQL query log", &command.MysqlDebug{})
	argparser.AddCommand("mysql:slowlog", "MySQL slow query log", "Show MySQL slow query log", &command.MysqlSlowLog{})
	argparser.AddCommand("mysql:convert", "MySQL convert schema charset/collation", "Convert a schema to a charset and collation", &command.MysqlConvert{})

	argparser.AddCommand("postgres:backup", "PostgreSQL dump schema", "Backup PostgreSQL schema to file", &command.PostgresBackup{})
	argparser.AddCommand("postgres:restore", "PostgreSQL restore schema", "Restore PostgreSQL schema from file", &command.PostgresRestore{})

	argparser.AddCommand("typo3:stubs", "TYPO3 create file stubs", "", &command.Typo3Stubs{})
	argparser.AddCommand("typo3:beuser", "TYPO3 create BE user", "", &command.Typo3BeUser{})

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
}

func main() {
	createArgparser()
	os.Exit(0)
}
