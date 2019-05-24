set -o xtrace

go get -v github.com/gin-contrib/cors
go get -v github.com/gin-gonic/gin
go get -v github.com/go-sql-driver/mysql
go get -v github.com/gomodule/redigo/redis
go get -v github.com/jinzhu/copier
go get -v github.com/jinzhu/gorm
go get -v github.com/jinzhu/gorm/dialects/mysql
go get -v github.com/jinzhu/gorm/dialects/sqlite
go get -v github.com/sirupsen/logrus
go get -v github.com/stretchr/testify/assert
go get -v github.com/swaggo/gin-swagger
go get -v github.com/swaggo/gin-swagger/swaggerFiles
go get -v github.com/dgrijalva/jwt-go
go get -v github.com/alecthomas/template

go get -v gitlab.com/robolucha/robolucha-api/auth
go get -v gitlab.com/robolucha/robolucha-api/docs
go get -v gitlab.com/robolucha/robolucha-api/test

go get -v gopkg.in/matryer/try.v1
