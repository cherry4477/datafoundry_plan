package handler

import (
	"crypto/rand"
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/api"
	"github.com/asiainfoLDP/datafoundry_plan/common"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"github.com/asiainfoLDP/datafoundry_plan/models"
	"github.com/julienschmidt/httprouter"
	mathrand "math/rand"
	"net/http"
	"time"
)

var logger = log.GetLogger()

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func CreatePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: POST %v.", r.URL)

	logger.Info("Begin create plan handler.")
	defer logger.Info("End create plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeDbNotInitlized), nil)
		return
	}

	plan := &models.Plan{}
	err := common.ParseRequestJsonInto(r, plan)
	if err != nil {
		logger.Error("Parse body err: %v", err)
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeParseJsonFailed, err.Error()), nil)
		return
	}

	plan.Plan_id = genUUID()

	//create plan in database
	planId, err := models.CreatePlan(db, plan)
	if err != nil {
		logger.Error("Create plan err: %v", err)
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeCreatePlan, err.Error()), nil)
		return
	}

	api.JsonResult(w, http.StatusOK, nil, planId)
}

func DeletePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: DELETE %v.", r.URL)

	logger.Info("Begin delete plan handler.")
	defer logger.Info("End delete plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeDbNotInitlized), nil)
		return
	}

	planId := params.ByName("id")
	logger.Debug("Plan id: %s.", planId)

	// /delete in database
	err := models.DeletePlan(db, planId)
	if err != nil {
		logger.Error("Delete plan err: %v", err)
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeDeletePlan, err.Error()), nil)
		return
	}

	api.JsonResult(w, http.StatusOK, nil, nil)
}

func ModifyPlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: PUT %v.", r.URL)

	logger.Info("Begin modify plan handler.")
	defer logger.Info("End modify plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeDbNotInitlized), nil)
		return
	}

	plan := &models.Plan{}
	err := common.ParseRequestJsonInto(r, plan)
	if err != nil {
		logger.Error("Parse body err: %v", err)
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeParseJsonFailed, err.Error()), nil)
		return
	}
	logger.Debug("Plan: %v", plan)

	planId := params.ByName("id")
	logger.Debug("Plan id: %s.", planId)

	plan.Plan_id = planId

	//update in database
	err = models.ModifyPlan(db, plan)
	if err != nil {
		logger.Error("Modify plan err: %v", err)
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeModifyPlan, err.Error()), nil)
		return
	}

	api.JsonResult(w, http.StatusOK, nil, nil)
}

func RetrievePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: GET %v.", r.URL)

	logger.Info("Begin retrieve plan handler.")
	defer logger.Info("End retrieve plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeDbNotInitlized), nil)
		return
	}

	planId := params.ByName("id")
	plan, err := models.RetrievePlanByID(db, planId)
	if err != nil {
		logger.Error("Get plan err: %v", err)
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeGetPlan), nil)
		return
	}

	api.JsonResult(w, http.StatusOK, nil, plan)
}

func QueryPlanList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: GET %v.", r.URL)

	logger.Info("Begin retrieve plan handler.")
	defer logger.Info("End retrieve plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeDbNotInitlized), nil)
		return
	}

	r.ParseForm()
	//
	//provider, e := validateAppProvider(r.Form.Get("provider"), false)
	//if e != nil {
	//	api.JsonResult(w, http.StatusBadRequest, e, nil)
	//	return
	//}
	//
	//category, e := validateAppCategory(r.Form.Get("category"), false)
	//if e != nil {
	//	api.JsonResult(w, http.StatusBadRequest, e, nil)
	//	return
	//}

	offset, size := api.OptionalOffsetAndSize(r, 30, 1, 100)
	orderBy := models.ValidateOrderBy(r.Form.Get("orderby"))
	sortOrder := models.ValidateSortOrder(r.Form.Get("sortorder"), false)

	count, apps, err := models.QueryPlans(db, orderBy, sortOrder, offset, size)
	if err != nil {
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeQueryPlans, err.Error()), nil)
		return
	}

	api.JsonResult(w, http.StatusOK, nil, api.NewQueryListResult(count, apps))
}

func genUUID() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		logger.Warn("genUUID error: ", err.Error())

		mathrand.Read(bs)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
}
