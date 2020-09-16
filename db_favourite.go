package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func (user *User) InsertFavourite(resultid int) {
	fmt.Println(Gray(8-1, "Starting InsertFavourite..."))
	statement := `INSERT INTO favouriteresults (userid, resultid, createdat)
                  VALUES ($1, $2, current_timestamp);`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		&user.Id,
		resultid,
	)
}

func (user *User) DeleteFavourite(resultid int) {
	fmt.Println(Gray(8-1, "Starting InsertFavourite..."))
	statement := `
                    UPDATE favouriteresults
                    SET deletedat = current_timestamp
                    WHERE userid = $1
                    AND resultid = $2
                    AND deletedat IS NULL;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		&user.Id,
		resultid,
	)
}

func (user *User) SelectJobsByUserAndFavourite() (jobs []Job) {
    fmt.Println(Gray(8-1, "Starting SelectJobsByUserAndFavourite..."))
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
                            INNER JOIN userfavouriteresults uft ON(r.id = uft.resultid)
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
