#!/bin/bash
export REDIS_HOST=localhost
export REDIS_PORT=6379
export INTERNAL_API_KEY=9239

export MYSQL_ROOT_PASSWORD=foo123123
export MYSQL_DATABASE=robolucha_db
export MYSQL_USER=robolucha_uzr
export MYSQL_PASSWORD=foo123123
export MYSQL_HOST=localhost

export API_PORT=8080

export GIN_MODE=debug
export GORM_DEBUG=true
export API_SECRET=343434343

echo "--- BUILD api"
go build -v -o $HOME/go/bin/robolucha-api

echo "--- START api"
$HOME/go/bin/robolucha-api
