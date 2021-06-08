package testutil

import (
	"log"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)

func InitTestConf() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	viper.SetConfigFile("../test_util/resources/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		glog.Fatalf(err.Error())
	}

	os.Setenv("MongoPassword", "test123")
}
