package command

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/webdevops/go-shell"
)

type MysqlSlowLog struct {
	Options MysqlCommonOptions `group:"common"`
	QueryTime int              `long:"querytime"    description:"Slow query time (seconds)"  default:"10"`
	QueryWithoutIndex bool     `long:"no-index"     description:"Log queries using NO index"`
}

func (conf *MysqlSlowLog) Execute(args []string) error {
	Logger.Main("Starting MySQL slow log")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	logfile := fmt.Sprintf("/tmp/mysql.debug.%d.log", r.Int63())

	conf.Options.Init()

	defer NewSigIntHandler(func() {
		Logger.Step("disabling mysql slow log")
		conf.Options.ExecStatement("mysql", "SET GLOBAL general_log = 'OFF'")
		conf.Options.ExecStatement("mysql", "SET GLOBAL log_queries_not_using_indexes = 'OFF'")
		Logger.Step("removing log file")
		shell.Cmd(conf.Options.connection.CommandBuilder("rm", "-f", logfile)...).Run()
	})()

	conf.Options.ExecStatement("mysql", fmt.Sprintf("SET GLOBAL slow_query_log_file = '%s'", logfile))
	conf.Options.ExecStatement("mysql", "SET GLOBAL slow_query_log = 'ON'")
	conf.Options.ExecStatement("mysql", fmt.Sprintf("SET GLOBAL long_query_time = %d", conf.QueryTime))

	if conf.QueryWithoutIndex {
		conf.Options.ExecStatement("mysql", "SET GLOBAL log_queries_not_using_indexes = 'ON'")
	} else {
		conf.Options.ExecStatement("mysql", "SET GLOBAL log_queries_not_using_indexes = 'OFF'")
	}

	Logger.Println("Starting log tail")
	Logger.Println("-----------------")
	cmd := shell.Cmd(conf.Options.connection.CommandBuilder("tail", "-n0", "-f", logfile)...)
	cmd.RunInteractive()
	return nil
}
