package command

import (
	"fmt"
)

type MysqlConvert struct {
	Options MysqlCommonOptions `group:"common"`
	Positional struct {
		Database string `description:"Database" required:"true"`
	} `positional-args:"true"`
	Charset   string           `long:"charset"      description:"MySQL charset"               default:"utf8mb4"`
	Collation string           `long:"collation"    description:"MySQL collation"             default:"utf8mb4_unicode_ci"`
}

func (conf *MysqlConvert) Execute(args []string) error {
	Logger.Main("Converting MySQL database \"%s\" to charset \"%s\" and collation \"%s\"", conf.Positional.Database, conf.Charset, conf.Collation)
	if err := conf.Options.Init(); err != nil {
		return err
	}


	defer NewSigIntHandler(func() {})()

	// Convert database
	Logger.Step("converting atabase \"%s\"", conf.Positional.Database)
	statement := fmt.Sprintf(
		"SET FOREIGN_KEY_CHECKS=0; ALTER DATABASE %s CHARACTER SET %s COLLATE %s",
		mysqlIdentifier(conf.Positional.Database),
		conf.Charset,
		conf.Collation,
	)
	conf.Options.ExecStatement("mysql", statement)


	// Convert tables
	tableList := conf.Options.GetTableList(conf.Positional.Database)
	for _, table := range tableList {
		Logger.Step("converting table \"%s\"", table)
		statement := fmt.Sprintf(
			"SET FOREIGN_KEY_CHECKS=0; ALTER TABLE %s.%s CONVERT TO CHARACTER SET %s COLLATE %s",
			mysqlIdentifier(conf.Positional.Database),
			mysqlIdentifier(table),
			conf.Charset,
			conf.Collation,
		)
		conf.Options.ExecStatement("mysql", statement)
	}

	return nil
}
