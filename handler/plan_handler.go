package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datafoundry_plan/api"
	"github.com/asiainfoLDP/datafoundry_plan/common"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"github.com/asiainfoLDP/datafoundry_plan/models"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	mathrand "math/rand"
	"net/http"
	"strconv"
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
		api.JsonResult(w, http.StatusInternalServerError, api.GetError(api.ErrorCodeDbNotInitlized), nil)
		return
	}

	plan := &models.Plan{}
	err := common.ParseRequestJsonInto(r, plan)
	if err != nil {
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeParseJsonFailed, err.Error()), nil)
		return
	}

	plan.Plan_number = genUUID()

	//create plan in database
	planId, err := models.CreatePlan(db, plan)
	if err != nil {
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

	id := params.ByName("id")
	logger.Debug("Plan id: %s.", id)

	planId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Strconv err: %v.", err)
		return
	}
	// /delete in database
	err = models.DeletePlan(db, planId)
	if err != nil {
		api.JsonResult(w, http.StatusBadRequest, api.GetError2(api.ErrorCodeDeletePlan, err.Error()), nil)
		return
	}

	api.JsonResult(w, http.StatusOK, nil, nil)
}

func ModifyPlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: PUT %v.", r.URL)

	logger.Info("Begin modify plan handler.")
	defer logger.Info("End modify plan handler.")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("IOioutil err, %v.", err)
	}

	plan := models.Plan{}
	err = json.Unmarshal(body, &plan)
	if err != nil {
		logger.Error("Unmarshal err: %v.", err)
		//api.JsonResult(w, )
	}
	logger.Debug("Plan: %v", plan)

	//todo update in database

	api.JsonResult(w, http.StatusOK, nil, nil)
}

func RetrievePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: GET %v.", r.URL)

	logger.Info("Begin retrieve plan handler.")
	defer logger.Info("End retrieve plan handler.")

	id := params.ByName("id")
	logger.Debug("Plan id: %s.", id)

	planId, err := strconv.Atoi(id)
	if err != nil {
		logger.Error("Strconv err: %v.", err)
		return
	}

	plan := models.Plan{
		Plan_id:        planId,
		Plan_number:    "1d3452ea-7f14-11e6-9fe0-2344dd5557c3",
		Plan_type:      "C",
		Specification1: "1Gi",
		Specification2: "8Gi",
		Price:          88.88,
		Cycle:          "M",
		Create_time:    time.Now(),
		Status:         "A",
	}

	api.JsonResult(w, http.StatusOK, nil, plan)
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
