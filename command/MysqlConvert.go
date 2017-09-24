package command

import (
	"fmt"
)

type MysqlConvert struct {
	Options MysqlCommonOptions `group:"common"`
	Schema    string           `long:"schema"       description:"Database/Schema to convert"  required:"true"`
	Charset   string           `long:"charset"      description:"MySQL charset"               default:"utf8"`
	Collation string           `long:"collation"    description:"MySQL collation"             default:"utf8_general_ci"`
}

func (conf *MysqlConvert) Execute(args []string) error {

	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	// Convert database
	fmt.Println(fmt.Sprintf(" - converting schema \"%s\"", conf.Schema))
	statement := fmt.Sprintf(
		"SET FOREIGN_KEY_CHECKS=0; ALTER DATABASE %s CHARACTER SET %s COLLATE %s",
		mysqlIdentifier(conf.Schema),
		conf.Charset,
		conf.Collation,
	)
	conf.Options.ExecMySqlStatement(statement)


	// Convert tables
	tableList := conf.Options.GetTableList(conf.Schema)
	for _, table := range tableList {
		fmt.Println(fmt.Sprintf(" - converting table \"%s\"", table))
		statement := fmt.Sprintf(
			"SET FOREIGN_KEY_CHECKS=0; ALTER TABLE %s.%s CONVERT TO CHARACTER SET %s COLLATE %s",
			mysqlIdentifier(conf.Schema),
			mysqlIdentifier(table),
			conf.Charset,
			conf.Collation,
		)
		conf.Options.ExecMySqlStatement(statement)
	}

	return nil
}
