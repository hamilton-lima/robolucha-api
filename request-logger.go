package main

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// From
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := ioutil.ReadAll(c.Request.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		log.WithFields(log.Fields{
			"body": readBody(rdr1),
		}).Debug("request body")

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}
