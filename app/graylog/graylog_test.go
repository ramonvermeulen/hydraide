package graylog

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	g := New("xxx.xxx.xxx.xxx:5140", "GraylogTest")
	g.SetLogLevel("trace")

	for i := 0; i < 500; i++ {

		log.WithFields(log.Fields{
			"testKey": "test value",
		}).Info("anakin ", i)

		time.Sleep(1 * time.Second)

	}

}
