package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type Row struct {
	Date  string
	Count int
}

func JobOffersPerDayPerTarget() (jobs []Row) {
	fmt.Println(Gray(8-1, "Starting JobOffersPerDayPerTarget..."))
	rows, err := Db.Query(`SELECT
                                updatedat,
                                countJobs
                            FROM (
                                SELECT
                                    TO_CHAR(r.updatedat, 'YYYY-MM-DD') AS updatedat,
                                    COUNT(DISTINCT r.id) AS countJobs,
                                    ROW_NUMBER() OVER () AS rn
                                FROM results r
                                LEFT JOIN scrapers s ON(r.scraperid = s.id)
                                WHERE s.name = 'Zalando'
                                GROUP BY 1
                                ORDER BY 1 DESC
                            ) as t
                            where rn != 1;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		row := Row{}
		if err = rows.Scan(
			&row.Date,
			&row.Count); err != nil {
			panic(err.Error())
		}
		jobs = append(jobs, row)
	}
	rows.Close()
	return
}
