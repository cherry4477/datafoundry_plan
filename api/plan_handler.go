package api

import (
	"crypto/rand"
	"fmt"
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

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	username, e := validateAuth(r.Header.Get("Authorization"))
	if e != nil {
		JsonResult(w, http.StatusUnauthorized, e, nil)
		return
	}
	logger.Info("username:%v", username)

	if !canEditSaasApps(username) {
		JsonResult(w, http.StatusUnauthorized, GetError(ErrorCodePermissionDenied), nil)
		return
	}

	plan := &models.Plan{}
	err := common.ParseRequestJsonInto(r, plan)
	if err != nil {
		logger.Error("Parse body err: %v", err)
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeParseJsonFailed, err.Error()), nil)
		return
	}

	plan.Plan_id = genUUID()

	logger.Info("plan: %v", plan)

	//create plan in database
	planId, err := models.CreatePlan(db, plan)
	if err != nil {
		logger.Error("Create plan err: %v", err)
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeCreatePlan, err.Error()), nil)
		return
	}

	logger.Info("End create plan handler.")
	JsonResult(w, http.StatusOK, nil, planId)
}

func DeletePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: DELETE %v.", r.URL)

	logger.Info("Begin delete plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	planId := params.ByName("id")
	logger.Debug("Plan id: %s.", planId)

	// /delete in database
	err := models.DeletePlan(db, planId)
	if err != nil {
		logger.Error("Delete plan err: %v", err)
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeDeletePlan, err.Error()), nil)
		return
	}

	logger.Info("End delete plan handler.")
	JsonResult(w, http.StatusOK, nil, nil)
}

func ModifyPlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: PUT %v.", r.URL)

	logger.Info("Begin modify plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	plan := &models.Plan{}
	err := common.ParseRequestJsonInto(r, plan)
	if err != nil {
		logger.Error("Parse body err: %v", err)
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeParseJsonFailed, err.Error()), nil)
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
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeModifyPlan, err.Error()), nil)
		return
	}

	logger.Info("End modify plan handler.")
	JsonResult(w, http.StatusOK, nil, nil)
}

func RetrievePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: GET %v.", r.URL)

	logger.Info("Begin retrieve plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	planId := params.ByName("id")
	plan, err := models.RetrievePlanByID(db, planId)
	if err != nil {
		logger.Error("Get plan err: %v", err)
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeGetPlan), nil)
		return
	}

	logger.Info("End retrieve plan handler.")
	JsonResult(w, http.StatusOK, nil, plan)
}

func QueryPlanList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: GET %v.", r.URL)

	logger.Info("Begin retrieve plan handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	r.ParseForm()

	region := r.Form.Get("region")
	ptype := r.Form.Get("type")
	belong := r.Form.Get("belong")

	offset, size := OptionalOffsetAndSize(r, 30, 1, 100)
	orderBy := models.ValidateOrderBy(r.Form.Get("orderby"))
	sortOrder := models.ValidateSortOrder(r.Form.Get("sortorder"), false)

	count, apps, err := models.QueryPlans(db, region, ptype, belong, orderBy, sortOrder, offset, size)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeQueryPlans, err.Error()), nil)
		return
	}

	logger.Info("End retrieve plan handler.")
	JsonResult(w, http.StatusOK, nil, NewQueryListResult(count, apps))
}

func RetrievePlanRegion(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: GET %v.", r.URL)

	logger.Info("Begin retrieve plans's region handler.")

	db := models.GetDB()
	if db == nil {
		logger.Warn("Get db is nil.")
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	regions, err := models.RetrievePlanRegion(db)
	if err != nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeGetPlansRegion), nil)
	}

	logger.Info("End retrieve plans's region handler.")
	JsonResult(w, http.StatusOK, nil, NewQueryListResult(int64(len(regions)), regions))
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

func validateAuth(token string) (string, *Error) {
	if token == "" {
		return "", GetError(ErrorCodeAuthFailed)
	}

	username, err := getDFUserame(token)
	if err != nil {
		return "", GetError2(ErrorCodeAuthFailed, err.Error())
	}

	return username, nil
}

func canEditSaasApps(username string) bool {
	return username == "wangmeng5"
}
