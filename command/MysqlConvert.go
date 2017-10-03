package command

import (
	"fmt"
)

type MysqlConvert struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Schema string `description:"Schema" required:"true"`
	} `positional-args:"true"`
	Charset   string           `long:"charset"      description:"MySQL charset"               default:"utf8"`
	Collation string           `long:"collation"    description:"MySQL collation"             default:"utf8_general_ci"`
}

func (conf *MysqlConvert) Execute(args []string) error {
	fmt.Println(fmt.Sprintf("Converting MySQL schema \"%s\" to charset \"%s\" and collation \"%s\"", conf.Positional.Schema, conf.Charset, conf.Collation))
	conf.Options.Init()

	defer NewSigIntHandler(func() {
	})()

	// Convert database
	fmt.Println(fmt.Sprintf(" - converting schema \"%s\"", conf.Positional.Schema))
	statement := fmt.Sprintf(
		"SET FOREIGN_KEY_CHECKS=0; ALTER DATABASE %s CHARACTER SET %s COLLATE %s",
		mysqlIdentifier(conf.Positional.Schema),
		conf.Charset,
		conf.Collation,
	)
	conf.Options.ExecStatement("mysql", statement)


	// Convert tables
	tableList := conf.Options.GetTableList(conf.Positional.Schema)
	for _, table := range tableList {
		fmt.Println(fmt.Sprintf(" - converting table \"%s\"", table))
		statement := fmt.Sprintf(
			"SET FOREIGN_KEY_CHECKS=0; ALTER TABLE %s.%s CONVERT TO CHARACTER SET %s COLLATE %s",
			mysqlIdentifier(conf.Positional.Schema),
			mysqlIdentifier(table),
			conf.Charset,
			conf.Collation,
		)
		conf.Options.ExecStatement("mysql", statement)
	}

	return nil
}
