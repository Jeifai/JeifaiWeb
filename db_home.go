package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type HomeData struct {
	UserName                     string
	CountKeywords                int
	CountTargets                 int
	CountTargetsKeywords         int
	CountMatchesLast7Days        int
	CountResultsLast7Days        int
	CountClosedPositionLast7Days int
	CountOpenPositions           int
}

func (user *User) GetHomeData() (home HomeData, err error) {
	fmt.Println(Gray(8-1, "Starting GetHomeData..."))
	err = Db.QueryRow(`
                        WITH
                            into_user_targets AS(
                                SELECT
                                    ut.userid,
                                    COUNT(DISTINCT ut.targetid) AS count_targets
                                FROM userstargets ut
                                WHERE ut.userid=$1
                                AND ut.deletedat IS NULL
                                GROUP BY 1),
                            info_user_keywords AS(
                                SELECT
                                    utk.userid,
                                    COUNT(DISTINCT utk.keywordid) AS count_keywords,
                                    COUNT(DISTINCT utk.id) AS count_targetskeywords
                                FROM userstargetskeywords utk
                                WHERE utk.userid=$1
                                AND utk.deletedat IS NULL
                                GROUP BY 1),
                            matches_last_7_days AS(
                                SELECT
                                    utk.userid,
                                    COUNT(DISTINCT m.id) AS count_matches_last_7_days
                                FROM userstargetskeywords utk
                                LEFT JOIN matches m ON(utk.keywordid = m.keywordid)
                                WHERE utk.userid=$1
                                AND m.createdat > current_date - interval '7' day
                                AND utk.deletedat IS NULL
                                GROUP BY 1),
                            info_targets AS(
                                WITH latest_scraping AS (
                                    SELECT
                                        s.scraperid,
                                        MAX(s.id) AS scrapingid
                                    FROM scrapings s
                                    WHERE s.createdat::date >= current_timestamp::date
                                    GROUP BY 1)
                                SELECT
                                    ut.userid,
                                    COUNT(DISTINCT r.id) AS count_open_positions
                                FROM userstargets ut
                                LEFT JOIN scrapers s ON(ut.targetid = s.targetid)
                                INNER JOIN latest_scraping ls ON(s.id = ls.scraperid)
                                LEFT JOIN results r ON(ls.scrapingid = r.scrapingid)
                                WHERE ut.userid=1
                                AND ut.deletedat IS NULL
                                GROUP BY 1),
                            closed_positions_last_7_days AS(
                                WITH latest_scraping AS (
                                    SELECT
                                        s.scraperid,
                                        MAX(s.id) AS scrapingid
                                    FROM scrapings s
                                    WHERE s.createdat::date >= current_timestamp::date
                                    GROUP BY 1)
                                SELECT
                                    ut.userid,
                                    COUNT(DISTINCT r.id) AS count_closed_positions_last_7_days
                                FROM userstargets ut
                                LEFT JOIN scrapers s ON(ut.targetid = s.targetid)
                                LEFT JOIN results r ON(s.id = r.scraperid)
                                LEFT JOIN latest_scraping ls ON(s.id = ls.scraperid)
                                WHERE ut.userid=$1
                                AND r.updatedat > current_date - interval '7' day
                                AND r.scrapingid != ls.scrapingid
                                AND ut.deletedat IS NULL
                                GROUP BY 1),
                            results_last_7_days AS(
                                SELECT
                                    ut.userid,
                                    COUNT(DISTINCT r.id) AS count_results_last_7_days
                                FROM userstargets ut
                                LEFT JOIN scrapers s ON(ut.targetid = s.targetid)
                                LEFT JOIN results r ON(r.scraperid = s.id)
                                WHERE ut.userid=$1
                                AND r.createdat > current_date - interval '7' day
                                AND ut.deletedat IS NULL
                                GROUP BY 1)
                        SELECT
                            CASE WHEN iuk.count_keywords IS NULL THEN 0 ELSE iuk.count_keywords END,
                            CASE WHEN iut.count_targets IS NULL THEN 0 ELSE iut.count_targets END,
                            CASE WHEN iuk.count_targetskeywords IS NULL THEN 0 ELSE iuk.count_targetskeywords END,
                            CASE WHEN mld.count_matches_last_7_days IS NULL THEN 0 ELSE mld.count_matches_last_7_days END,
                            CASE WHEN rld.count_results_last_7_days IS NULL THEN 0 ELSE rld.count_results_last_7_days END,
                            CASE WHEN cpld.count_closed_positions_last_7_days IS NULL THEN 0 ELSE cpld.count_closed_positions_last_7_days END,
                            CASE WHEN it.count_open_positions IS NULL THEN 0 ELSE it.count_open_positions END
                        FROM into_user_targets iut
                        LEFT JOIN info_user_keywords iuk ON(iut.userid = iuk.userid)
                        LEFT JOIN info_targets it ON(iuk.userid = it.userid)
                        LEFT JOIN results_last_7_days rld ON(iuk.userid = rld.userid)
                        LEFT JOIN matches_last_7_days mld ON(rld.userid = mld.userid)
                        LEFT JOIN closed_positions_last_7_days cpld ON(mld.userid = cpld.userid)`, user.Id).
		Scan(
			&home.CountKeywords,
			&home.CountTargets,
			&home.CountTargetsKeywords,
			&home.CountMatchesLast7Days,
			&home.CountResultsLast7Days,
			&home.CountClosedPositionLast7Days,
			&home.CountOpenPositions,
		)
	return
}
