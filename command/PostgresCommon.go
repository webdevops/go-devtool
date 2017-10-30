package command

import (
	"fmt"
	"errors"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type PostgresCommonOptions struct {
	SSH      string `long:"ssh"`
	Docker   string `long:"docker"`
	Postgres string `long:"postgres"`
	PostgresOptions struct {
		Hostname string `long:"hostname"`
		Port     string `long:"port"`
		Username string `long:"user"`
		Password string `long:"password"`
	} `group:"postgres" namespace:"postgres"`

	connection commandbuilder.Connection
	dumpCompression string
}

func postgresQuote(value string) string {
	return "'" + strings.Replace(value, "'", "\\'", -1) + "'"
}

func postgresIdentifier(value string) string {
	return "\"" + strings.Replace(value, "\"", "\\\"", -1) + "\""
}

func  (conf *PostgresCommonOptions) Init() error {
	Logger.Step("init connection settings")

	// --ssh
	if conf.SSH != "" {
		if err := conf.connection.SetSsh(conf.SSH); err != nil {
			return err
		}
		Logger.Item("using ssh connection \"%s\"", conf.SSH)
	}

	// --docker
	if conf.Docker != "" {
		if err := conf.connection.SetDocker(conf.Docker); err != nil {
			return err
		}
		conf.InitDockerSettings()
	}

	// --mysql
	// parse DSN/URL value
	if conf.Postgres != "" {
		postgresConf, err := commandbuilder.ParseArgument(conf.Postgres)
		if err != nil {
			return err
		}

		if postgresConf.Scheme != "mysql" {
			return errors.New(fmt.Sprintf("Scheme \"%v\" is not allowed, only mysql is supported in --mysql", postgresConf.Scheme))
		}

		if postgresConf.Hostname != "" {
			conf.PostgresOptions.Hostname = postgresConf.Hostname
		}

		if postgresConf.Port != "" {
			conf.PostgresOptions.Port = postgresConf.Port
		}

		if postgresConf.Username != "" {
			conf.PostgresOptions.Username = postgresConf.Username
		}

		if postgresConf.Password != "" {
			conf.PostgresOptions.Username = postgresConf.Password
		}
	}

	return nil
}

func (conf *PostgresCommonOptions) PsqlCommandBuilder(args ...string) []interface{} {
	connection := conf.connection.Clone()
	cmd := []string{}

	cmd = append(cmd, "psql")

	if conf.PostgresOptions.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.PostgresOptions.Hostname))
	}

	if conf.PostgresOptions.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.PostgresOptions.Port))
	}

	if conf.PostgresOptions.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.PostgresOptions.Username))
	}

	if conf.PostgresOptions.Password != "" {
		connection.Environment.Set("PGPASSWORD", conf.PostgresOptions.Password)
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

	if conf.PostgresOptions.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.PostgresOptions.Hostname))
	}

	if conf.PostgresOptions.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.PostgresOptions.Port))
	}

	if conf.PostgresOptions.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.PostgresOptions.Username))
	}

	if conf.PostgresOptions.Password != "" {
		connection.Environment.Set("PGPASSWORD", conf.PostgresOptions.Password)
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

	if conf.PostgresOptions.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.PostgresOptions.Hostname))
	}

	if conf.PostgresOptions.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.PostgresOptions.Port))
	}

	if conf.PostgresOptions.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.PostgresOptions.Username))
	}

	if conf.PostgresOptions.Password != "" {
		connection.Environment.Set("PGPASSWORD", conf.PostgresOptions.Password)
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

	if conf.PostgresOptions.Hostname != "" {
		cmd = append(cmd, "-h", shell.Quote(conf.PostgresOptions.Hostname))
	}

	if conf.PostgresOptions.Port != "" {
		cmd = append(cmd, "-p", shell.Quote(conf.PostgresOptions.Port))
	}

	if conf.PostgresOptions.Username != "" {
		cmd = append(cmd, "-U", shell.Quote(conf.PostgresOptions.Username))
	}

	if conf.PostgresOptions.Password != "" {
		connection.Environment.Set("PGPASSWORD", conf.PostgresOptions.Password)
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
	conn := conf.connection

	containerId := conn.DockerGetContainerId()
	Logger.Item("using docker container \"%s\"", containerId)

	containerEnv := conn.DockerGetEnvironment()

	// try to guess user/password
	if conf.PostgresOptions.Username == "" {
		// get superuser pass from env
		if pass, ok := containerEnv["POSTGRES_PASSWORD"]; ok {
			if user, ok := containerEnv["POSTGRES_USER"]; ok {
				Logger.Item("using postgres superadmin account \"%s\" (from env:POSTGRES_USER and env:POSTGRES_PASSWORD)", user)
				conf.PostgresOptions.Username = user
				conf.PostgresOptions.Password = pass
			} else {
				Logger.Item("using postgres superadmin account \"postgres\" (from env:POSTGRES_PASSWORD)")
				// only password available
				conf.PostgresOptions.Username = "postgres"
				conf.PostgresOptions.Password = pass
			}
		}
	}

}
