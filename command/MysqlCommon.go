package command

import (
	"bufio"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
	"fmt"
	"regexp"
	"encoding/xml"
	"errors"
)

type MysqlCommonOptions struct {
	Hostname string `          long:"hostname"`
	Port     string `short:"P" long:"port"`
	Username string `short:"u" long:"user"`
	Password string `short:"p" long:"password"`
	Docker   string `          long:"docker"`
	SSH      string `          long:"ssh"`

	connection commandbuilder.Connection
	dumpCompression string
}

type xmlMysqlQueryResult struct {
	XMLName xml.Name `xml:"resultset"`
	Row []xmlMysqlQueryRow `xml:"row"`
}

type xmlMysqlQueryRow struct {
	Field []xmlMysqlQueryField `xml:"field"`
}

type xmlMysqlQueryField struct {
	Name string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

func  (row *xmlMysqlQueryRow) GetList() map[string]string {
	ret := map[string]string{}

	for _, field := range row.Field {
		ret[field.Name] = field.Value
	}
	
	return ret
}

func  (row *xmlMysqlQueryRow) GetField(name string) (string, error) {
	for _, field := range row.Field {
		if name == field.Name {
			return field.Value, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Field %s not found", name))
}


func mysqlQuote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func mysqlIdentifier(value string) string {
	return "`" + strings.Replace(value, "`", "\\`", -1) + "`"
}

func  (conf *MysqlCommonOptions) Init() {
	Logger.Step("init connection settings")
	
	if conf.SSH != "" {
		conf.connection.Hostname = conf.SSH
		Logger.Item("using ssh connection \"%s\"", conf.SSH)
	}

	if conf.Docker != "" {
		conf.connection.Docker = conf.Docker
		conf.InitDockerSettings()
	}
}

func (conf *MysqlCommonOptions) MysqlInteractiveCommandBuilder(args ...string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{""}

	if conf.Hostname != "" {
		cmd = append(cmd, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, shell.Quote("-P" + conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["MYSQL_PWD"] = conf.Password
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return connection.RawCommandBuilder("mysql", cmd...)
}

func (conf *MysqlCommonOptions) MysqlCommandBuilder(args ...string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{"-NB"}

	if conf.Hostname != "" {
		cmd = append(cmd, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, shell.Quote("-P" + conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["MYSQL_PWD"] = conf.Password
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return connection.RawCommandBuilder("mysql", cmd...)
}

func (conf *MysqlCommonOptions) MysqlDumpCommandBuilder(args ...string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{"mysqldump", "--single-transaction"}

	if conf.Hostname != "" {
		cmd = append(cmd, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, shell.Quote("-P" + conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["MYSQL_PWD"] = conf.Password
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	switch conf.dumpCompression {
	case "gzip":
		cmd = append(cmd, "| gzip")
	case "bzip2":
		cmd = append(cmd, "| bzip2")
	case "xz":
		cmd = append(cmd, "| xz --compress --stdout")
	}

	return connection.RawShellCommandBuilder(cmd...)
}

func (conf *MysqlCommonOptions) MysqlRestoreCommandBuilder(args ...string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{}

	switch conf.dumpCompression {
	case "gzip":
		cmd = append(cmd, "gzip -dc |")
	case "bzip2":
		cmd = append(cmd, "bzcat |")
	case "xz":
		cmd = append(cmd, "xzcat |")
	}

	cmd = append(cmd, "mysql", "-NB")

	if conf.Hostname != "" {
		cmd = append(cmd, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, shell.Quote("-P" + conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["MYSQL_PWD"] = conf.Password
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return connection.RawShellCommandBuilder(cmd...)
}

func (conf *MysqlCommonOptions) ExecStatement(database string, statement string) string {
	cmd := shell.Cmd(conf.MysqlCommandBuilder(shell.Quote(database), "-e", shell.Quote(statement))...)
	return cmd.Run().Stdout.String()
}

func (conf *MysqlCommonOptions) ExecQuery(database string, statement string) xmlMysqlQueryResult {
	re := regexp.MustCompile("\\n")
	re.ReplaceAllString(statement, " ")

	cmd := shell.Cmd(conf.MysqlCommandBuilder(shell.Quote(database), "--xml", "-e", shell.Quote(statement))...)
	stdout := cmd.Run().Stdout.String()

	// parse result as xml
	var result xmlMysqlQueryResult
	xml.Unmarshal([]byte(stdout), &result)

	return result
}

func  (conf *MysqlCommonOptions) GetTableList(database string) []string {
	var ret []string

	sql := fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = %s", mysqlQuote(database))
	output := conf.ExecStatement("mysql", sql)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		ret = append(ret, line)
	}

	return ret
}

func  (conf *MysqlCommonOptions) InitDockerSettings() {
	containerName := conf.connection.Docker

	connectionClone := conf.connection.Clone()
	connectionClone.Docker = ""
	connectionClone.Type  = "auto"

	containerId := connectionClone.DockerGetContainerId(containerName)
	Logger.Item("using docker container \"%s\"", containerId)

	containerEnv := connectionClone.DockerGetEnvironment(containerId)

	if conf.Username == "" {
		if val, ok := containerEnv["MYSQL_ROOT_PASSWORD"]; ok {
			// get root pass from env
			if conf.Username == "" && conf.Password == "" {
				Logger.Item("using mysql root account (from env:MYSQL_ROOT_PASSWORD)")
				conf.Username = "root"
				conf.Password = val
			}
		} else if val, ok := containerEnv["MYSQL_ALLOW_EMPTY_PASSWORD"]; ok {
			// get root without password from env
			if val == "yes" && conf.Username == "" {
				Logger.Item("using mysql root account (from env:MYSQL_ALLOW_EMPTY_PASSWORD)")
				conf.Username = "root"
				conf.Password = ""
			}
		} else if user, ok := containerEnv["MYSQL_USER"]; ok {
			if pass, ok := containerEnv["MYSQL_PASSWORD"]; ok {
				if conf.Username == "" && conf.Password == "" {
					Logger.Item("using mysql user account \"%s\" (from env:MYSQL_USER and env:MYSQL_PASSWORD)", user)
					conf.Username = user
					conf.Password = pass
				}
			}
		}
	}
}
