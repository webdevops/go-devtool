package command

import (
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlServerRestore struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Filename string `description:"Backup filename" required:"1"`
	} `positional-args:"true"`
}

func (conf *MysqlServerRestore) Execute(args []string) error {
	Logger.Main("Restoring MySQL dump \"%s\"", conf.Positional.Filename)
	conf.Options.Init()

	defer NewSigIntHandler(func() {})()

	shell.SetDefaultShell("bash")

	conf.Options.dumpCompression = GetCompressionByFilename(conf.Positional.Filename)
	if (conf.Options.dumpCompression != "") {
		Logger.Step("using %s decompression", conf.Options.dumpCompression)
	}

	cmd := shell.Cmd(fmt.Sprintf("cat %s", shell.Quote(conf.Positional.Filename))).Pipe(conf.Options.MysqlRestoreCommandBuilder()...)
	cmd.Run()

	return nil
}
