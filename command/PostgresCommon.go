package command

import (
	"bufio"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
	"fmt"
)

type PostgresCommonOptions struct {
	Hostname string `long:"hostname"`
	Username string `short:"u" long:"user"`
	Password string `short:"p" long:"password"`
	Docker   string `          long:"docker"`
	SSH      string `          long:"ssh"`

	connection commandbuilder.Connection
	dumpCompression string
}

func postgresQuote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func postgresIdentifier(value string) string {
	return "\"" + strings.Replace(value, "\"", "\\\"", -1) + "\""
}

func  (conf *PostgresCommonOptions) Init() {
	if conf.SSH != "" {
		conf.connection.Hostname = conf.SSH
		fmt.Println(fmt.Sprintf(" - Using ssh connection \"%s\"", conf.SSH))
	}

	if conf.Docker != "" {
		conf.connection.Docker = conf.Docker
		conf.InitDockerSettings()
	}
}

func (conf *PostgresCommonOptions) PsqlCommandBuilder(args ...string) []interface{} {
	cmd := []string{}

	if conf.Password != "" {
		cmd = append(cmd, "PGPASSWORD=" + shell.Quote(conf.Password))
	}

	cmd = append(cmd, "psql")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return conf.connection.RawShellCommandBuilder(cmd...)
}

func (conf *PostgresCommonOptions) PgDumpCommandBuilder(schema string) []interface{} {
	cmd := []string{}

	if conf.Password != "" {
		cmd = append(cmd, "PGPASSWORD=" + shell.Quote(conf.Password))
	}

	cmd = append(cmd, "pg_dump")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	cmd = append(cmd, shell.Quote(schema))

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

func (conf *PostgresCommonOptions) PostgresRestoreCommandBuilder(args ...string) []interface{} {
	cmd := []string{}

	switch conf.dumpCompression {
	case "gzip":
		cmd = append(cmd, "gzip -dc |")
	case "bzip2":
		cmd = append(cmd, "bzcat |")
	case "xz":
		cmd = append(cmd, "xzcat |")
	}

	if conf.Password != "" {
		cmd = append(cmd, "PGPASSWORD=" + shell.Quote(conf.Password))
	}

	cmd = append(cmd, "pg_dump")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return conf.connection.RawShellCommandBuilder(cmd...)
}

func (conf *PostgresCommonOptions) ExecStatement(statement string) string {
	cmd := shell.Cmd(conf.PsqlCommandBuilder("postgres", "-c", shell.Quote(statement))...)
	return cmd.Run().Stdout.String()
}

func  (conf *PostgresCommonOptions) InitDockerSettings() {
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

			if varName == "POSTGRES_USER" && conf.Username == "" {
				conf.Username = varValue
			}

			if varName == "POSTGRES_PASSWORD" && conf.Password == ""  {
				conf.Password = varValue
			}
		}
	}
}
