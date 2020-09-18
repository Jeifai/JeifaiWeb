package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

type Target struct {
	Id          int
	Url         string
	Host        string
	Name        string `validate:"required,max=30,min=3"`
	CreatedAt   time.Time
	CreatedDate string
}

type TargetInfo struct {
	Name               string
	CreatedDate        string
	LastExtractionDate string
	Employees          int
	JobsAll            int
	JobsNow            int
	Opened             int
	Closed             int
}

type Job struct {
	Id          string
	IsSaved     bool
	CreatedDate string
	TargetName  string
	KeywordText string
	Title       string
	Location    string
	Url         string
}

func (target *Target) InsertTarget() {
	fmt.Println(Gray(8-1, "Starting InsertTarget..."))
	statement := `INSERT INTO targets (name, createdat)
                  VALUES ($1, $2)
                  RETURNING id, name, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		target.Name,
		time.Now(),
	).Scan(
		&target.Id,
		&target.Name,
		&target.CreatedAt,
	)
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) InsertUserTarget(target Target) {
	fmt.Println(Gray(8-1, "Starting InsertUserTarget..."))

	statement := `INSERT INTO userstargets (userid, targetid, createdat)
                  VALUES ($1, $2, current_timestamp);`
	stmt, err := Db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow(user.Id, target.Id)
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) TargetsNamesByUser() (targetsNames []string) {
	fmt.Println(Gray(8-1, "Starting TargetsNamesByUser..."))

	err := Db.QueryRow(`
                SELECT
                  ARRAY_AGG(t.name)
                FROM users u
                INNER JOIN userstargets ut ON(u.id = ut.userid) 
                INNER JOIN targets t ON(ut.targetid = t.id)
                WHERE ut.deletedat IS NULL
                AND u.id=$1;`, user.Id).
		Scan(
			pq.Array(&targetsNames),
		)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (user *User) SelectTargetsByUser() (targets []Target) {
	fmt.Println(Gray(8-1, "Starting SelectTargetsByUser..."))
	rows, err := Db.Query(`
							SELECT
								t.name,
								TO_CHAR(MIN(ut.createdat::date), 'YYYY-MM-DD')
							FROM userstargets ut
							LEFT JOIN targets t ON(ut.targetid = t.id)
							WHERE ut.userid = $1
							AND ut.deletedat IS NULL
							GROUP BY 1;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(
			&target.Name,
			&target.CreatedDate); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

func SelectTargetsByAll() (targets []string) {
	fmt.Println(Gray(8-1, "Starting SelectTargetsByAll..."))
	rows, err := Db.Query(`
							SELECT
								DISTINCT t.name
							FROM targets t;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var target string
		if err = rows.Scan(
			&target); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

func (user *User) InfoUsersTargetsByUser() (targetsinfo []TargetInfo) {
	fmt.Println(Gray(8-1, "Starting InfoUsersTargetsByUser..."))
	rows, err := Db.Query(`
                            WITH
                                linkedindata AS (
                                    WITH latest_linkedin AS(
                                        SELECT
                                            l.targetid,
                                            MAX(l.id) AS latest_id
                                        FROM linkedin l
                                        GROUP BY 1)
                                    SELECT
                                        l.targetid,
                                        l.employees
                                    FROM linkedin l
                                    INNER JOIN latest_linkedin ll ON(l.id = ll.latest_id)),
                                latest_scraping AS (
                                    SELECT
                                        s.scraperid,
                                        MAX(s.createdat) AS createdat,
                                        MAX(s.id) AS scrapingid
                                    FROM scrapings s
                                    GROUP BY 1),
                                usertargets AS (
                                    SELECT
                                        s.id AS scraperid,
                                        t.name,
                                        TO_CHAR(t.createdat, 'YYYY-MM-DD') AS createdat,
                                        ld.employees
                                    FROM users u
                                    INNER JOIN userstargets ut ON(u.id = ut.userid) 
                                    INNER JOIN targets t ON(ut.targetid = t.id)
                                    LEFT JOIN scrapers s ON(t.id = s.targetid)
                                    LEFT JOIN linkedindata ld ON(t.id = ld.targetid)
                                    WHERE ut.deletedat IS NULL
                                    AND u.id=$1
                                    ORDER BY t.createdat DESC)
                            SELECT
                                ut.name,
                                ut.createdat,
                                CASE WHEN ut.employees IS NULL THEN 0 ELSE ut.employees END AS employees,
                                CASE WHEN ls.createdat IS NULL THEN '' ELSE TO_CHAR(MAX(ls.createdat), 'YYYY-MM-DD') END AS last_extraction,
                                COUNT(DISTINCT r.url) AS all_time_job,
                                SUM(CASE WHEN r.scrapingid = ls.scrapingid THEN 1 ELSE 0 END) AS actual_job_opens,
                                SUM(CASE WHEN (r.createdat > current_date - interval '7' day) THEN 1 ELSE 0 END) AS open_positions_last_7_days,
                                SUM(CASE WHEN (r.updatedat > current_date - interval '7' day AND r.updatedat < current_date - interval '1' day) THEN 1 ELSE 0 END) AS close_positions_last_7_days
                            FROM usertargets ut
                            LEFT JOIN results r ON(ut.scraperid = r.scraperid)
                            LEFT JOIN latest_scraping ls ON(r.scraperid = ls.scraperid)
                            GROUP BY 1, 2, 3, ls.createdat
                            ORDER BY 2;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		targetinfo := TargetInfo{}
		if err = rows.Scan(
			&targetinfo.Name,
			&targetinfo.CreatedDate,
			&targetinfo.Employees,
			&targetinfo.LastExtractionDate,
			&targetinfo.JobsAll,
			&targetinfo.JobsNow,
			&targetinfo.Opened,
			&targetinfo.Closed); err != nil {
			panic(err.Error())
		}
		targetsinfo = append(targetsinfo, targetinfo)
	}
	rows.Close()
	if err != nil {
		panic(err.Error())
	}
	return
}

func (target *Target) SelectTargetByName() {
	fmt.Println(Gray(8-1, "Starting SelectTargetByName..."))
	_ = Db.QueryRow(`SELECT
                         t.id
                       FROM targets t
                       WHERE t.name=$1`, target.Name).Scan(&target.Id)
}

func TargetsByNames(targetNames []string) (targets []Target) {
	fmt.Println(Gray(8-1, "Starting TargetsByNames..."))

	rows, err := Db.Query(`SELECT
                                t.id,
                                t.name
                            FROM targets t
                            WHERE t.name LIKE ANY($1)`, pq.Array(targetNames))
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(&target.Id, &target.Name); err != nil {
			panic(err.Error())
		}
		targets = append(targets, target)
	}
	rows.Close()

	return
}

func (user *User) SelectUserTargetByUserAndTarget(target Target) (userTargetId int) {
	fmt.Println(Gray(8-1, "Starting SelectUserTargetByUserAndTarget..."))
	_ = Db.QueryRow(`SELECT
                         ut.id
                       FROM userstargets ut
                       WHERE ut.userid = $1
                       AND ut.targetid = $2
                       AND ut.deletedat IS NULL;`, user.Id, target.Id).Scan(&userTargetId)
	return
}

func (user *User) SelectTargetsKeywordsByUser() (utks []map[string]interface{}) {
	fmt.Println(Gray(8-1, "Starting SelectTargetsKeywordsByUser..."))
	rows, err := Db.Query(`
							WITH
								complete AS (
									WITH
										utks AS (
											WITH
												userkeywords AS(
													SELECT
														uk.id,
														uk.keywordid
													FROM userskeywords uk
													WHERE uk.userid = $1
													AND uk.deletedat IS NULL),
												usertargets AS(
													SELECT
														ut.id,
														ut.targetid
													FROM userstargets ut
													WHERE ut.userid = $1
													AND ut.deletedat IS NULL)
											SELECT
												k.text AS keyword_text,
												t.name AS target_name
											FROM userstargetskeywords utk
											INNER JOIN userkeywords uk ON(utk.userkeywordid = uk.id)
											INNER JOIN usertargets ut ON(utk.usertargetid = ut.id)
											LEFT JOIN keywords k ON(uk.keywordid = k.id)
											LEFT JOIN targets t ON(ut.targetid = t.id)),
										pivot_by_keyword AS (
										    SELECT 
												keyword_text, 
												json_agg(target_name) AS target_name
											FROM utks
											GROUP BY 1),
										pivot_by_target AS (
										    SELECT 
												target_name, 
												json_agg(keyword_text) AS keyword_text
											FROM utks
											GROUP BY 1)
								SELECT
									json_object_agg(keyword_text, target_name) AS agg_data
								FROM pivot_by_keyword
								UNION ALL
								SELECT
									json_object_agg(target_name, keyword_text) AS agg_data
								FROM pivot_by_target)
							SELECT
								*
							FROM complete
							WHERE agg_data IS NOT NULL;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var temp_utk string
		if err = rows.Scan(
			&temp_utk); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		utk := map[string]interface{}{}
		if err := json.Unmarshal([]byte(temp_utk), &utk); err != nil {
			panic(err)
		}
		utks = append(utks, utk)
	}
	rows.Close()
	return
}

func (user *User) UpdateDeletedAtInUsersTargets(target Target) {
	fmt.Println(Gray(8-1, "Starting UpdateDeletedAtInUsersTargets..."))

	statement := `UPDATE userstargets
				  SET deletedat = current_timestamp
				  WHERE userid = $1
				  AND targetid = $2;`
	stmt, err := Db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow(user.Id, target.Id)
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) DeleteUserTargetsKeywordsByKeywords(keywords []string) {
	fmt.Println(Gray(8-1, "Starting DeleteUserTargetsKeywordsByKeywords..."))

	statement := `DELETE FROM userstargetskeywords
					WHERE userkeywordid IN(
						SELECT
							uk.id AS userkeywordid
						FROM userskeywords uk
						WHERE uk.userid = $1
						AND uk.keywordid = (
							SELECT
								k.id
							FROM keywords k
							WHERE k.text = ANY($2)));`
	stmt, err := Db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow(user.Id, pq.Array(keywords))
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) DeleteUserTargetsKeywordsByTargets(targets []string) {
	fmt.Println(Gray(8-1, "Starting DeleteUserTargetsKeywordsByTargets..."))

	statement := `DELETE FROM userstargetskeywords
					WHERE usertargetid IN(
						SELECT
							ut.id AS usertargetid
						FROM userstargets ut
						WHERE ut.userid = $1
						AND ut.targetid = (
							SELECT
								t.id
							FROM targets t
							WHERE t.name = ANY($2)));`
	stmt, err := Db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow(user.Id, pq.Array(targets))
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) InsertUserTargetsKeywords(keywords []string, targets []string) {
	fmt.Println(Gray(8-1, "Starting InsertUserTargetsKeywords..."))

	var k_ids []int
	k_rows, err := Db.Query(`
					SELECT
						uk.id
					FROM userskeywords uk
					LEFT JOIN keywords k ON(uk.keywordid = k.id)
					WHERE uk.userid = $1
					AND k.text = ANY($2);`, user.Id, pq.Array(keywords))
	if err != nil {
		panic(err.Error())
	}
	for k_rows.Next() {
		var k_id int
		err = k_rows.Scan(&k_id)
		if err != nil {
			panic(err.Error())
		}
		k_ids = append(k_ids, k_id)
	}
	k_rows.Close()

	var t_ids []int
	t_rows, err := Db.Query(`
					SELECT
						ut.id
					FROM userstargets ut
					LEFT JOIN targets t ON(ut.targetid = t.id)
					WHERE ut.userid = $1
					AND t.name = ANY($2);`, user.Id, pq.Array(targets))
	if err != nil {
		panic(err.Error())
	}
	for t_rows.Next() {
		var t_id int
		err = t_rows.Scan(&t_id)
		if err != nil {
			panic(err.Error())
		}
		t_ids = append(t_ids, t_id)
	}
	t_rows.Close()

	valueStrings := []string{}
	for _, k_elem := range k_ids {
		for _, t_elem := range t_ids {
			str1 := "(" + strconv.Itoa(k_elem)
			str2 := "," + strconv.Itoa(t_elem)
			str3 := ",current_timestamp)"
			str_n := str1 + str2 + str3
			valueStrings = append(valueStrings, str_n)
		}
	}
	smt := `INSERT INTO userstargetskeywords (userkeywordid, usertargetid, createdat) VALUES %s;`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err = Db.Exec(smt)
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) SelectJobsByTargetsAndKeywords(targets []Target, keywords []Keyword) (jobs []Job) {
	fmt.Println(Gray(8-1, "Starting SelectJobsByTargetsAndKeywords..."))
	rows, err := Db.Query(`
							WITH
								usertargets AS(
									SELECT
										DISTINCT s.id,
										s.name
									FROM userstargets ut
									LEFT JOIN targets t ON(ut.targetid = t.id)
									LEFT JOIN scrapers s ON(t.id = s.targetid)
									WHERE ut.userid = $1
									AND ut.deletedat IS NULL),
								userkeywords AS(
									SELECT
										k.text
									FROM userskeywords ut
									LEFT JOIN keywords k ON(ut.keywordid = k.id)
									WHERE ut.userid = $1
									AND ut.deletedat IS NULL),
								userfavouriteresults AS(
									SELECT
										ft.resultid
									FROM favouriteresults ft
									WHERE ft.userid = $1
									AND ft.deletedat IS NULL)
							SELECT
								r.id,
								CASE WHEN uft.resultid IS NULL THEN FALSE ELSE TRUE END,
								TO_CHAR(r.createdat, 'YYYY-MM-DD') AS createdat,
								ut.name,
								uk.text,
								r.title,
								CASE WHEN r.location IS NULL THEN '/' ELSE r.location END,
								r.url
							FROM results r
							INNER JOIN usertargets ut ON(r.scraperid = ut.id)
							INNER JOIN userkeywords uk ON(LOWER(r.title) LIKE('%' || uk.text || '%'))
							LEFT JOIN userfavouriteresults uft ON(r.id = uft.resultid)
							WHERE r.createdat > NOW() - INTERVAL '7 days'
							ORDER BY 3 DESC;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		job := Job{}
		if err = rows.Scan(
			&job.Id,
			&job.IsSaved,
			&job.CreatedDate,
			&job.TargetName,
			&job.KeywordText,
			&job.Title,
			&job.Location,
			&job.Url); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		jobs = append(jobs, job)
	}
	rows.Close()
	return
}
