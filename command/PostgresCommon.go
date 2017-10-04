package command

import (
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

func (conf *PostgresCommonOptions) PgDumpAllCommandBuilder() []interface{} {
	cmd := []string{}

	if conf.Password != "" {
		cmd = append(cmd, "PGPASSWORD=" + shell.Quote(conf.Password))
	}

	cmd = append(cmd, "pg_dumpall", "-c")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
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

	containerEnv := GetDockerEnvList(connectionClone, containerId)

	// get user from env
	if val, ok := containerEnv["POSTGRES_USER"]; ok {
		if conf.Username == "" {
			conf.Username = val
		}
	}

	// get user from env
	if val, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
		if conf.Username == "" {
			conf.Password = val
		}
	}
}
