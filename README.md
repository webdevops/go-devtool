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
- Show query log for local, remote and docker MySQL installations
- Show slow log for local, remote and docker MySQL installations
- Convert MySQL Schema and tables to specific charset and collation
- Backup MySQL Schema to file with automatic compression
- Restore MySQL Schema to file with automatic decompression

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
  mysql:backup   MySQl dump schema
  mysql:convert  MySQl convert schema charset/collation
  mysql:debug    MySQl debug
  mysql:restore  MySQl restore schema
  mysql:slowlog  MySQl slow query log
  version        MySQl dump schema

```
