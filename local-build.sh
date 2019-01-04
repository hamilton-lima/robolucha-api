#!/bin/bash
export REDIS_HOST=localhost
export REDIS_PORT=6379
export INTERNAL_API_KEY=9239

export MYSQL_ROOT_PASSWORD=foo123123
export MYSQL_DATABASE=robolucha_db
export MYSQL_USER=robolucha_uzr
export MYSQL_PASSWORD=foo123123

go build -o $HOME/go/bin/robolucha-api
$HOME/go/bin/robolucha-api