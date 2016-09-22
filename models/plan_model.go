package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Plan struct {
	Plan_id        int `json:"plan_id, omitempty"`
	Plan_number    string
	Plan_type      string
	Specification1 string
	Specification2 string
	Price          float32
	Cycle          string
	Create_time    time.Time
	Status         string
}

func CreatePlan(db *sql.DB, planInfo *Plan) (int64, error) {
	logger.Info("Model create a plan.")

	nowstr := time.Now().Format("2006-01-02 15:04:05.999999")
	sqlstr := fmt.Sprintf(`insert into DF_PLAN (
				PLAN_NUMBER, PLAN_TYPE, SPECIFICATION1, SPECIFICATION2,
				PRICE, CYCLE, CREATE_TIME, STATUS
				) values (
				?, ?, ?, ?, ?, ?,
				'%s', '%s')`,
		nowstr, "A")

	result, err := db.Exec(sqlstr,
		planInfo.Plan_number, planInfo.Plan_type, planInfo.Specification1, planInfo.Specification2,
		planInfo.Price, planInfo.Cycle)

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}
