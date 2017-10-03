package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type PostgresRestore struct {
	Options PostgresCommonOptions `group:"common"`
	Positional struct {
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *PostgresRestore) Execute(args []string) error {
	fmt.Println(fmt.Sprintf("Restoring PostgreSQL dump \"%s\"", conf.Positional.Filename))
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		fmt.Println(fmt.Sprintf(" - Using %s decompression", conf.Options.dumpCompression))
	}

	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.PostgresRestoreCommandBuilder("postgres")...)
	cmd.Run()

	return nil
}
