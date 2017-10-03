package command

import (
	"bufio"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
	"fmt"
	"regexp"
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

func mysqlQuote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func mysqlIdentifier(value string) string {
	return "`" + strings.Replace(value, "`", "\\`", -1) + "`"
}

func  (conf *MysqlCommonOptions) Init() {
	if conf.SSH != "" {
		conf.connection.Hostname = conf.SSH
		fmt.Println(fmt.Sprintf(" - Using ssh connection \"%s\"", conf.SSH))
	}

	if conf.Docker != "" {
		conf.connection.Docker = conf.Docker
		conf.InitDockerSettings()
	}
}

func (conf *MysqlCommonOptions) MysqlCommandBuilder(args ...string) []interface{} {
	cmd := []string{"-N", "-B"}

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
		cmd = append(cmd, shell.Quote("-p" + conf.Password))
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return conf.connection.CommandBuilder("mysql", cmd...)
}

func (conf *MysqlCommonOptions) MysqlDumpCommandBuilder(args ...string) []interface{} {
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
		cmd = append(cmd, shell.Quote("-p" + conf.Password))
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

	return conf.connection.RawShellCommandBuilder(cmd...)
}

func (conf *MysqlCommonOptions) MysqlRestoreCommandBuilder(args ...string) []interface{} {
	cmd := []string{}

	switch conf.dumpCompression {
	case "gzip":
		cmd = append(cmd, "gzip -dc |")
	case "bzip2":
		cmd = append(cmd, "bzcat |")
	case "xz":
		cmd = append(cmd, "xzcat |")
	}

	cmd = append(cmd, "mysql", "-N", "-B")

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
		cmd = append(cmd, shell.Quote("-p" + conf.Password))
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return conf.connection.RawShellCommandBuilder(cmd...)
}

func (conf *MysqlCommonOptions) ExecStatement(database string, statement string) string {
	cmd := shell.Cmd(conf.MysqlCommandBuilder(shell.Quote(database), "-e", shell.Quote(statement))...)
	return cmd.Run().Stdout.String()
}

func (conf *MysqlCommonOptions) ExecQuery(database string, statement string) map[int][]string {
	ret := map[int][]string{}

	re := regexp.MustCompile("\\n")
	re.ReplaceAllString(statement, " ")

	cmd := shell.Cmd(conf.MysqlCommandBuilder(shell.Quote(database), "-e", shell.Quote(statement))...)
	stdout := cmd.Run().Stdout.String()

	resultRegex := regexp.MustCompile("\t+")
	scanner := bufio.NewScanner(strings.NewReader(stdout))
	lineNumber := 0
	for scanner.Scan() {
		ret[lineNumber] = resultRegex.Split(scanner.Text(), -1)
		lineNumber++
	}

	return ret
}

func  (conf *MysqlCommonOptions) GetTableList (schema string) []string {
	var ret []string

	output := conf.ExecStatement("mysql", fmt.Sprintf("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = %s", mysqlQuote(schema)))

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
	fmt.Println(fmt.Sprintf(" - Using docker container \"%s\"", containerId))

	cmd := shell.Cmd(connectionClone.CommandBuilder("docker", "inspect",  "-f", shell.Quote("{{range .Config.Env}}{{println .}}{{end}}"), shell.Quote(containerId))...)
	envList := cmd.Run().Stdout.String()

	scanner := bufio.NewScanner(strings.NewReader(envList))
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2)
		if len(split) == 2 {
			varName, varValue := split[0], split[1]

			if varName == "MYSQL_ROOT_PASSWORD" && conf.Username == "" && conf.Password == "" {
				conf.Username = "root"
				conf.Password = varValue
				conf.Hostname = ""
			}
		}
	}
}
