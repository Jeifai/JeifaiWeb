package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

type Match struct {
	CreatedDate string
	Target      string
	Title       string
	Url         string
	KeywordText string
}

func (user *User) MatchesByUser() (matches []Match) {
	fmt.Println(Gray(8-1, "Starting MatchesByUser..."))
	rows, err := Db.Query(`SELECT DISTINCT
							    TO_CHAR(r.createdat, 'YYYY-MM-DD'),
							    s.name,
							    r.title,
							    r.url,
							    k.text
							FROM matches m
							INNER JOIN keywords k ON(m.keywordid = k.id)
							INNER JOIN results r ON(m.resultid = r.id)
							INNER JOIN scrapers s ON(r.scraperid = s.id)
							INNER JOIN userstargetskeywords utk ON(k.id = utk.keywordid)
							WHERE m.createdat > current_date - interval '3' day
							AND utk.userid = $1
							AND utk.deletedat IS NULL
							ORDER BY 1 DESC;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		match := Match{}
		if err = rows.Scan(
			&match.CreatedDate,
			&match.Target,
			&match.Title,
			&match.Url,
			&match.KeywordText); err != nil {
			return
		}
		matches = append(matches, match)
	}
	rows.Close()
	if err != nil {
		panic(err.Error())
	}
	return
}
