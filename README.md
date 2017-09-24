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


Help
====

```
Usage:
  main [OPTIONS] <mysql:convert | mysql:debug | mysql:slowlog>

Application Options:
  -v, --verbose      verbose mode
  -V, --version      show version and exit
      --dumpversion  show only version number and exit
      --help         show this help message

Help Options:
  -h, --help         Show this help message

Available commands:
  mysql:convert  MySQl convert schema charset/collation
  mysql:debug    MySQl debug
  mysql:slowlog  MySQl slow query log

```
