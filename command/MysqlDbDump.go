package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlDbDump struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Database string `description:"Database" required:"1"`
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *MysqlDbDump) Execute(args []string) error {
	Logger.Main("Dumping MySQL database \"%s\" to \"%s\"", conf.Positional.Database, conf.Positional.Filename)
	if err := conf.Options.Init(); err != nil {
		return err
	}

	defer NewSigIntHandler(func() {})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		Logger.Step("using %s compression", conf.Options.dumpCompression)
	}

	cmd := shell.Cmd(conf.Options.MysqlDumpCommandBuilder(shell.Quote(conf.Positional.Database))...).Pipe(fmt.Sprintf("cat > %s", shell.Quote(conf.Positional.Filename)))
	cmd.Run()

	return nil
}
