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
	Name        string
	CreatedDate string
	JobsAll     int
	JobsNow     int
	Opened      int
	Closed      int
}

// Add a new target
func (target *Target) CreateTarget() (err error) {
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
	return err
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

func (user *User) UsersTargetsByUser() (targets []Target, err error) {
	fmt.Println(Gray(8-1, "Starting UsersTargetsByUser..."))
	rows, err := Db.Query(`SELECT
                            t.id,
                            t.name,
                            t.createdat,
                            TO_CHAR(t.createdat, 'YYYY-MM-DD')
                           FROM users u
                           INNER JOIN userstargets ut ON(u.id = ut.userid) 
                           INNER JOIN targets t ON(ut.targetid = t.id)
                           WHERE ut.deletedat IS NULL
                           AND u.id=$1
                           ORDER BY t.createdat DESC`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(
			&target.Id,
			&target.Name,
			&target.CreatedAt,
			&target.CreatedDate); err != nil {
			panic(err.Error())
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

func (user *User) InfoUsersTargetsByUser() (targetsinfo []TargetInfo, err error) {
	fmt.Println(Gray(8-1, "Starting InfoUsersTargetsByUser..."))
	rows, err := Db.Query(`WITH usertargets AS (
                                SELECT
                                    s.id AS scraperid,
                                    t.name,
                                    TO_CHAR(t.createdat, 'YYYY-MM-DD') AS createdat
                                FROM users u
                                INNER JOIN userstargets ut ON(u.id = ut.userid) 
                                INNER JOIN targets t ON(ut.targetid = t.id)
                                LEFT JOIN scrapers s ON(t.id = s.targetid)
                                WHERE ut.deletedat IS NULL
                                AND u.id=$1
                                ORDER BY t.createdat DESC)
                                SELECT
                                    ut.name,
                                    ut.createdat,
                                    COUNT(DISTINCT r.url) AS all_time_job,
                                    SUM(CASE WHEN r.updatedat::date = current_date::date THEN 1 ELSE 0 END) AS actual_job_opens,
                                    SUM(CASE WHEN (r.createdat > current_date - interval '7' day) THEN 1 ELSE 0 END) AS open_positions_last_7_days,
                                    SUM(CASE WHEN (r.updatedat > current_date - interval '7' day AND r.updatedat < current_date - interval '1' day) THEN 1 ELSE 0 END) AS close_positions_last_7_days
                                FROM usertargets ut
                                LEFT JOIN results r ON(ut.scraperid = r.scraperid)
                                GROUP BY 1, 2
                                ORDER BY 2;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		targetinfo := TargetInfo{}
		if err = rows.Scan(
			&targetinfo.Name,
			&targetinfo.CreatedDate,
			&targetinfo.JobsAll,
			&targetinfo.JobsNow,
			&targetinfo.Opened,
			&targetinfo.Closed); err != nil {
			panic(err.Error())
		}
		targetsinfo = append(targetsinfo, targetinfo)
	}
	rows.Close()
	return
}

func (user *User) NotSelectedUsersTargetsByUser() (targets []Target, err error) {
	fmt.Println(Gray(8-1, "Starting NotSelectedUsersTargetsByUser..."))
	rows, err := Db.Query(`WITH usertargets AS (
                                SELECT
                                    ut.targetid
                                FROM userstargets ut
                                WHERE ut.deletedat IS NULL
                                AND ut.userid=$1)
                            SELECT
                            t.name
                            FROM targets t
                            LEFT JOIN usertargets ut ON(t.id = ut.targetid)
                            WHERE ut.targetid IS NULL
                            ORDER BY 1;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(
			&target.Name); err != nil {
			panic(err.Error())
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

func (target *Target) TargetByName() (err error) {
	fmt.Println(Gray(8-1, "Starting TargetByName..."))
	err = Db.QueryRow(`SELECT
                         t.id
                       FROM targets t
                       WHERE t.name=$1`, target.Name).Scan(&target.Id)
	return
}

func TargetsByNames(targetNames []string) (targets []Target, err error) {
	fmt.Println(Gray(8-1, "Starting TargetsByNames..."))

	rows, err := Db.Query(`SELECT
                                t.id,
                                t.name
                            FROM targets t
                            WHERE t.name LIKE ANY($1)`, pq.Array(targetNames))
	if err != nil {
		return
	}
	for rows.Next() {
		target := Target{}
		if err = rows.Scan(&target.Id, &target.Name); err != nil {
			return
		}
		targets = append(targets, target)
	}
	rows.Close()
	return
}

func (user *User) UsersTargetsByUserAndName(name string) (target Target, err error) {
	fmt.Println(Gray(8-1, "Starting UsersTargetsByUserAndName..."))
	err = Db.QueryRow(`SELECT
                         t.id, 
                         t.name, 
                         t.createdat 
                       FROM users u
                       INNER JOIN userstargets ut ON(u.id = ut.userid) 
                       INNER JOIN targets t ON(ut.targetid = t.id)
                       WHERE u.id=$1
                       AND t.name=$2
                       AND ut.deletedat IS NULL`, user.Id, name).Scan(&target.Id, &target.Name, &target.CreatedAt)
	return
}

// Update userstargets in column deletedat
func (target *Target) SetDeletedAtInUsersTargetsByUserAndTarget(
	user User) (err error) {
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
	return
}

func (target *Target) SetDeletedAtIntUserTargetKeywordByUserAndTarget(user User) (err error) {
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
	return
}
