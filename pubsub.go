package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func connect() redis.Conn {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	serverAddr := fmt.Sprintf("%v:%v", host, port)

	log.WithFields(log.Fields{
		"serverAddr": serverAddr,
	}).Info("Connecting to REDIS")

	readTimeout := time.Minute + (10 * time.Second)
	writeTimeout := 10 * time.Second

	var connection redis.Conn

	connection, err := redis.Dial("tcp", serverAddr,
		redis.DialReadTimeout(readTimeout),
		redis.DialWriteTimeout(writeTimeout))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Info("Error connecting to REDIS")
		return nil
	}

	return connection
}

// Publish message to REDIS
func Publish(channel string, message string) {
	conn := connect()
	_, err := conn.Do("PUBLISH", channel, message)

	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"channel": channel,
			"message": message,
		}).Error("Error publishing message to REDIS")
	}
}
