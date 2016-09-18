package router

import "github.com/julienschmidt/httprouter"

func newRouter() *httprouter.Router {
	router := httprouter.New()
	router.RedirectFixedPath = false
	router.RedirectTrailingSlash = false

	return
}
