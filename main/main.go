package main

import (
	"flag"
	"github.com/alecthomas/log4go"
	"github.com/robfig/config"
	"goReptile/slave"
)

func init() {
	c, _ := config.ReadDefault("config/config.ini")
	level, _ := c.Int("log4go", "level")
	log4go.AddFilter("stdout", log4go.Level(level), log4go.NewConsoleLogWriter())
}

func main() {
	var tag string
	var version int

	flag.StringVar(&tag, "tag", "", "Input tag")
	flag.IntVar(&version, "version", 0, "Input version")

	flag.Parse()

	slave.Run(tag, version)
}
