package command

import (
	"time"
	"math/rand"
	"fmt"
	"github.com/webdevops/go-shell"
)

type MysqlSlowLog struct {
	Options MysqlCommonOptions `group:"common"`
	QueryTime int              `long:"querytime"    description:"Slow query time (seconds)"  default:"10"`
	QueryWithoutIndex bool     `long:"no-index"     description:"Log queries using NO index"`
}

func (conf *MysqlSlowLog) Execute(args []string) error {
	fmt.Println("Starting MySQL slow log")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	logfile := fmt.Sprintf("/tmp/mysql.debug.%d.log", r.Int63());

	conf.Options.Init()

	defer NewSigIntHandler(func() {
		conf.Options.ExecMySqlStatement("SET GLOBAL general_log = 'OFF'")
		conf.Options.ExecMySqlStatement("SET GLOBAL log_queries_not_using_indexes = 'OFF'")
		shell.Cmd(conf.Options.connection.CommandBuilder("rm", "-f", logfile)...).Run()
	})()

	conf.Options.ExecMySqlStatement(fmt.Sprintf("SET GLOBAL slow_query_log_file = '%s'", logfile))
	conf.Options.ExecMySqlStatement("SET GLOBAL slow_query_log = 'ON'")
	conf.Options.ExecMySqlStatement(fmt.Sprintf("SET GLOBAL long_query_time = %d", conf.QueryTime))

	if conf.QueryWithoutIndex {
		conf.Options.ExecMySqlStatement("SET GLOBAL log_queries_not_using_indexes = 'ON'")
	} else {
		conf.Options.ExecMySqlStatement("SET GLOBAL log_queries_not_using_indexes = 'OFF'")
	}

	fmt.Println("Starting log tail")
	fmt.Println("-----------------")
	cmd := shell.Cmd(conf.Options.connection.CommandBuilder("tail", "-n0", "-f", logfile)...)
	cmd.RunInteractive()
	return nil
}
