package main

import (
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/api"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"github.com/asiainfoLDP/datafoundry_plan/router"
	"github.com/asiainfoLDP/datahub_commons/httputil"
	"net/http"
	"time"
)

const SERVERPORT = 8574

var (
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

	//todo init db

	service := newService(SERVERPORT)
	address := fmt.Sprintf(":%d", service.httpPort)
	logger.Debug("address: %v", address)

	err := http.ListenAndServe(address, httputil.TimeoutHandler(initRouter, 250*time.Millisecond, ""))
	if err != nil {
		logger.Error("http listen and server err: %v", err)
		return
	}
	logger.Info("Listening http at:%s", address)

	return
}

func init() {
	api.InitMQ()
}
