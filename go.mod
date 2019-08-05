module gitlab.com/robolucha/robolucha-api

go 1.12

require (
	github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc
	github.com/bxcodec/faker v2.0.1+incompatible
	github.com/bxcodec/faker/v3 v3.1.0
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.9
	github.com/matryer/try v0.0.0-20161228173917-9ac251b645a2 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.3.0
	github.com/swaggo/gin-swagger v1.1.0
	github.com/swaggo/swag v1.5.1
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
	gotest.tools v2.2.0+incompatible
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
