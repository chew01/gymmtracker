package main

import (
	"errors"
	"flag"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"gymmtracker/data"
	"log"
	"os"
	"os/signal"
)

var stdout = flag.Bool("stdout", true, "enable/disable logging to console")
var logPath = flag.String("logfile", "logs/client.log", "file to output log to")

func main() {
	// Load environment variables to be used in all other functions
	err := godotenv.Load()
	if err != nil {
		err = errors.New("error loading environment variables: " + err.Error())
		panic(err)
	}

	// Parse cmd flags and instantiate new client
	flag.Parse()
	logFlags := log.Lshortfile | log.Ltime

	options := &data.Options{
		Stdout:   *stdout,
		LogPath:  *logPath,
		LogFlags: logFlags,
	}
	client := data.NewClient(options)

	// Add cron scheduling for every 8am, 11am, 2pm, 5pm, 8pm
	c := cron.New()
	err = c.AddFunc("0 0 8,11,14,17,20 * * *", client.Record)
	if err != nil {
		err = errors.New("error adding cronjob: " + err.Error())
	}

	go c.Start()

	// Add interrupt/kill condition for process to keep running
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig

}
