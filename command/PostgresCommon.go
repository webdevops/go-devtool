package command

import (
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type PostgresCommonOptions struct {
	Hostname string `long:"hostname"`
	Port     string `short:"P" long:"port"`
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

func (conf *PostgresCommonOptions) PsqlCommandBuilder(args ...string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{}

	cmd = append(cmd, "psql")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["PGPASSWORD"] = conf.Password
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return connection.RawShellCommandBuilder(cmd...)
}

func (conf *PostgresCommonOptions) PgDumpCommandBuilder(database string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{}

	cmd = append(cmd, "pg_dump")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["PGPASSWORD"] = conf.Password
	}

	cmd = append(cmd, shell.Quote(database))

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

func (conf *PostgresCommonOptions) PgDumpAllCommandBuilder() []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{}

	cmd = append(cmd, "pg_dumpall", "-c")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["PGPASSWORD"] = conf.Password
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

func (conf *PostgresCommonOptions) PostgresRestoreCommandBuilder(args ...string) []interface{} {
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

	cmd = append(cmd, "psql")

	if conf.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.Hostname))
	}

	if conf.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.Port))
	}

	if conf.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.Username))
	}

	if conf.Password != "" {
		connection.Environment["PGPASSWORD"] = conf.Password
	}

	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	return connection.RawShellCommandBuilder(cmd...)
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
	Logger.Item("using docker container \"%s\"", containerId)

	containerEnv := connectionClone.DockerGetEnvironment(containerId)

	// try to guess user/password
	if conf.Username == "" {
		// get superuser pass from env
		if pass, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
			if user, ok := containerEnv["POSTGRES_USER"]; ok {
				Logger.Item("using postgres superadmin account \"%s\" (from env:POSTGRES_USER and env:POSTGRES_PASSWORD)", user)
				conf.Username = user
				conf.Password = pass
			} else {
				Logger.Item("using postgres superadmin account \"postgres\" (from env:POSTGRES_PASSWORD)")
				// only password available
				conf.Username = "postgres"
				conf.Password = pass
			}
		}
	}

}
