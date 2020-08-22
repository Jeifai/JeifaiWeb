package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/logrusorgru/aurora"
)

type Keyword struct {
	Id        int
	Text      string `validate:"required,max=30,min=3"`
	CreatedAt time.Time
}

type UserTargetKeyword struct {
	Id          int
	UserId      int
	TargetId    int
	KeywordId   int
	CreatedAt   time.Time
	CreatedDate string
	UpdatedAt   time.Time
	KeywordText string
	TargetUrl   string
	TargetName  string
}

type KeywordInfo struct {
	Name          string
	CreatedDate   string
	CountTargets  int
	CountMatches  int
	AvgMatchesDay float32
}

func (keyword *Keyword) CreateKeyword() (err error) {
	fmt.Println(Gray(8-1, "Starting CreateKeyword..."))
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
	return err
}

func (keyword *Keyword) KeywordByText() (err error) {
	fmt.Println(Gray(8-1, "Starting KeywordByText..."))
	err = Db.QueryRow(`SELECT
                         k.id
                       FROM keywords k
                       WHERE k.text=$1`, keyword.Text).Scan(&keyword.Id)
	return
}

func (user *User) GetUserTargetKeyword() (
	utks []UserTargetKeyword, err error) {
	fmt.Println(Gray(8-1, "Starting GetUserTargetKeyword..."))

	rows, err := Db.Query(`SELECT
                                utk.id,
                                utk.userid,
                                utk.targetid,
                                utk.keywordid,
                                utk.createdat,
                                TO_CHAR(utk.createdat, 'YYYY-MM-DD'),
                                utk.updatedat,
                                k.text,
                                t.name
                            FROM userstargetskeywords utk
                            LEFT JOIN keywords k ON(utk.keywordid = k.id)
                            LEFT JOIN targets t ON(utk.targetid = t.id)
                            WHERE utk.userid = $1
                            AND utk.deletedat IS NULL
                            ORDER BY utk.updatedat DESC`, user.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		utk := UserTargetKeyword{}
		if err = rows.Scan(
			&utk.Id,
			&utk.UserId,
			&utk.TargetId,
			&utk.KeywordId,
			&utk.CreatedAt,
			&utk.CreatedDate,
			&utk.UpdatedAt,
			&utk.KeywordText,
			&utk.TargetName); err != nil {
			return
		}
		utks = append(utks, utk)
	}
	rows.Close()
	return
}

func (user *User) InfoKeywordsByUser() (keywordsInfo []KeywordInfo) {
	fmt.Println(Gray(8-1, "Starting InfoKeywordsByUser..."))
	rows, err := Db.Query(`
						SELECT DISTINCT
							k.text,
							MAX(utk.createdat::date),
							COUNT(DISTINCT utk.targetid),
							COUNT(DISTINCT m.id),
							ROUND(1.0 * COUNT(DISTINCT m.id) / COUNT(DISTINCT m.createdat::date), 1)
						FROM userstargetskeywords utk
						LEFT JOIN keywords k ON(utk.keywordid = k.id)
						LEFT JOIN matches m ON(m.keywordid = k.id AND m.createdat > utk.createdat)
						WHERE utk.userid = $1
						AND utk.deletedat IS NULL
						GROUP BY 1;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		keywordInfo := KeywordInfo{}
		if err = rows.Scan(
			&keywordInfo.Name,
			&keywordInfo.CreatedDate,
			&keywordInfo.CountTargets,
			&keywordInfo.CountMatches,
			&keywordInfo.AvgMatchesDay); err != nil {
			panic(err.Error())
		}
		keywordsInfo = append(keywordsInfo, keywordInfo)
	}
	rows.Close()
	if err != nil {
		panic(err.Error())
	}
	return
}

func SetUserTargetKeyword(
	user User, targets []Target, keyword Keyword) (err error) {
	fmt.Println(Gray(8-1, "Starting SetUserTargetKeyword..."))

	valueStrings := []string{}
	valueArgs := []interface{}{}
	timeNow := time.Now()
	for i, elem := range targets {
		str1 := "$" + strconv.Itoa(1+i*4) + ","
		str2 := "$" + strconv.Itoa(2+i*4) + ","
		str3 := "$" + strconv.Itoa(3+i*4) + ","
		str4 := "$" + strconv.Itoa(4+i*4)
		str_n := "(" + str1 + str2 + str3 + str4 + ")"
		valueStrings = append(valueStrings, str_n)
		valueArgs = append(valueArgs, user.Id)
		valueArgs = append(valueArgs, elem.Id)
		valueArgs = append(valueArgs, keyword.Id)
		valueArgs = append(valueArgs, timeNow)
	}
	smt := `INSERT INTO userstargetskeywords (
            userid, targetid, keywordid, createdat)
            VALUES %s
            ON CONFLICT (userid, targetid, keywordid) 
            DO UPDATE SET deletedat = NULL, updatedat = current_timestamp`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err = Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
	return
}

func SetDeletedAtInUserTargetKeywordMultiple(utks []UserTargetKeyword) (err error) {
	fmt.Println(Gray(8-1, "Starting SetDeletedAtInUserTargetKeywordMultiple..."))

	valueArray := []string{}
	timeNow := time.Now()
	_ = timeNow
	for _, elem := range utks {
		str1 := "(userid=" + strconv.Itoa(elem.UserId) + " AND "
		str2 := "keywordid=(SELECT id FROM keywords WHERE text='" + elem.KeywordText + "') AND "
		str3 := "targetid=(SELECT id FROM targets WHERE name='" + elem.TargetName + "'))"
		str_n := str1 + str2 + str3
		valueArray = append(valueArray, str_n)
	}
	valueString := strings.Join(valueArray, " OR ")

	smt := `UPDATE userstargetskeywords
            SET deletedat = current_timestamp
            WHERE id IN (
                SELECT
                    id
                FROM userstargetskeywords
                WHERE (%s));`

	smt = fmt.Sprintf(smt, valueString)
	_, err = Db.Exec(smt)
	if err != nil {
		panic(err.Error())
	}
	return
}
