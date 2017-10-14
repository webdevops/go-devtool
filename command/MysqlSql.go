package command

import (
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type MysqlSql struct {
	Options MysqlCommonOptions `group:"common"`
	Interactive bool           `short:"i" long:"interactive"     description:"Run interactive shell"`
}

func (conf *MysqlSql) Execute(args []string) error {
	conf.Options.Init()

	defer NewSigIntHandler(func() {})()

	commandbuilder.ConnectionDockerArguments = append(commandbuilder.ConnectionDockerArguments, "-t")

	if conf.Interactive {
		commandbuilder.ConnectionSshArguments = []string{"-oPasswordAuthentication=no"}

		cmd := shell.Cmd(conf.Options.MysqlInteractiveCommandBuilder(args...)...)
		cmd.RunInteractive()
	} else {
		cmd := shell.Cmd(conf.Options.MysqlCommandBuilder(args...)...)
		cmd.RunInteractive()
	}

	return nil
}
