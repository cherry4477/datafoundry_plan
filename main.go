package main

import (
	"flag"
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/api"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"github.com/asiainfoLDP/datafoundry_plan/models"
	"github.com/asiainfoLDP/datafoundry_plan/router"
	"net/http"
)

const SERVERPORT = 8574

var (
	debug = flag.Bool("debug", false, "debug mode")
	local = flag.Bool("local", false, "running on local")

	logger = log.GetLogger()

	//init a router
	initRouter = router.InitRouter()
)

type Service struct {
	httpPort int
}

func newService(httpPort int) *Service {
	service := &Service{
		httpPort: httpPort,
	}

	return service
}

func main() {

	//new a router
	router.NewRouter(initRouter)

	// init db
	models.InitDB()

	service := newService(SERVERPORT)
	address := fmt.Sprintf(":%d", service.httpPort)
	logger.Debug("address: %v", address)

	logger.Info("Listening http at: %s", address)
	err := http.ListenAndServe(address, initRouter)
	if err != nil {
		logger.Error("http listen and server err: %v", err)
		return
	}

	return
}

func init() {
	//api.InitMQ()

	flag.Parse()
	log.SetDebug = *debug
	models.SetPlatform = *local

	//init log
	log.InitLog()

	//init remote
	api.InitGateWay()

}
