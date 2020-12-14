package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

// AWS
// var mysql_connect = "bitly:bitly@tcp(10.0.2.245:3306)/bitly"

// GCP
var mysql_server = os.Getenv("mysql_server")
var mysql_connect = "bitly:bitly@tcp(" + mysql_server + ":3306)/bitly"

// Local
// var mysql_connect = "root:bitly@tcp(localhost:3306)/bitly"

// Create Connection to Database
func connectToMySQL() *sql.DB {

	db, err := sql.Open("mysql", mysql_connect)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connection to DB Successfull")
	}
	return db
}

//Insert new shortlink in database
func insertShortLink(shorturl_code string, longurl string) {
	dbConn := connectToMySQL()
	stmt, err := dbConn.Prepare("INSERT into links(shorturl_code, longurl) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(shorturl_code, longurl); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Insert short link successful")
	}
}

func updateStatistics(shorturl_code string) {
	dbConn := connectToMySQL()

	insertStmt, err := dbConn.Prepare("insert into trend_stats(shorturl_code, accessed_at) values (?, current_timestamp)")
	if err != nil {
		log.Fatal(err)
	}
	defer insertStmt.Close()
	if _, err := insertStmt.Exec(shorturl_code); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Insert into trends successful")
	}

	updateStmt, err := dbConn.Prepare("update links set no_of_clicks = no_of_clicks+1 where shorturl_code = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer updateStmt.Close()
	if _, err := updateStmt.Exec(shorturl_code); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Update count in links successful")
	}
}

/*
create database bitly;

use bitly;

create table links(
	shorturl_code varchar(10) primary key,
	longurl mediumtext,
	no_of_clicks mediumint default 0,
	created_at timestamp default current_timestamp
	);

create table trend_stats(
	shorturl_code varchar(10),
	accessed_at timestamp
	);
*/
