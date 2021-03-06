package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

func LogSetup() *os.File {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.MkdirAll("logs", 0755)
		if err != nil {
			panic(err)
		}
	}
	t := time.Now()
	nameLog := fmt.Sprintf("./logs/deployer-%d-%02d-%02d.log", t.Year(), t.Month(), t.Day())
	// Config Log path
	file, err := os.OpenFile(nameLog, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	return file
}

func LogClose(f *os.File) {
	f.Close()
}
