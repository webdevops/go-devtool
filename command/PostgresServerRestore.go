package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type PostgresServerRestore struct {
	Options PostgresCommonOptions `group:"common"`
	Positional struct {
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *PostgresServerRestore) Execute(args []string) error {
	Logger.Main("Restoring PostgreSQL dump \"%s\"", conf.Positional.Filename)
	if err := conf.Options.Init(); err != nil {
		return err
	}

	defer NewSigIntHandler(func() {})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		Logger.Step("using %s decompression", conf.Options.dumpCompression)
	}

	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.PostgresRestoreCommandBuilder("postgres")...)
	cmd.Run()

	return nil
}
