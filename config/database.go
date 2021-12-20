package config

import (
	"database/sql"
	"friday/config/utils"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DBCon *sql.DB
var JwtRedis *redis.Client
var PubSubRedis *redis.Client
var Sqlite3 *sql.DB

type DatabaseInfo struct {
	Name     string
	Host     string
	Password string
	Root     string
}

func getSourceName(db DatabaseInfo) string {
	return db.Root + ":" + db.Password + "@tcp(" + db.Host + ":3306)/" + db.Name
}

func InitDB() string {
	err := godotenv.Load(".env")
	utils.FatalError{Error: err}.Handle()

	go initMysql()
	go initRedis()

	return "InitDB Success"
}

func initRedis() {
	//Initializing redis
	jwtDsn := os.Getenv("REDIS_JWT_DSN")
	if len(jwtDsn) == 0 {
		jwtDsn = "localhost:6379"
	}
	JwtRedis = redis.NewClient(&redis.Options{
		Addr: jwtDsn, //redis port
	})

	psDsn := os.Getenv("REDIS_PUB_SUB_DSN")
	if len(psDsn) == 0 {
		psDsn = "localhost:6380"
	}
	PubSubRedis = redis.NewClient(&redis.Options{
		Addr: psDsn, //redis port
	})
}

func InitSqlite3() *sql.DB {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `	
	create table if not exists users(
    	id    	   varchar(255) not null primary key,
    	level      integer default 1 not null,
    	name       varchar(255) not null,
    	email      varchar(255) not null,
    	password   varchar(255),
    	created_at datetime default CURRENT_TIMESTAMP,
    	updated_at datetime default CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt = `	
	CREATE TABLE IF NOT EXISTS rooms (
		id VARCHAR(255) NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		private TINYINT NULL
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}

	Sqlite3 = db

	return db
}

func initMysql() {
	databaseInfo := DatabaseInfo{
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ROOT"),
	}

	db, err := sql.Open("mysql", getSourceName(databaseInfo))
	utils.FatalError{Error: err}.Handle()

	DBCon = db
}