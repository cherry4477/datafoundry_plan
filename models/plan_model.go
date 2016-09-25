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

func DeletePlan(db *sql.DB, planId int) error {
	logger.Info("Model delete a plan.")

	sqlstr := fmt.Sprintf(`update DF_PLAN set status = "N" where PLAN_ID = %d`, planId)

	_, err := db.Exec(sqlstr)
	if err != nil {
		return err
	}

	return err
}

func ModifyPlan(db *sql.DB, planInfo *Plan) error {
	logger.Info("Model modify a plan.")

	plan, err := RetrievePlanByID(db, planInfo.Plan_id)
	if err != nil {
		return err
	}
	logger.Debug("Retrieve plan: %v", plan)

	err = modifyPlanStatusToN(db, plan.Plan_id)
	if err != nil {
		return err
	}

	planInfo.Plan_id = 0
	_, err = CreatePlan(db, planInfo)
	if err != nil {
		return err
	}

	return err
}

func RetrievePlanByID(db *sql.DB, planID int) (*Plan, error) {
	return getSinglePlan(db, fmt.Sprintf("where PLAN_ID = %d", planID))
}

func getSinglePlan(db *sql.DB, sqlWhere string) (*Plan, error) {
	apps, err := queryPlans(db, sqlWhere, 1, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if len(apps) == 0 {
		return nil, nil
	}

	return apps[0], nil
}

func queryPlans(db *sql.DB, sqlWhereAll string, limit int, offset int64, sqlParams ...interface{}) ([]*Plan, error) {
	offset_str := ""
	if offset > 0 {
		offset_str = fmt.Sprintf("offset %d", offset)
	}
	sql_str := fmt.Sprintf(`select
					PLAN_ID,
					PLAN_NUMBER, PLAN_TYPE,
					SPECIFICATION1,
					SPECIFICATION2,
					PRICE, CYCLE,
					CREATE_TIME, STATUS
					from DF_PLAN
					%s
					limit %d
					%s
					`,
		sqlWhereAll,
		limit,
		offset_str)
	rows, err := db.Query(sql_str, sqlParams...)

	fmt.Println(">>> ", sql_str)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]*Plan, 0, 100)
	for rows.Next() {
		plan := &Plan{}
		err := rows.Scan(
			&plan.Plan_id,
			&plan.Plan_number, &plan.Plan_type, &plan.Specification1, &plan.Specification2,
			&plan.Price, &plan.Cycle, &plan.Create_time, &plan.Status,
		)
		if err != nil {
			return nil, err
		}
		//validateApp(s) // already done in scanAppWithRows
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func modifyPlanStatusToN(db *sql.DB, planId int) error {
	sqlstr := fmt.Sprintf(`update DF_PLAN set status = "N" where PLAN_ID = %d`, planId)

	_, err := db.Exec(sqlstr)
	if err != nil {
		return err
	}

	return err
}
