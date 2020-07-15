package main

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"os"

	"github.com/dchest/uniuri"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

const isLocal = false

var Db *sql.DB

func init() {
	if isLocal {
		err := godotenv.Load()
		if err != nil {
			panic(err.Error())
		}
	}

	psqlInfo := os.Getenv("POSTGRES_CONNECTION")

	db, err := sql.Open("postgres", psqlInfo)

	Db = db

	if err != nil {
		panic(err.Error())
	}
	if err = Db.Ping(); err != nil {
		Db.Close()
		fmt.Println(Red("Unsuccessfully connected to the database"))
		return
	}
	fmt.Println(Green("Successfully connected to the database"))
}

func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		return
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

func GenerateToken() (token string) {
	token = uniuri.NewLen(40)
	return
}