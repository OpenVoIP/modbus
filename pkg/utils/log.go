package utils

import (
	"log"
	"os"
	"sync"
)

var logger *log.Logger
var once sync.Once

func GetLogger() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "app: ", log.LstdFlags)
	})
	return logger
}
