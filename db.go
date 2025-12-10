package main

import (
	"fmt"
	"time"
	_ "embed"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Country struct {
	Code string
	Name string
}

type Sport struct {
	Code string
	Name string
	IsTeam bool
}

type Athlete struct {
	ID uint
	Name string
	Gender string
	Birthday time.Time

	CountryCode string
	Country Country
}

type Team struct {
	ID uint
	Name string

	CountryCode string
	Country Country

	SportCode string
	Sport Sport

	Members []Athlete
}

type Site struct {
	ID uint
	Name string
}

type Competition struct {
	ID uint
	Time time.Time

	SportCode string
	Sport Sport

	SiteID string
	Site Site

	Athletes []Athlete
	Teams []Team
}

//go:embed schema.sql
var dbSchema string

var db *sql.DB

func dbOpen(path string) error {
	var err error

	if db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s", path)); err != nil {
		return fmt.Errorf("failed to open database: %s", err)
	}

	if _, err = db.Exec(dbSchema); err != nil {
		return fmt.Errorf("failed to init db schema: %s", err)
	}

	return nil
}

func getCountries() ([]Country, error) {
	var countries []Country

	rows, err := db.Query("SELECT code, name FROM countries;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		country := Country{}
		if err := rows.Scan(&country.Code, &country.Name); err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return countries, nil
}

func addCountry(code string, name string) error {
	_, err := db.Exec("INSERT INTO countries ( code, name ) VALUES ( ?, ? );", code, name)
	return err
}

func deleteCountry(code string) error {
	_, err := db.Exec("DELETE FROM countries WHERE code = ?;", code)
	return err
}

func getSports() ([]Sport, error) {
	var sports []Sport

	rows, err := db.Query("SELECT code, name, is_team FROM sports;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sport := Sport{}
		if err := rows.Scan(&sport.Code, &sport.Name, &sport.IsTeam); err != nil {
			return nil, err
		}
		sports = append(sports, sport)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sports, nil
}

func addSport(code string, name string, team bool) error {
	_, err := db.Exec("INSERT INTO sports ( code, name, is_team ) VALUES ( ?, ?, ? );", code, name, team)
	return err
}

func deleteSport(code string) error {
	_, err := db.Exec("DELETE FROM sports WHERE code = ?;", code)
	return err
}

