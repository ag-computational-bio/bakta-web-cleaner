package main

import (
	"log"

	"github.com/ag-computational-bio/bakta-web-cleaner/cleaner"
	"github.com/golang/glog"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
)

var opts struct {
	ConfigFile string `short:"c" long:"configfile" description:"File of the config file" default:"./config/config.yaml"`
}

//Version Version tag
var Version string

func main() {
	// Enable line numbers in logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	_, err := flags.Parse(&opts)
	if err != nil {
		glog.Fatalf(err.Error())
	}

	viper.SetConfigFile(opts.ConfigFile)
	err = viper.ReadInConfig()
	if err != nil {
		glog.Fatalf(err.Error())
	}

	cleaner, err := cleaner.Init()
	if err != nil {
		glog.Fatalf(err.Error())
	}

	err = cleaner.RemoveExpired()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
