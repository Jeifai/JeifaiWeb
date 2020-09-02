package main

import (
	"fmt"
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

func (target *Target) CreateTarget() {
	fmt.Println(Gray(8-1, "Starting CreateTarget..."))
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

func (target *Target) CreateUserTarget(user User) {
	fmt.Println(Gray(8-1, "Starting CreateUserTarget..."))
	statement := `INSERT INTO userstargets (userid, targetid, createdat) 
                  VALUES ($1, $2, $3)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(user.Id, target.Id, time.Now())
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

func (user *User) NotSelectedTargetsNamesByUser() (notSelectedTargetsNames []string) {
	fmt.Println(Gray(8-1, "Starting NotSelectedTargetsNamesByUser..."))

	err := Db.QueryRow(`
                WITH usertargets AS (
				  SELECT
				      ut.targetid
				  FROM userstargets ut
				  WHERE ut.deletedat IS NULL
				  AND ut.userid=$1)
				SELECT
				  ARRAY_AGG(t.name)
				FROM targets t
				LEFT JOIN usertargets ut ON(t.id = ut.targetid)
				WHERE ut.targetid IS NULL
				ORDER BY 1;`, user.Id).
		Scan(
			pq.Array(&notSelectedTargetsNames),
		)
	if len(notSelectedTargetsNames) == 0 { // vue-taggable-select does not work if name_targets is empty
		notSelectedTargetsNames = []string{" "}
	}
	if err != nil {
		panic(err.Error())
	}
	return
}

func (user *User) TargetsByUser() (targets []Target) {
	fmt.Println(Gray(8-1, "Starting TargetsByUser..."))
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
                                TO_CHAR(MAX(ls.createdat), 'YYYY-MM-DD') AS last_extraction,
                                COUNT(DISTINCT r.url) AS all_time_job,
                                SUM(CASE WHEN r.scrapingid = ls.scrapingid THEN 1 ELSE 0 END) AS actual_job_opens,
                                SUM(CASE WHEN (r.createdat > current_date - interval '7' day) THEN 1 ELSE 0 END) AS open_positions_last_7_days,
                                SUM(CASE WHEN (r.updatedat > current_date - interval '7' day AND r.updatedat < current_date - interval '1' day) THEN 1 ELSE 0 END) AS close_positions_last_7_days
                            FROM usertargets ut
                            LEFT JOIN results r ON(ut.scraperid = r.scraperid)
                            LEFT JOIN latest_scraping ls ON(r.scraperid = ls.scraperid)
                            GROUP BY 1, 2, 3
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

func (target *Target) TargetByName() {
	fmt.Println(Gray(8-1, "Starting TargetByName..."))
	err := Db.QueryRow(`SELECT
                         t.id
                       FROM targets t
                       WHERE t.name=$1`, target.Name).Scan(&target.Id)
	if err != nil {
		panic(err.Error())
	}
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

func (user *User) UsersTargetsByUserAndName(target Target) (err error) {
	fmt.Println(Gray(8-1, "Starting UsersTargetsByUserAndName..."))
	err = Db.QueryRow(`SELECT
                         t.id,
                         t.createdat 
                       FROM users u
                       INNER JOIN userstargets ut ON(u.id = ut.userid) 
                       INNER JOIN targets t ON(ut.targetid = t.id)
                       WHERE u.id=$1
                       AND t.name=$2
                       AND ut.deletedat IS NULL`, user.Id, target.Name).Scan(
		&target.Id, &target.CreatedAt)
	return
}

func (target *Target) SetDeletedAtInUsersTargetsByUserAndTarget(user User) {
	fmt.Println(Gray(8-1, "Starting SetDeletedAtInUsersTargetsByUserAndTarget..."))

	statement := `UPDATE userstargets
                  SET deletedat = current_timestamp
                  WHERE userid = $1
                  AND targetid = $2;`

	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id, target.Id)
	if err != nil {
		panic(err.Error())
	}
}

func (target *Target) SetDeletedAtIntUserTargetKeywordByUserAndTarget(user User) {
	fmt.Println(Gray(8-1, "Starting SetDeletedAtIntUserTargetKeywordByUserAndTarget..."))
	statement := `UPDATE userstargetskeywords
                  SET deletedat = current_timestamp
                  WHERE userid = $1
                  AND targetid = $2`

	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Id, target.Id)
	if err != nil {
		panic(err.Error())
	}
}
