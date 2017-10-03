package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type PostgresBackupAll struct {
	Options PostgresCommonOptions `group:"common"`
	Positional struct {
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *PostgresBackupAll) Execute(args []string) error {
	fmt.Println(fmt.Sprintf("Dumping PostgreSQL to \"%s\"", conf.Positional.Filename))
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		fmt.Println(fmt.Sprintf(" - Using %s compression", conf.Options.dumpCompression))
	}

	cmd := shell.Cmd(conf.Options.PgDumpAllCommandBuilder()...).Pipe(fmt.Sprintf("cat > %s", shell.Quote(conf.Positional.Filename)))
	cmd.Run()

	return nil
}
