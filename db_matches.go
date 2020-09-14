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
    rows, err := Db.Query(`
                            WITH
                                userkeywords AS(
                                    SELECT
                                        uk.keywordid
                                    FROM userskeywords uk
                                    WHERE uk.userid = $1
                                    AND uk.deletedat IS NULL),
                                usertargets AS(
                                    SELECT
                                        ut.targetid
                                    FROM userstargets ut
                                    WHERE ut.userid = $1
                                    AND ut.deletedat IS NULL)
                            SELECT
                                TO_CHAR(r.createdat, 'YYYY-MM-DD'),
                                s.name,
                                r.title,
                                r.url,
                                k.text
                            FROM userkeywords uk
                            INNER JOIN keywords k ON(uk.keywordid = k.id)
                            INNER JOIN matches m ON(uk.keywordid = m.keywordid)
                            INNER JOIN results r ON(m.resultid = r.id)
                            INNER JOIN scrapers s ON(r.scraperid = s.id)
                            INNER JOIN targets t ON(s.targetid = t.id)
                            INNER JOIN usertargets ut ON(t.id = ut.targetid)
                            WHERE m.createdat > current_date - interval '7' day
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
