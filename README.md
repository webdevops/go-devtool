# go-devtool

[![GitHub release](https://img.shields.io/github/release/webdevops/go-devtool.svg)](https://github.com/webdevops/go-devtool/releases)
[![license](https://img.shields.io/github/license/webdevops/go-devtool.svg)](https://github.com/webdevops/go-devtool/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-devtool.svg?branch=master)](https://travis-ci.org/webdevops/go-devtool)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-devtool/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-devtool/latest/total.svg)]()

Easy development tools for MySQL, PostgreSQL and TYPO3

Successor for [CliTools](https://github.com/webdevops/clitools) written on Golang

Features
========

MySQL:
- Connection to local, remote (`--ssh user@example.com`) and dockerized (`--docker containerid` or `--docker compose:mysql`) MySQL installations supported
- Show query log
- Show slow log
- Convert MySQL Schema and tables to specific charset and collation
- Backup MySQL Schema to file with automatic compression
- Restore MySQL Schema to file with automatic decompression

Postgres:
- Connection to local, remote (`--ssh user@example.com`) and dockerized (`--docker containerid` or `--docker compose:postgres`) PostgreSQL installations supported
- Backup PostgreSQL Schema to file with automatic compression
- Restore PostgreSQL Schema to file with automatic decompression

TYPO3:
- Create file stubs from existing FAL informations
- Create TYPO3 backend user

File:
- Create file stubs (small file examples) based on filelist

Install
=======

The binary file can be found in the [project releases](https://github.com/webdevops/go-devtool/releases).

```
DOWNLOAD_VERSION=0.3.3
DOWNLOAD_OS=linux
DOWNLOAD_ARCH=x64

wget -O/usr/local/bin/gdt "https://github.com/webdevops/go-devtool/releases/download/${DOWNLOAD_VERSION}/gdt-${DOWNLOAD_OS}-${DOWNLOAD_ARCH}"
chmod +x /usr/local/bin/gdt
```

Help
====

```
Usage:
  godevtool [OPTIONS] <command>

Application Options:
  -v, --verbose  verbose mode

Help Options:
  -h, --help     Show this help message

Available commands:
  file:stubs               Create file stubs from list of files
  mysql:convert            MySQL convert database charset/collation
  mysql:debug              MySQL debug
  mysql:dump               MySQL dump database
  mysql:restore            MySQL restore database
  mysql:server:dump        MySQL dump instance
  mysql:server:restore     MySQL restore instance
  mysql:slowlog            MySQL slow query log
  mysql:sql                MySQL shell
  postgres:dump            PostgreSQL dump database
  postgres:restore         PostgreSQL restore database
  postgres:server:dump     PostgreSQL dump server
  postgres:server:restore  PostgreSQL restore server
  postgres:sql             PostgreSQL shell
  self-update              Self update
  typo3:beuser             TYPO3 create BE user
  typo3:stubs              TYPO3 create file stubs
  version                  Show version

```

Docker support
==============

Using the parameter ``--docker=configuration`` this commands can be
execued with docker containers. If the container id is passed the
container is used without lookup using eg. `docker-compose`.

**docker-compose:**

*CONTAINER* is the name of the docker-compose container.

| DSN style configuration                                             | Description                                                                                     |
|:--------------------------------------------------------------------|:------------------------------------------------------------------------------------------------|
| ``compose:CONTAINER``                                               | Use container with docker-compose in current directory                                          |
| ``compose:CONTAINER;path=/path/to/project``                         | Use container with docker-compose in `/path/to/project` directory                               |
| ``compose:CONTAINER;path=/path/to/project;file=custom-compose.yml`` | Use container with docker-compose in `/path/to/project` directory and `custom-compose.yml` file |
| ``compose:CONTAINER;project-name=foobar``                           | Use container with docker-compose in current directory with project name `foobar`               |
| ``compose:CONTAINER;host=example.com``                              | Use container with docker-compose in current directory with docker host `example.com`           |
| ``compose:CONTAINER;env[FOOBAR]=BARFOO``                            | Use container with docker-compose in current directory with env var `FOOBAR` set to `BARFOO`    |

| Query style configuration                                             | Description                                                                                     |
|:----------------------------------------------------------------------|:------------------------------------------------------------------------------------------------|
| ``compose://CONTAINER``                                               | Use container with docker-compose in current directory                                          |
| ``compose://CONTAINER?path=/path/to/project``                         | Use container with docker-compose in `/path/to/project` directory                               |
| ``compose://CONTAINER?path=/path/to/project&file=custom-compose.yml`` | Use container with docker-compose in `/path/to/project` directory and `custom-compose.yml` file |
| ``compose://CONTAINER?project-name=foobar``                           | Use container with docker-compose in current directory with project name `foobar`               |
| ``compose://CONTAINER?host=example.com``                              | Use container with docker-compose in current directory with docker host `example.com`           |
| ``compose://CONTAINER?env[FOOBAR]=BARFOO``                            | Use container with docker-compose in current directory with env var `FOOBAR` set to `BARFOO`    |

Examples
========

MySQL commands
--------------

```bash
# Dump db1 into db1.sql.gz using local MySQL with user root and password dev
gdt mysql:schema:dump --mysql.user root --mysql.password dev db1 db1.sql.gz
gdt mysql:schema:dump --mysql mysql://root:dev@localhost db1 db1.sql.gz

# Dump db1 into db1.sql.gz using remote MySQL on host example.com with user root and password dev
gdt mysql:schema:dump --hostname example.com --mysql.user root --mysql.password dev db1 db1.sql.gz

# Dump db1 into db1.sql.gz using remote MySQL with user root and password dev on host example.com using SSH with user foobar 
gdt mysql:schema:dump --ssh foobar@example.com -u root -p dev db1 db1.sql.gz

# Dump db1 into db1.sql.gz using docker-compose container mysql
gdt mysql:schema:dump --docker compose:mysql db1 db1.sql.gz

# Dump db1 into db1.sql.gz using docker container 081e7bfaada1
gdt mysql:schema:dump --docker 081e7bfaada1 db1 db1.sql.gz

# Restore db1 from db1.sql.gz using docker container 081e7bfaada1
gdt mysql:schema:restore --docker 081e7bfaada1 db1 db1.sql.gz

```

PostgreSQL commands
-------------------

(same as mysql)

```bash
# Dump db1 into db1.sql.gz using docker container 081e7bfaada1
gdt postgres:schema:backup --docker=081e7bfaada1 db1 db1.sql.gz

# Restore db1 from db1.sql.gz using docker container 081e7bfaada1
gdt postgres:schema:restore --docker=081e7bfaada1 db1 db1.sql.gz

```

TYPO3 commands
--------------

```bash
# Create FAL stubs (example files) from existing TYPO3 database (Docker container is the MySQL container)
gdt typo3:stubs --docker=081e7bfaada1 typo3 /path/to/typo3/root/

# Inject BE user (user: dev, password: dev) into database
gdt typo3:beuser --docker=081e7bfaada1 typo3 

```


FILE commands
-------------

```bash
# Create stubs from stdin
cat filelist | gdt file:stubs --stdin --path /tmp/foobar/

# Create stubs from file content (from filelist)
gdt file:stubs --stdin --path /tmp/foobar/ filelist

```
