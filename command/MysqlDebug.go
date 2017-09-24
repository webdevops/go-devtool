package command

import (
	"time"
	"math/rand"
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlDebug struct {
	Options MysqlCommonOptions `group:"common"`
}

func (conf *MysqlDebug) Execute(args []string) error {

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	logfile := fmt.Sprintf("/tmp/mysql.debug.%d.log", r.Int63());

	conf.Options.Init()

	defer NewSigIntHandler(func() {
		conf.Options.ExecMySqlStatement("SET GLOBAL general_log = 'OFF'")
		shell.Cmd(conf.Options.connection.CommandBuilder("rm", "-f", logfile)...).Run()
	})()

	conf.Options.ExecMySqlStatement(fmt.Sprintf("SET GLOBAL general_log_file = '%s'", logfile))
	conf.Options.ExecMySqlStatement("SET GLOBAL general_log = 'ON'")

	fmt.Println("Starting log tail")
	fmt.Println("-----------------")
	cmd := shell.Cmd(conf.Options.connection.CommandBuilder("tail", "-n0", "-f", logfile)...)
	cmd.RunInteractive()
	return nil
}
