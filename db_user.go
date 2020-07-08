package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Id                int    `db:"id"`
	UserName          string `db:"username"    validate:"min=1"`
	Email             string `db:"email"       validate:"email"`
	LoginPassword     string
	Password          string         `db:"password"`
	CreatedAt         time.Time      `db:"createdat"`
	UpdatedAt         time.Time      `db:"updatedat"`
	DeletedAt         time.Time      `db:"deletedat"`
	FirstName         sql.NullString `db:"firstname"`
	LastName          sql.NullString `db:"lastname"`
	DateOfBirth       sql.NullString `db:"dateofbirth"`
	Country           sql.NullString `db:"country"`
	City              sql.NullString `db:"city"`
	Gender            sql.NullString `db:"gender"`
	CurrentPassword   string         `                 validate:"required,eqfield=Password"`
	NewPassword       string         `db:"newpassword" validate:"eqfield=RepeatNewPassword"`
	RepeatNewPassword string         `                 validate:"eqfield=NewPassword"`
	InvitationCode    string
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

func session(request *http.Request) (sess Session, err error) {
	cookie, err := request.Cookie("_cookie")
	if err == nil {
		sess = Session{Uuid: cookie.Value}
		if ok, _ := sess.Check(); !ok {
			panic(err.Error())
		}
	}
	return
}

// Create a new session for an existing user
func (user *User) CreateSession() (session Session, err error) {
	statement := `INSERT INTO sessions (uuid, email, userid, createdat)
                  VALUES ($1, $2, $3, $4)
                  RETURNING id, uuid, email, userid, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
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
	return
}

// Get the session for an existing user
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow(`SELECT
                        id, 
                        uuid, 
                        email, 
                        userid, 
                        createdat
                      FROM sessions
                      WHERE userid = $1
                      AND deletedat IS NULL`,
		user.Id,
	).
		Scan(
			&session.Id,
			&session.Uuid,
			&session.Email,
			&session.UserId,
			&session.CreatedAt,
		)
	return
}

// Check if session is valid in the database
func (session *Session) Check() (valid bool, err error) {
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
func (session *Session) SetSessionDeletedAtByUUID() (err error) {
	statement := `UPDATE sessions SET deletedat = current_timestamp WHERE uuid = $1;`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return
}

// Create a new user, save user info into the database
func (user *User) Create() (err error) {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := `INSERT INTO users
                  (username, email, password, createdat, updatedat)
                  VALUES ($1, $2, $3, $4, $5)
                  RETURNING id, createdat, updatedat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(
		user.UserName,
		user.Email,
		Encrypt(user.Password),
		time.Now(),
		time.Now(),
	).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return
}

// Get a single user given the email
func (user *User) UserByEmail() (err error) {
	err = Db.QueryRow(`SELECT
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
                      WHERE email = $1`,
		user.Email,
	).
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
	return
}

// Get a single user given the email
func (user *User) UserById() (err error) {
	err = Db.QueryRow(`SELECT
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
                      WHERE id = $1`,
		user.Id,
	).
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
	return
}

func (user User) UpdateUser() {
	fmt.Println("Starting UpdateUser...")

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
	fmt.Println("Starting UpdateUserUpdates...")
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
}

func (user *User) CreateToken(token string) (err error) {
	fmt.Println("Starting CreateToken...")
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
	return
}

func UserByToken(token string) (user User) {
	current_timestamp := time.Now()
	err := Db.QueryRow(`SELECT
                        u.id,
                        u.email,
                        u.username
                      FROM resetpasswords r
                      LEFT JOIN users u ON (r.userid = u.id)
                      WHERE r.token = $1
                      AND $2 < r.expiredat
                      AND r.consumedat IS NULL`,
		token, current_timestamp).Scan(&user.Id, &user.Email, &user.UserName)
	_ = err
	return
}

func (user *User) ChangePassword(password string) (err error) {
	fmt.Println("Starting ChangePassword...")
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
	return
}
