package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type HomeData struct {
	Targets  int
	Keywords int
	Results  int
	Matches  int
}

func (user *User) GetHomeData() (home HomeData, err error) {
	fmt.Println(Gray(8-1, "Starting GetHomeData..."))
	err = Db.QueryRow(`
                        WITH 
                            results_last_days AS(
                                SELECT
                                    ut.userid,
                                    COUNT(DISTINCT r.id) AS count_results
                                FROM userstargets ut
                                LEFT JOIN targets t ON(ut.targetid = t.id)
                                LEFT JOIN scrapers s ON(t.id = s.targetid)
                                LEFT JOIN results r ON(r.scraperid = s.id)
                                WHERE ut.userid=$1
                                AND r.createdat > current_date - interval '7' day
                                GROUP BY 1),
                            matches_last_days AS(
                                SELECT
                                    utk.userid,
                                    COUNT(DISTINCT m.id) AS count_matches
                                FROM userstargetskeywords utk
                                LEFT JOIN matches m ON(utk.keywordid = m.keywordid)
                                WHERE utk.userid=$1
                                AND m.createdat > current_date - interval '7' day
                                GROUP BY 1)
                        SELECT
                            rld.count_results,
                            mld.count_matches,
                            COUNT(DISTINCT utk.targetid),
                            COUNT(DISTINCT utk.keywordid)
                        FROM userstargetskeywords utk
                        LEFT JOIN results_last_days rld ON(utk.userid = rld.userid)
                        LEFT JOIN matches_last_days mld ON(utk.userid = mld.userid)
                        WHERE utk.userid=$1
                        AND utk.deletedat IS NULL
                        GROUP BY 1, 2;`, user.Id).
		Scan(
			&home.Results,
			&home.Matches,
			&home.Targets,
			&home.Keywords,
		)
	return
}
