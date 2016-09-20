package handler

import (
	"encoding/json"
	"github.com/asiainfoLDP/datafoundry_plan/api"
	"github.com/asiainfoLDP/datafoundry_plan/log"
	"github.com/asiainfoLDP/datafoundry_plan/models"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var logger = log.GetLogger()

func CreatePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: POST %v.", r.URL)

	logger.Info("Begin create plan handler.")
	defer logger.Info("End create plan handler.")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("IOioutil err, %v.", err)
		return
	}

	plan := models.Plan{}
	err = json.Unmarshal(body, &plan)
	if err != nil {
		logger.Error("Unmarshal err: %v.", err)
		return
		//api.JsonResult(w, )
	}
	plan.Create_time = time.Now()
	plan.Status = "A"
	logger.Debug("Plan: %v", plan)

	planId := 1111

	//todo create in database

	api.JsonResult(w, http.StatusOK, nil, planId)
}

func DeletePlan(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Info("Request url: DELETE %v.", r.URL)

	logger.Info("Begin delete plan handler.")
	defer logger.Info("End delete plan handler.")

	planId := params.ByName("id")
	logger.Debug("Plan id: %s.", planId)

	//todo delete in database

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
