package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	stmtCreateTable = `create table if not exists wallpaper (filename TEXT primary key, content BLOB);`
	stmtInsertData  = `insert into wallpaper values (?, ?);`
	stmtQueryData   = `select * from wallpaper;`
)

type SqliteStorage struct {
	database *sql.DB
}

func NewSqliteStorage(dbFile string) *SqliteStorage {
	db, err := sql.Open("sqlite3", dbFile)
	checkErr(err)
	_, err = db.Exec(stmtCreateTable)
	checkErr(err)

	return &SqliteStorage{database: db}
}

func (s *SqliteStorage) Destroy() {
	s.database.Close()
}

func (s *SqliteStorage) Load(name string) []byte {
	tx, err := s.database.Begin()
	checkErr(err)
	stmt, err := tx.Prepare(stmtQueryData)
	checkErr(err)
	defer stmt.Close()
	var result []byte
	err = stmt.QueryRow("name").Scan(&result)
	checkErr(err)
	return result
}

func (s *SqliteStorage) Save(img []byte, name string) {
	tx, err := s.database.Begin()
	checkErr(err)
	stmt, err := tx.Prepare(stmtInsertData)
	checkErr(err)
	defer stmt.Close()
	_, err = stmt.Exec(name, img)
	checkErr(err)
	tx.Commit()
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
