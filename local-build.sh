#!/bin/bash
export REDIS_HOST=localhost
export REDIS_PORT=6379

go build -o $HOME/go/bin/robolucha-api
$HOME/go/bin/robolucha-api