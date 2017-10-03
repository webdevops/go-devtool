package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlSchemaRestore struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *MysqlSchemaRestore) Execute(args []string) error {
	fmt.Println(fmt.Sprintf("Restoring MySQL dump \"%s\" to schema \"%s\"", conf.Positional.Filename, conf.Positional.Schema))

	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		fmt.Println(fmt.Sprintf(" - Using %s decompression", conf.Options.dumpCompression))
	}

	conf.Options.ExecStatement("mysql", fmt.Sprintf("DROP DATABASE IF EXISTS %s", mysqlIdentifier(conf.Positional.Schema)))
	conf.Options.ExecStatement("mysql", fmt.Sprintf("CREATE DATABASE %s", mysqlIdentifier(conf.Positional.Schema)))
	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.MysqlRestoreCommandBuilder(conf.Positional.Schema)...)
	cmd.Run()

	return nil
}