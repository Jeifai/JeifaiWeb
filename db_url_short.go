package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func SelectUrlByShortUrl(shorturl string) (resultid int, url string) {
	fmt.Println(Gray(8-1, "Starting SelectUrlByShortUrl..."))
	_ = Db.QueryRow(`SELECT id, url FROM results WHERE urlshort = $1;`, shorturl).Scan(&resultid, &url)
	return
}

func InsertUserResultVisit(userid int, resultid int) {
	fmt.Println(Gray(8-1, "Starting InsertUserResultVisit..."))
	statement := `INSERT INTO usersresultsvisits (userid, resultid, createdat)
                  VALUES ($1, $2, current_timestamp);`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(userid, resultid)
}
