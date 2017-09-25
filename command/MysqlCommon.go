package command

import (
	"bufio"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
	"fmt"
)

type MysqlCommonOptions struct {
	Hostname string `long:"hostname"`
	Username string `short:"u" long:"user"`
	Password string `short:"p" long:"password"`
	Docker   string `          long:"docker"`
	SSH      string `          long:"ssh"`

	connection commandbuilder.Connection
}

func mysqlQuote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func mysqlIdentifier(value string) string {
	return "`" + strings.Replace(value, "`", "\\`", -1) + "`"
}

func  (conf *MysqlCommonOptions) Init() {
	if conf.SSH != "" {
		conf.connection.Hostname = conf.SSH
	}

	if conf.Docker != "" {
		conf.connection.Docker = conf.Docker
		conf.InitDockerSettings()
	}
}

func (conf *MysqlCommonOptions) MysqlCommandBuilder(args ...string) []interface{} {
	mysqlArgs := []string{"-N", "-B"}

	if conf.Hostname != "" {
		mysqlArgs = append(mysqlArgs, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Username != "" {
		mysqlArgs = append(mysqlArgs, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		mysqlArgs = append(mysqlArgs, shell.Quote("-p" + conf.Password))
	}

	if len(args) > 0 {
		mysqlArgs = append(mysqlArgs, args...)
	}

	return conf.connection.CommandBuilder("mysql", mysqlArgs...)
}

func (conf *MysqlCommonOptions) MysqlDumpCommandBuilder(schema string) []interface{} {
	mysqlArgs := []string{"--single-transaction"}

	if conf.Hostname != "" {
		mysqlArgs = append(mysqlArgs, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Username != "" {
		mysqlArgs = append(mysqlArgs, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		mysqlArgs = append(mysqlArgs, shell.Quote("-p" + conf.Password))
	}

	mysqlArgs = append(mysqlArgs, shell.Quote(schema))
	return conf.connection.CommandBuilder("mysqldump", mysqlArgs...)
}


func (conf *MysqlCommonOptions) ExecMySqlStatement(statement string) string {
	cmd := shell.Cmd(conf.MysqlCommandBuilder("-e", shell.Quote(statement))...)
	return cmd.Run().Stdout.String()
}

func  (conf *MysqlCommonOptions) GetTableList (schema string) []string {
	var ret []string

	output := conf.ExecMySqlStatement(fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = %s", mysqlQuote(schema)))

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

	cmd := shell.Cmd(connectionClone.CommandBuilder("docker", "inspect",  "-f", shell.Quote("{{range .Config.Env}}{{println .}}{{end}}"), shell.Quote(containerId))...)
	envList := cmd.Run().Stdout.String()

	scanner := bufio.NewScanner(strings.NewReader(envList))
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2)
		if len(split) == 2 {
			varName, varValue := split[0], split[1]

			if varName == "MYSQL_ROOT_PASSWORD" {
				conf.Username = "root"
				conf.Password = varValue
				conf.Hostname = ""
			}
		}
	}
}
