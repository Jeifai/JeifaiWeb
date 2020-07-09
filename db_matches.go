package main

import (
	"fmt"
	"time"

	. "github.com/logrusorgru/aurora"
)

type Match struct {
	Target      string
	Id          int
	ScraperId   int
	Title       string
	Url         string
	CreatedAt   time.Time
	CreatedDate string
}

func (user *User) MatchesByUser() (matches []Match, err error) {
	fmt.Println(Gray(8-1, "Starting MatchesByUser..."))
	rows, err := Db.Query(`SELECT DISTINCT
                                s.name,
                                r.createdat,
                                TO_CHAR(r.createdat, 'YYYY-MM-DD'),
                                r.title,
                                r.url
                            FROM matches m
                            INNER JOIN keywords k ON(m.keywordid = k.id)
                            INNER JOIN results r ON(m.resultid = r.id)
                            INNER JOIN scrapers s ON(r.scraperid = s.id)
                            INNER JOIN userstargetskeywords utk ON(k.id = utk.keywordid)
                            WHERE m.createdat > current_date - interval '3' day
                            AND utk.userid = $1
                            ORDER BY r.createdat DESC;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		match := Match{}
		if err = rows.Scan(
			&match.Target,
			&match.CreatedAt,
			&match.CreatedDate,
			&match.Title,
			&match.Url); err != nil {
			return
		}
		matches = append(matches, match)
	}
	rows.Close()
	return
}
