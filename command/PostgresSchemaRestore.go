package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type PostgresSchemaRestore struct {
	Options PostgresCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *PostgresSchemaRestore) Execute(args []string) error {
	fmt.Println(fmt.Sprintf("Restoring PostgreSQL dump \"%s\" to schema \"%s\"", conf.Positional.Filename, conf.Positional.Schema))
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		fmt.Println(fmt.Sprintf(" - Using %s decompression", conf.Options.dumpCompression))
	}

	conf.Options.ExecStatement(fmt.Sprintf("DROP DATABASE IF EXISTS %s", postgresIdentifier(conf.Positional.Schema)))
	conf.Options.ExecStatement(fmt.Sprintf("CREATE DATABASE %s", postgresIdentifier(conf.Positional.Schema)))

	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.PostgresRestoreCommandBuilder(conf.Positional.Schema)...)
	cmd.Run()

	return nil
}
