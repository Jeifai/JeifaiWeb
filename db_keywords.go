package main

import (
	"fmt"
	"time"

	. "github.com/logrusorgru/aurora"
)

type Keyword struct {
	Id          int
	Text        string `validate:"required,max=30,min=2"`
	CreatedAt   time.Time
	CreatedDate string
}

type UserTargetKeyword struct {
	CreatedDate string
	KeywordText string
	TargetName  string
}

func (keyword *Keyword) InsertKeyword() {
	fmt.Println(Gray(8-1, "Starting InsertKeyword..."))
	statement := `INSERT INTO keywords (text, createdat)
                  VALUES ($1, current_timestamp)
                  RETURNING id, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		keyword.Text,
	).Scan(
		&keyword.Id,
		&keyword.CreatedAt,
	)
	if err != nil {
		panic(err.Error())
	}
}

func (keyword *Keyword) SelectKeywordByText() {
	fmt.Println(Gray(8-1, "Starting SelectKeywordByText..."))
	_ = Db.QueryRow(`SELECT
                         k.id
                       FROM keywords k
                       WHERE k.text=$1`, keyword.Text).Scan(&keyword.Id)
}

func (user *User) SelectUserKeywordByUserAndKeyword(keyword Keyword) (userKeywordId int) {
	fmt.Println(Gray(8-1, "Starting SelectUserKeywordByUserAndKeyword..."))
	_ = Db.QueryRow(`SELECT
                         uk.id
                       FROM userskeywords uk
                       WHERE uk.userid = $1
                       AND uk.keywordid = $2
                       AND uk.deletedat IS NULL;`, user.Id, keyword.Id).Scan(&userKeywordId)
	return
}

func (user *User) SelectKeywordsByUser() (keywords []Keyword) {
	fmt.Println(Gray(8-1, "Starting SelectKeywordsByUser..."))
	rows, err := Db.Query(`
							SELECT
								k.text,
								TO_CHAR(MIN(ut.createdat::date), 'YYYY-MM-DD')
							FROM userskeywords ut
							LEFT JOIN keywords k ON(ut.keywordid = k.id)
							WHERE ut.userid = $1
							AND ut.deletedat IS NULL
							GROUP BY 1;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		keyword := Keyword{}
		if err = rows.Scan(
			&keyword.Text,
			&keyword.CreatedDate); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		keywords = append(keywords, keyword)
	}
	rows.Close()
	return
}

func SelectKeywordsByAll() (keywords []string) {
	fmt.Println(Gray(8-1, "Starting SelectKeywordsByAll..."))
	rows, err := Db.Query(`
							SELECT
								DISTINCT k.text
							FROM keywords k;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var keyword string
		if err = rows.Scan(
			&keyword); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		keywords = append(keywords, keyword)
	}
	rows.Close()
	return
}

func (user *User) InsertUserKeyword(keyword Keyword) {
	fmt.Println(Gray(8-1, "Starting InsertUserKeyword..."))

	statement := `INSERT INTO userskeywords (userid, keywordid, createdat)
                  VALUES ($1, $2, current_timestamp);`
	stmt, err := Db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow(user.Id, keyword.Id)
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) UpdateDeletedAtInUsersKeywords(keyword Keyword) {
	fmt.Println(Gray(8-1, "Starting UpdateDeletedAtInUsersKeywords..."))

	statement := `UPDATE userskeywords
				  SET deletedat = current_timestamp
				  WHERE userid = $1
				  AND keywordid = $2;`
	stmt, err := Db.Prepare(statement)
	defer stmt.Close()
	stmt.QueryRow(user.Id, keyword.Id)
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) GetUserTargetKeyword() (utks []UserTargetKeyword) {
	fmt.Println(Gray(8-1, "Starting GetUserTargetKeyword..."))

	rows, err := Db.Query(`SELECT
                                TO_CHAR(utk.createdat, 'YYYY-MM-DD'),
                                k.text,
                                t.name
                            FROM userstargetskeywords utk
                            LEFT JOIN keywords k ON(utk.keywordid = k.id)
                            LEFT JOIN targets t ON(utk.targetid = t.id)
                            WHERE utk.userid = $1
                            AND utk.deletedat IS NULL
                            ORDER BY utk.updatedat DESC`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		utk := UserTargetKeyword{}
		if err = rows.Scan(
			&utk.CreatedDate,
			&utk.KeywordText,
			&utk.TargetName); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		utks = append(utks, utk)
	}
	rows.Close()
	return
}
