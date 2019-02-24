# robolucha-api

## Generate API swagger file 

Follow the instalation procedures in https://github.com/swaggo/swag

```
code ~/.profile
export PATH=$PATH:$HOME/go/bin

go get -u github.com/swaggo/swag/cmd/swag
cd ~/code/robolucha/robolucha-api
swag init

cd ~/go/src/gitlab.com/robolucha/robolucha-api
go get -v
go run main.go
open http://localhost:8080/public/swagger/index.html

```
- run swag init in the api folder
- use the generated file in api/docs/swagger/swagger.json to generate API clients

## Local environment setup

Create symbolic link from workspace to gopath
```
	export WIN=/mnt/c/Users/hamil/code
	mkdir -p $WIN/go/src/gitlab.com/robolucha
	ln -s $WIN/robolucha-api $WIN/go/src/gitlab.com/robolucha
	export PATH=$PATH:$WIN/go/bin
	cd $WIN/go/src/gitlab.com/robolucha/robolucha-api
	go get -v	
```

Edit .profile
```
	nano $HOME/.profile
	export WIN=/mnt/c/Users/hamil/code
	export PATH=$PATH:$WIN/go/bin
	export WIN=/mnt/c/Users/hamil/code
	export GOPATH=$WIN/go
	export DOCKER_HOST=tcp://0.0.0.0:2375
```

## Create users

Run script at database 
```
insert into users(email) values('hamilton.lima@gmail.com');
```

## Enable SQL log mode 

```
GORM_DEBUG=true
```
