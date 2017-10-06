package main

import (
	"os"
	"log"
	"fmt"
	"os/signal"
	"runtime/debug"
	flags "github.com/jessevdk/go-flags"
	"github.com/webdevops/go-shell"
	"./logger"
	"./command"
)

const (
	// application informations
	Name    = "godevtool"
	Author  = "webdevops.io"
	Version = "0.3.0"

	// self update informations
	GithubOrganization  = "webdevops"
	GithubRepository    = "go-devtool"
	GithubAssetTemplate = "gdt-%ARCH%-%OS%"
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
	
	argparser.AddCommand("version", "Show version", "Show current app version", &command.Version{Name:Name, Version:Version, Author:Author})
	argparser.AddCommand("self-update", "Self update", "Run self update of this application", &command.SelfUpdate{GithubOrganization:GithubOrganization, GithubRepository:GithubRepository, GithubAssetTemplate:GithubAssetTemplate})

	argparser.AddCommand("mysql:debug", "MySQL debug", "Show MySQL query log", &command.MysqlDebug{})
	argparser.AddCommand("mysql:slowlog", "MySQL slow query log", "Show MySQL slow query log", &command.MysqlSlowLog{})
	argparser.AddCommand("mysql:dump", "MySQL dump instance", "Backup MySQL instance (all schemas) to file", &command.MysqlDump{})
	argparser.AddCommand("mysql:restore", "MySQL restore instance", "Restore MySQL instance (all schemas) from file", &command.MysqlRestore{})

	argparser.AddCommand("mysql:schema:dump", "MySQL dump schema", "Backup MySQL schema to file", &command.MysqlSchemaDump{})
	argparser.AddCommand("mysql:schema:restore", "MySQL restore schema", "Restore MySQL schema from file", &command.MysqlSchemaRestore{})
	argparser.AddCommand("mysql:schema:convert", "MySQL convert schema charset/collation", "Convert a schema to a charset and collation", &command.MysqlConvert{})

	argparser.AddCommand("postgres:dump", "PostgreSQL dump instance", "Backup PostgreSQL schema to file", &command.PostgresDump{})
	argparser.AddCommand("postgres:restore", "PostgreSQL restore instance", "Restore PostgreSQL instance from file", &command.PostgresRestore{})
	argparser.AddCommand("postgres:schema:dump", "PostgreSQL dump schema", "Backup PostgreSQL schema to file", &command.PostgresSchemaDump{})
	argparser.AddCommand("postgres:schema:restore", "PostgreSQL restore schema", "Restore PostgreSQL schema from file", &command.PostgresSchemaRestore{})

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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			if len(opts.Verbose) >= 2 {
				fmt.Println(r)
				debug.PrintStack()
			} else {
				fmt.Println(r)
			}
			os.Exit(255)
		}
	}()

	createArgparser()
	os.Exit(0)
}
