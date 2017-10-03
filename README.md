# go-devtool

[![GitHub release](https://img.shields.io/github/release/webdevops/go-devtool.svg)](https://github.com/webdevops/go-devtool/releases)
[![license](https://img.shields.io/github/license/webdevops/go-devtool.svg)](https://github.com/webdevops/go-devtool/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/webdevops/go-devtool.svg?branch=master)](https://travis-ci.org/webdevops/go-devtool)
[![Github All Releases](https://img.shields.io/github/downloads/webdevops/go-devtool/total.svg)]()
[![Github Releases](https://img.shields.io/github/downloads/webdevops/go-devtool/latest/total.svg)]()

Easy development tools for MySQL

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

Help
====

```
Usage:
  main [OPTIONS] <command>

Application Options:
  -v, --verbose  verbose mode

Help Options:
  -h, --help     Show this help message

Available commands:
  mysql:backup      MySQL dump schema
  mysql:convert     MySQL convert schema charset/collation
  mysql:debug       MySQL debug
  mysql:restore     MySQL restore schema
  mysql:slowlog     MySQL slow query log
  postgres:backup   PostgreSQL dump schema
  postgres:restore  PostgreSQL restore schema
  typo3:stubs       TYPO3 create file stubs
  version           MySQL dump schema
```
