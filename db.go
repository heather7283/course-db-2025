package main

import (
	"fmt"
	"time"
	_ "embed"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Country struct {
	Code string `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type Sport struct {
	Code string `gorm:"primaryKey"`
	Name string `gorm:"not null;unique"`
	IsTeam bool `gorm:"not null;default:0"`
}

type Athlete struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Gender string `gorm:"not null;check:gender in ('F', 'M');index"`
	Birthday time.Time `gorm:"check:birthday < current_timestamp"`

	CountryCode string `gorm:"index"`
	Country Country `gorm:"foreignKey:CountryCode;references:Code"`
}

type Team struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null"`

	CountryCode string `gorm:"index"`
	Country Country `gorm:"foreignKey:CountryCode;references:Code"`

	SportCode string `gorm:"index"`
	Sport Sport `gorm:"foreignKey:SportCode;references:Code"`

	Members []Athlete `gorm:"many2many:team_members;foreignKey:ID;references:ID"`
}

type Site struct {
	ID uint `gorm:"primaryKey"`
	Name string `gorm:"not null;unique"`
}

type Competition struct {
	ID uint `gorm:"primaryKey"`
	Time time.Time `gorm:"not null;check:time < current_timestamp"`

	SportCode string `gorm:"index"`
	Sport Sport `gorm:"foreignKey:SportCode;references:Code"`

	SiteID string `gorm:"index"`
	Site Site `gorm:"foreignKey:SiteID;references:ID"`

	Athletes []Athlete `gorm:"many2many:competition_athletes;foreignKey:ID;references:ID"`
	Teams []Team `gorm:"many2many:competition_teams;foreignKey:ID;references:ID"`
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

