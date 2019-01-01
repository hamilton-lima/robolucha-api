# robolucha-api

## Generate API swagger file 

Follow the instalation procedures in https://github.com/swaggo/swag

```
code $HOME/.profile
export PATH=$PATH:$HOME/go/bin

go get -u github.com/swaggo/swag/cmd/swag
cd ~/Code/robolucha/robolucha-api
swag init

cd ~/go/src/gitlab.com/robolucha/robolucha-api
go get -v
go run main.go
open http://localhost:8080/public/swagger/index.html

```
- run swag init in the api folder
- use the generated file in api/docs/swagger/swagger.json to generate API clients

## Create users

Run script at database 
```
insert into users(email) values('foo@bar.com');
```

## Enable SQL log mode 

```
GORM_DEBUG=true
```
