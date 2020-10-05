package main

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/logrusorgru/aurora"
)

type User struct {
	Id            int
	Email         string
	UserName      string
	LoginPassword string
	Password      string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func GetSession(request *http.Request) (sess Session) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = Session{Uuid: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			panic(err.Error())
		}
	}
	if err != nil {
		return Session{}
	}
	return
}

func (user *User) CreateSession() (session Session) {
	fmt.Println(Gray(8-1, "Starting CreateSession..."))
	statement := `INSERT INTO sessions (uuid, email, userid, createdat)
                  VALUES ($1, $2, $3, $4)
                  RETURNING id, uuid, email, userid, createdat;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(
		createUUID(),
		user.Email,
		user.Id,
		time.Now(),
	).Scan(
		&session.Id,
		&session.Uuid,
		&session.Email,
		&session.UserId,
		&session.CreatedAt,
	)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (session *Session) CheckSession() (valid bool, err error) {
	fmt.Println(Gray(8-1, "Starting Check..."))
	err = Db.QueryRow(`SELECT
                        id,
                        uuid,
                        email,
                        userid,
                        createdat
                      FROM sessions
                      WHERE uuid = $1;`,
		session.Uuid,
	).
		Scan(
			&session.Id,
			&session.Uuid,
			&session.Email,
			&session.UserId,
			&session.CreatedAt,
		)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// Set deletedat by uuid
func (session *Session) SetSessionDeletedAtByUUID() {
	fmt.Println(Gray(8-1, "Starting SetSessionDeletedAtByUUID..."))
	statement := `UPDATE sessions SET deletedat = current_timestamp WHERE uuid = $1;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	if err != nil {
		panic(err.Error())
	}
}

func InsertUser(email string, username string, password string) {
	fmt.Println(Gray(8-1, "Starting InsertUser..."))
	statement := `INSERT INTO users (email, username, password, createdat, updatedat)
                  VALUES ($1, $2, $3, current_timestamp, current_timestamp)
                  RETURNING id, createdat, updatedat;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	stmt.QueryRow(email, username, Encrypt(password))
	if err != nil {
		panic(err.Error())
	}
}

func UserByEmail(email string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByEmail..."))
	_ = Db.QueryRow(`
					SELECT
						id,
						email,
                        password
					FROM users
					WHERE email = $1;`, email).Scan(&user.Id, &user.Email, &user.Password)
	return
}

func (user *User) CreateToken(token string) {
	fmt.Println(Gray(8-1, "Starting CreateToken..."))
	statement := `INSERT INTO resetpasswords (userid, token, createdat, expiredat) VALUES ($1, $2, $3, $4);`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	created_at := time.Now()
	expired_at := time.Now().Local().Add(time.Minute * time.Duration(10))
	stmt.QueryRow(user.Id, token, created_at, expired_at)
}

func UserByToken(token string) (user User) {
	fmt.Println(Gray(8-1, "Starting UserByToken..."))
	_ = Db.QueryRow(`
					SELECT
						u.id,
                        u.email,
                        u.username
					FROM resetpasswords r
					LEFT JOIN users u ON (r.userid = u.id)
					WHERE r.token = $1
					AND current_timestamp < r.expiredat
					AND r.consumedat IS NULL;`, token).Scan(&user.Id, &user.Email, &user.UserName)
	return
}

func (user *User) UpdatePassword(password string) {
	fmt.Println(Gray(8-1, "Starting UpdatePassword..."))
	statement := `UPDATE users SET password=$1 WHERE id=$2;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(
		password,
		user.Id,
	)
}

func (user *User) UpdateResetPassword(token string) {
	fmt.Println(Gray(8-1, "Starting UpdateResetPassword..."))
	statement := `UPDATE resetpasswords SET consumedat=current_timestamp WHERE userid=$1 AND token=$2;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	stmt.QueryRow(user.Id, token)
}
