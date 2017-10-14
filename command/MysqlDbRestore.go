package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlDbRestore struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Database string `description:"Database" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *MysqlDbRestore) Execute(args []string) error {
	Logger.Main("Restoring MySQL dump \"%s\" to database \"%s\"", conf.Positional.Filename, conf.Positional.Database)

	conf.Options.Init()

	defer NewSigIntHandler(func() {})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		Logger.Step("using %s decompression", conf.Options.dumpCompression)
	}

	conf.Options.ExecStatement("mysql", fmt.Sprintf("DROP DATABASE IF EXISTS %s", mysqlIdentifier(conf.Positional.Database)))
	conf.Options.ExecStatement("mysql", fmt.Sprintf("CREATE DATABASE %s", mysqlIdentifier(conf.Positional.Database)))
	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.MysqlRestoreCommandBuilder(conf.Positional.Database)...)
	cmd.Run()

	return nil
}
