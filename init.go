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
	"github.com/sirupsen/logrus"
)

const isLocal = true

var log = logrus.New()

var Db *sql.DB

func init() {
	init_log()
	init_db()
}

func init_log() {
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)
}

func init_db() {
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
		log.Warn("Unsuccessfully connected to the database")
		return
	}
	log.Info("Successfully connected to the database")
	fmt.Println(Green("Successfully connected to the database"))
}

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

func GenerateToken() (token string) {
	token = uniuri.NewLen(40)
	return
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
