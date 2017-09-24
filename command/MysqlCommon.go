package command

import (
	"bufio"
	"strings"
	"github.com/webdevops/go-shell"
	"github.com/webdevops/go-shell/commandbuilder"
)

type MysqlCommonOptions struct {
	Hostname string `long:"hostname"`
	Username string `short:"u" long:"user"`
	Password string `short:"p" long:"password"`
	Docker   string `          long:"docker"`
	SSH      string `          long:"ssh"`

	connection commandbuilder.Connection
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

func  (conf *MysqlCommonOptions) ExecMySqlStatement(statement string) {
	args := []string{"-N", "-B"}

	if conf.Hostname != "" {
		args = append(args, shell.Quote("-h" + conf.Hostname))
	}

	if conf.Username != "" {
		args = append(args, shell.Quote("-u" + conf.Username))
	}

	if conf.Password != "" {
		args = append(args, shell.Quote("-p" + conf.Password))
	}

	args = append(args, "-e", shell.Quote(statement))

	cmd := shell.Cmd(conf.connection.CommandBuilder("mysql", args...)...)
	cmd.Run()
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
