package api

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	//"github.com/asiainfoLDP/datahub_commons/common"

	"github.com/asiainfoLDP/datafoundry_appmarket/market"
)

//==================================================================
//
//==================================================================

func CreateApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	appId := 1

	JsonResult(w, http.StatusOK, nil, appId)
}

func DeleteApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	JsonResult(w, http.StatusOK, nil, nil)
}

func ModifyApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	app := &market.SaasApp{
		App_id:      123,
		Provider:    "ABCD ltd.",
		Name:        "Quick mail",
		Version:     "1.0.0",
		Category:    "email",
		Description: "cool SaaS app",
		Icon_url:    "/components/header/img/logo.png",
		Create_time: time.Now(),
	}

	JsonResult(w, http.StatusOK, nil, app)
}

func RetrieveApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	app := &market.SaasApp{
		App_id:      123,
		Provider:    "ABCD ltd.",
		Name:        "Quick mail",
		Version:     "1.0.0",
		Category:    "email",
		Description: "cool SaaS app",
		Icon_url:    "/components/header/img/logo.png",
		Create_time: time.Now(),
	}

	JsonResult(w, http.StatusOK, nil, app)
}

func QueryAppList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	apps := []*market.SaasApp{
		&market.SaasApp{
			App_id:      123,
			Provider:    "ABCD ltd.",
			Name:        "Quick mail",
			Version:     "1.0.0",
			Category:    "email",
			Description: "cool SaaS app",
			Icon_url:    "/components/header/img/logo.png",
			Create_time: time.Now(),
		},
		&market.SaasApp{
			App_id:      789,
			Provider:    "WXYZ ltd.",
			Name:        "net disk",
			Version:     "2.0.0",
			Category:    "storage",
			Description: "reliable storage app",
			Icon_url:    "/components/header/img/logo.png",
			Create_time: time.Now(),
		},
	}

	JsonResult(w, http.StatusOK, nil, newQueryListResult(int64(len(apps)), apps))
}
