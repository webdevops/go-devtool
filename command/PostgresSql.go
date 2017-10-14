package command

import (
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type PostgresSql struct {
	Options PostgresCommonOptions `group:"common"`
	Interactive bool           `short:"i" long:"interactive"     description:"Run interactive shell"`
}

func (conf *PostgresSql) Execute(args []string) error {
	conf.Options.Init()

	defer NewSigIntHandler(func() {})()

	commandbuilder.ConnectionDockerArguments = append(commandbuilder.ConnectionDockerArguments, "-t")
	if conf.Interactive {
		commandbuilder.ConnectionSshArguments = []string{"-oPasswordAuthentication=no"}
		cmd := shell.Cmd(conf.Options.PsqlCommandBuilder(args...)...)
		cmd.RunInteractive()
	} else {
		cmd := shell.Cmd(conf.Options.PsqlCommandBuilder(args...)...)
		cmd.RunInteractive()
	}

	return nil
}
