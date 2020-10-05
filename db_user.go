package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/logrusorgru/aurora"
)

type User struct {
	Id                int    `db:"id"`
	UserName          string `db:"username"    validate:"required,min=5,max=15"`
	Email             string `db:"email"       validate:"required,email"`
	LoginPassword     string
	Password          string         `db:"password" validate:"required,min=8,`
	CreatedAt         time.Time      `db:"createdat"`
	UpdatedAt         time.Time      `db:"updatedat"`
	DeletedAt         time.Time      `db:"deletedat"`
	FirstName         sql.NullString `db:"firstname"`
	LastName          sql.NullString `db:"lastname"`
	DateOfBirth       sql.NullString `db:"dateofbirth"`
	Country           sql.NullString `db:"country"`
	City              sql.NullString `db:"city"`
	Gender            sql.NullString `db:"gender"`
	CurrentPassword   string         `                 validate:"required,min=8,eqfield=Password"`
	NewPassword       string         `db:"newpassword" validate:"eqfield=RepeatNewPassword"`
	RepeatNewPassword string         `                 validate:"eqfield=NewPassword"`
	InvitationCode    string         `                 validate:"required"`
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
                  RETURNING id, uuid, email, userid, createdat`
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
                      WHERE uuid = $1`,
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

func CreateUser(email string, username string, password string) {
	fmt.Println(Gray(8-1, "Starting Create..."))
	statement := `INSERT INTO users (email, username, password, createdat, updatedat)
                  VALUES ($1, $2, $3, current_timestamp, current_timestamp)
                  RETURNING id, createdat, updatedat`
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
	err := Db.QueryRow(`
					SELECT
						id,
                        password,
					FROM users
					WHERE email = $1`,user.Email).Scan(user.Id,&user.Password,)
	if err != nil {
		panic(err.Error())
	}
	return
}

func UserById(userId int) (user User) {
	fmt.Println(Gray(8-1, "Starting UserById..."))
	err := Db.QueryRow(`SELECT
                        id,
                        username,
                        email,
                        password,
                        createdat,
                        updatedat,
                        firstname,
                        lastname,
                        dateofbirth,
                        country,
                        city,
                        gender
                      FROM users
                      WHERE id = $1`, userId).
		Scan(
			&user.Id,
			&user.UserName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.FirstName,
			&user.LastName,
			&user.DateOfBirth,
			&user.Country,
			&user.City,
			&user.Gender,
		)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (user User) UpdateUser() {
	fmt.Println(Gray(8-1, "Starting UpdateUser..."))
	statement := `UPDATE users SET 
                    username=$2,
                    email=$3,
                    password=$4,
                    firstname=$5,
                    lastname=$6,
                    country=$7,
                    city=$8,
                    gender=$9,
                    dateofbirth=$10,
                    updatedat=$11
                  WHERE id=$1`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		user.Id,
		user.UserName,
		user.Email,
		user.NewPassword,
		user.FirstName.String,
		user.LastName.String,
		user.Country.String,
		user.City.String,
		user.Gender.String,
		user.DateOfBirth.String,
		time.Now())

	if err != nil {
		panic(err.Error())
	}
}

func (user User) UpdateUserUpdates() {
	fmt.Println(Gray(8-1, "Starting UpdateUserUpdates..."))
	statement := `INSERT INTO usersupdates (userid, data, createdat) 
                    VALUES ($1, $2, $3)`

	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	user_json, err := json.Marshal(user)
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(
		user.Id,
		user_json,
		time.Now())
	if err != nil {
		panic(err.Error())
	}
}

func (user *User) CreateToken(token string) {
	fmt.Println(Gray(8-1, "Starting CreateToken..."))
	statement := `INSERT INTO resetpasswords (userid, token, createdat, expiredat)
                  VALUES ($1, $2, $3, $4)`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	created_at := time.Now()
	expired_at := time.Now().Local().Add(time.Minute * time.Duration(10))

	stmt.QueryRow(
		user.Id,
		token,
		created_at,
		expired_at,
	)
}

func UserByToken(token string) (user User)  {
	fmt.Println(Gray(8-1, "Starting UserByToken..."))
	err := Db.QueryRow(`SELECT
                        u.id,
                        u.email,
                        u.username
                      FROM resetpasswords r
                      LEFT JOIN users u ON (r.userid = u.id)
                      WHERE r.token = $1
                      AND current_timestamp < r.expiredat
                      AND r.consumedat IS NULL`, token).
		Scan(user.Id, user.Email, user.UserName)
	if err != nil {
		panic(err.Error())
	}
	return
}

func (user *User) ChangePassword(password string) {
	fmt.Println(Gray(8-1, "Starting ChangePassword..."))
	statement := `UPDATE users SET password=$1 WHERE id=$2`
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
