package main

import (
	"encoding/json"
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type TargetJobsTrend struct {
	CountTotalMinY      int
	CountCreated        json.RawMessage
	CountClosed         json.RawMessage
	CountTotal          json.RawMessage
}

type CompanyData struct {
	Employees           int
	Industry            string
	Companysize         string
	Headquarters        string
}

type TargetEmployeesTrend struct {
    CountEmployeesMinY      int
	CountEmployees          json.RawMessage
}

func (target *Target) GetTargetJobsTrend() (targetJobsTrend TargetJobsTrend) {
	fmt.Println(Gray(8-1, "Starting JobsPerDayPerTarget..."))
	err := Db.QueryRow(`WITH table_ready AS (
                            WITH view_ready AS (
                                SELECT
                                    t.createdat,
                                    t.countCreated,
                                    t.countClosed,
                                    sum(t.countCreated - t.countClosed) over (ORDER BY t.createdat) AS countTotal,
                                    ROW_NUMBER() OVER () AS rn
                                FROM (
                                    WITH
                                        jobs_created AS(
                                            SELECT
                                                r.createdat::date AS createdat,
                                                COUNT(DISTINCT r.id) AS countCreated
                                            FROM results r
                                            LEFT JOIN scrapers s ON(r.scraperid = s.id)
                                            WHERE s.name = $1
                                            GROUP BY 1),
                                        jobs_closed AS(
                                            SELECT
                                                r.updatedat::date AS closedat,
                                                COUNT(DISTINCT r.id) AS countClosed
                                            FROM results r
                                            LEFT JOIN scrapers s ON(r.scraperid = s.id)
                                            WHERE s.name = $1
                                            GROUP BY 1),
                                        consecutive_dates AS(
                                                SELECT
                                                    date_trunc('day', dd)::date AS consdate
                                                FROM generate_series(
                                                    (SELECT MIN(s.createdat) FROM scrapers s WHERE s.name = $1),
                                                    (SELECT MAX(r.updatedat) - INTERVAL '1 DAY' FROM scrapers s LEFT JOIN results r ON(s.id = r.scraperid) WHERE s.name = $1), 
                                                    '1 day'::interval) dd)
                                    SELECT
                                        cd.consdate AS createdat,
                                        CASE WHEN jcr.countCreated IS NULL THEN 0 ELSE jcr.countCreated END AS countCreated,
                                        CASE WHEN jcl.countClosed IS NULL THEN 0 ELSE jcl.countClosed END AS countClosed
                                    FROM consecutive_dates cd
                                    LEFT JOIN jobs_created jcr ON(cd.consdate = jcr.createdat)
                                    LEFT JOIN jobs_closed jcl ON(cd.consdate = jcl.closedat)) AS t)
                            SELECT
                                createdat,
                                countCreated,
                                countClosed,
                                countTotal
                            FROM view_ready
                            WHERE rn != 1)
                        SELECT
                            ROUND((MIN(CASE WHEN countTotal > 0 THEN countTotal END) * 0.96)::numeric::integer, -1),
                            json_object(array_agg(t.createdat::text), array_agg(t.countCreated::text)),
                            json_object(array_agg(t.createdat::text), array_agg(t.countClosed::text)),
                            json_object(array_agg(t.createdat::text), array_agg(t.countTotal::text))
                        FROM table_ready t`, target.Name).
		Scan(
			&targetJobsTrend.CountTotalMinY,
			&targetJobsTrend.CountCreated,
			&targetJobsTrend.CountClosed,
			&targetJobsTrend.CountTotal)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (target *Target) LinkedinDataPerTarget() (linkedinData CompanyData) {
	fmt.Println(Gray(8-1, "Starting LinkedinDataPerTarget..."))
	err := Db.QueryRow(`
        WITH latest_linkedin AS(
            SELECT
                l.targetid,
                MAX(l.id) AS latest_id
            FROM linkedin l
            LEFT JOIN targets t ON(l.targetid = t.id)
            WHERE t.id = $1
            GROUP BY 1)
        SELECT
            l.employees,
            l.industry,
            l.companysize,
            l.headquarters
        FROM latest_linkedin ll
        LEFT JOIN linkedin l ON(ll.latest_id = l.id)`, target.Id).
		Scan(
			&linkedinData.Employees,
			&linkedinData.Industry,
			&linkedinData.Companysize,
			&linkedinData.Headquarters,
		)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (target *Target) EmployeesTrendPerTarget() (targetEmployeesTrend TargetEmployeesTrend) {
	fmt.Println(Gray(8-1, "Starting EmployeesTrendPerTarget..."))
	err := Db.QueryRow(`
                WITH linkedin_data AS(
                    SELECT DISTINCT
                        l.createdat::date,
                        MAX(l.employees) AS count_employees
                    FROM linkedin l
                    WHERE targetid = $1
                    GROUP BY 1)
                SELECT
                    ROUND((MIN(CASE WHEN count_employees > 0 THEN count_employees END) * 0.96)::numeric::integer, -1),
                    json_object(array_agg(t.createdat::text), array_agg(t.count_employees::text))
                FROM linkedin_data t;`, target.Id).
                Scan(
                    &targetEmployeesTrend.CountEmployeesMinY,
                    &targetEmployeesTrend.CountEmployees,
                )
	if err != nil {
		panic(err.Error())
	}
	return
}
