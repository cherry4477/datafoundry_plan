package main

import (
	"flag"
	"github.com/asiainfoLDP/datafoundry_plan/initialize"
	"github.com/astaxie/beego/logs"
	"fmt"
)

var (
	debug = flag.Bool("debug", false, "is debug mode?")
	log   *logs.BeeLogger
)

func main() {

	log.Debug("debug....")
	log.Info("info....")
	log.Notice("notice...")
	log.Warn("warn...")
	log.Error("error...")
	log.Critical("critical...")
	//log.Alert("alert...")
	//log.Emergency("emergency...")

}

func init() {
	flag.Parse()
	initialize.Debug = *debug
	fmt.Println(*debug)
	fmt.Println(initialize.Debug)

	log = initialize.GetLog()
}
