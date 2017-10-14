package command

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/webdevops/go-shell"
)

type MysqlDebug struct {
	Options MysqlCommonOptions `group:"common"`
}

func (conf *MysqlDebug) Execute(args []string) error {
	Logger.Main("Starting MySQL query log")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	logfile := fmt.Sprintf("/tmp/mysql.debug.%d.log", r.Int63())

	conf.Options.Init()

	defer NewSigIntHandler(func() {
		Logger.Step("disabling mysql general log")
		conf.Options.ExecStatement("mysql", "SET GLOBAL general_log = 'OFF'")
		Logger.Step("removing log file")
		shell.Cmd(conf.Options.connection.CommandBuilder("rm", "-f", logfile)...).Run()
	})()

	Logger.Step("enabling mysql general log")
	conf.Options.ExecStatement("mysql", fmt.Sprintf("SET GLOBAL general_log_file = '%s'", logfile))
	conf.Options.ExecStatement("mysql", "SET GLOBAL general_log = 'ON'")

	Logger.Println("Starting log tail")
	Logger.Println("-----------------")
	cmd := shell.Cmd(conf.Options.connection.CommandBuilder("tail", "-n0", "-f", logfile)...)
	cmd.RunInteractive()
	return nil
}
