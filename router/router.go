package router

import (
	"github.com/asiainfoLDP/datafoundry_plan/api"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	Platform_Local  = "local"
	Platform_DataOS = "dataos"
)

var (
	Platform = Platform_DataOS
	logger   = log.GetLogger()
)

//==============================================================
//
//==============================================================

func handler_Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	api.JsonResult(w, http.StatusNotFound, api.GetError(api.ErrorCodeUrlNotSupported), nil)
}

func httpNotFound(w http.ResponseWriter, r *http.Request) {
	api.JsonResult(w, http.StatusNotFound, api.GetError(api.ErrorCodeUrlNotSupported), nil)
}

type HttpHandler struct {
	handler http.HandlerFunc
}

func (httpHandler *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if httpHandler.handler != nil {
		httpHandler.handler(w, r)
	}
}

//==============================================================
//
//==============================================================

func InitRouter() *httprouter.Router {
	router := httprouter.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false

	router.POST("/", handler_Index)
	router.DELETE("/", handler_Index)
	router.PUT("/", handler_Index)
	router.GET("/", handler_Index)

	router.NotFound = &HttpHandler{httpNotFound}
	router.MethodNotAllowed = &HttpHandler{httpNotFound}

	return router
}

func NewRouter(router *httprouter.Router) {
	logger.Info("new router.")
	//router.POST("/saasappapi/v1/apps", api.TimeoutHandle(500*time.Millisecond, CreateApp))
	//router.DELETE("/saasappapi/v1/apps/:id", api.TimeoutHandle(500*time.Millisecond, DeleteApp))
	//router.PUT("/saasappapi/v1/apps/:id", api.TimeoutHandle(500*time.Millisecond, ModifyApp))
	//router.GET("/saasappapi/v1/apps/:id", api.TimeoutHandle(500*time.Millisecond, RetrieveApp))
	//router.GET("/saasappapi/v1/apps", api.TimeoutHandle(500*time.Millisecond, QueryAppList))
}
