package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type PostgresSchemaDump struct {
	Options PostgresCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *PostgresSchemaDump) Execute(args []string) error {
	fmt.Println(fmt.Sprintf("Dumping PostgreSQL schema \"%s\" to \"%s\"", conf.Positional.Schema, conf.Positional.Filename))
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		fmt.Println(fmt.Sprintf(" - Using %s compression", conf.Options.dumpCompression))
	}

	cmd := shell.Cmd(conf.Options.PgDumpCommandBuilder(conf.Positional.Schema)...).Pipe(fmt.Sprintf("cat > %s", shell.Quote(conf.Positional.Filename)))
	cmd.Run()

	return nil
}
