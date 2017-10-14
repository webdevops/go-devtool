package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type PostgresDbRestore struct {
	Options PostgresCommonOptions `group:"common"`
	Positional struct {
		Database string `description:"Database" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *PostgresDbRestore) Execute(args []string) error {
	Logger.Main("Restoring PostgreSQL dump \"%s\" to database \"%s\"", conf.Positional.Filename, conf.Positional.Database)
	conf.Options.Init()

	defer NewSigIntHandler(func() {})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		Logger.Step("using %s decompression", conf.Options.dumpCompression)
	}

	conf.Options.ExecStatement(fmt.Sprintf("DROP DATABASE IF EXISTS %s", postgresIdentifier(conf.Positional.Database)))
	conf.Options.ExecStatement(fmt.Sprintf("CREATE DATABASE %s", postgresIdentifier(conf.Positional.Database)))

	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.PostgresRestoreCommandBuilder(conf.Positional.Database)...)
	cmd.Run()

	return nil
}
