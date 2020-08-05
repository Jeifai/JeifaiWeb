package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type Row struct {
	Date         string
	CountCreated int
	CountClosed  int
	CountTotal   int
}

func JobsPerDayPerTarget(target string) (jobs []Row) {
	fmt.Println(Gray(8-1, "Starting JobsPerDayPerTarget..."))
	rows, err := Db.Query(`WITH view_ready AS (
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
                                                FROM generate_series((SELECT s.createdat FROM scrapers s WHERE s.name = $1), current_timestamp, '1 day'::interval) dd)
                                    SELECT
                                        cd.consdate AS createdat,
                                        CASE WHEN jcr.countCreated IS NULL THEN 0 ELSE jcr.countCreated END AS countCreated,
                                        CASE WHEN jcl.countClosed IS NULL THEN 0 ELSE jcl.countClosed END AS countClosed
                                    FROM consecutive_dates cd
                                    LEFT JOIN jobs_created jcr ON(cd.consdate = jcr.createdat)
                                    LEFT JOIN jobs_closed jcl ON(jcr.createdat = jcl.closedat)) AS t)
                            SELECT
                                createdat,
                                countCreated,
                                countClosed,
                                countTotal
                            FROM view_ready
                            WHERE rn != 1;`, target)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		row := Row{}
		if err = rows.Scan(
			&row.Date,
			&row.CountCreated,
			&row.CountClosed,
			&row.CountTotal); err != nil {
			panic(err.Error())
		}
		jobs = append(jobs, row)
	}
	rows.Close()
	return
}
