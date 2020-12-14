package main

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

// GCP
var gateway = os.Getenv("gateway")

var mysql_server = os.Getenv("mysql_server")
var mysql_connect = "bitly:bitly@tcp(" + mysql_server + ":3306)/bitly"

// AWS
// var mysql_connect = "bitly:bitly@tcp(10.0.2.245:3306)/bitly"
// var gateway = "34.215.253.104:8000/lr"

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

func getLongURLFromDB(shorturl_code string) string {
	dbConn := connectToMySQL()
	var (
		longurl string
	)
	rows, err := dbConn.Query("select longurl from links where shorturl_code = ?", shorturl_code)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&longurl)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(longurl)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return longurl
}

func getStatistics(args ...interface{}) map[string]map[string]string {
	dbConn := connectToMySQL()

	var url string
	var linkstats = make(map[string]map[string]string)

	var (
		shorturl_code string
		no_of_clicks  string
	)
	rows, err := dbConn.Query("select shorturl_code, no_of_clicks from links order by no_of_clicks desc limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&shorturl_code, &no_of_clicks)
		if err != nil {
			log.Fatal(err)
		}

		url = "http://"+ gateway + "/" + shorturl_code

		lastOneMinCount := getAccessCount(shorturl_code, "minute", dbConn)
		lastOneHourCount := getAccessCount(shorturl_code, "hour", dbConn)
		lastOneDayCount := getAccessCount(shorturl_code, "day", dbConn)

		linkstats[url] = map[string]string{
			"Last 1 Minute": lastOneMinCount,
			"Last 1 Hour":   lastOneHourCount,
			"Last 1 Day":    lastOneDayCount,
			"All Time":      no_of_clicks,
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return linkstats
}

func getAccessCount(shorturl_code string, metric string, dbConn *sql.DB) string {
	rows, err := dbConn.Query("select count(shorturl_code) from trend_stats where shorturl_code='" + shorturl_code + "' and accessed_at > date_sub(now(), interval 1 " + metric + ")")
	if err != nil {
		log.Fatal(err)
	}
	var count string
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&count)
	}
	return count
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