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
	CreatedDate string
	KeywordText string
	TargetName  string
}

type KeywordInfo struct {
	Name             string
	CountResults     int
	CountResultsPerc string
	CountResultsDay  int
	CreatedDate      string
	CountTargets     int
	LastWeekMatches  int
	TotalMatches     int
	AvgMatchesDay    float32
}

func (keyword *Keyword) CreateKeyword() {
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
	if err != nil {
		panic(err.Error())
	}
}

func (keyword *Keyword) KeywordByText() (err error) {
	fmt.Println(Gray(8-1, "Starting KeywordByText..."))
	err = Db.QueryRow(`SELECT
                         k.id
                       FROM keywords k
                       WHERE k.text=$1`, keyword.Text).Scan(&keyword.Id)
	return
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

func (user *User) InfoKeywordsByUser() (keywordsInfo []KeywordInfo) {
	fmt.Println(Gray(8-1, "Starting InfoKeywordsByUser..."))
	rows, err := Db.Query(`
						WITH 
							keywords_macro AS (
								WITH
									userkeywords AS (
										SELECT DISTINCT
											k.text,
											k.id
										FROM userstargetskeywords utk
										LEFT JOIN keywords k ON(utk.keywordid = k.id)
										LEFT JOIN users u ON(utk.userid = u.id)
										WHERE u.id = $1
										AND utk.deletedat IS NULL),
									all_results AS (
										SELECT
											COUNT(DISTINCT r.id) AS count_id
										FROM results r)
								SELECT
									uk.text,
									COUNT(DISTINCT r.id) AS count_results,
									ROUND(100.0 * COUNT(DISTINCT r.id) / ar.count_id, 2)::text || '%' AS count_results_perc,
									COUNT(DISTINCT r.id) / COUNT(DISTINCT r.createdat::date) AS count_results_per_day
								FROM userkeywords uk
								LEFT JOIN results r ON(LOWER(r.title) LIKE('%' || uk.text || '%'))
								LEFT JOIN all_results ar ON(1 = 1)
								GROUP BY 1, ar.count_id)
							SELECT
								k.text,
								km.count_results,
								km.count_results_perc,
								km.count_results_per_day,
								TO_CHAR(MIN(utk.createdat), 'YYYY-MM-DD'),
								COUNT(DISTINCT utk.targetid),
								COUNT(DISTINCT CASE WHEN m.createdat >= now() - interval '1 week' THEN m.id END),
								COUNT(DISTINCT m.id),
								CASE 
									WHEN COUNT(DISTINCT m.id) = 0 THEN 0
									ELSE ROUND(1.0 * COUNT(DISTINCT m.id) / COUNT(DISTINCT m.createdat::date), 1)
								END
							FROM userstargetskeywords utk
							LEFT JOIN keywords k ON(utk.keywordid = k.id)
							LEFT JOIN matches m ON(utk.keywordid = m.keywordid AND utk.createdat < m.createdat)
							LEFT JOIN keywords_macro km ON(k.text = km.text)
							WHERE utk.userid = $1
							AND utk.deletedat IS NULL
							GROUP BY 1, 2, 3, 4;`, user.Id)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		keywordInfo := KeywordInfo{}
		if err = rows.Scan(
			&keywordInfo.Name,
			&keywordInfo.CountResults,
			&keywordInfo.CountResultsPerc,
			&keywordInfo.CountResultsDay,
			&keywordInfo.CreatedDate,
			&keywordInfo.CountTargets,
			&keywordInfo.LastWeekMatches,
			&keywordInfo.TotalMatches,
			&keywordInfo.AvgMatchesDay); err != nil {
			panic(err.Error())
		}
		keywordsInfo = append(keywordsInfo, keywordInfo)
	}
	rows.Close()
	return
}

func SetUserTargetKeyword(
	user User, targets []Target, keyword Keyword) {
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

	_, err := Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (user *User) SetDeletedAtInUserTargetKeywordMultiple(utks []UserTargetKeyword) {
	fmt.Println(Gray(8-1, "Starting SetDeletedAtInUserTargetKeywordMultiple..."))

	valueArray := []string{}
	timeNow := time.Now()
	_ = timeNow
	for _, elem := range utks {
		str1 := "(userid=" + strconv.Itoa(user.Id) + " AND "
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
	_, err := Db.Exec(smt)
	if err != nil {
		panic(err.Error())
	}
}
