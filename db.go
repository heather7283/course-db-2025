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
	ID int
	Name string
	Gender string
	Birthday time.Time
	CountryName string
}

type Team struct {
	ID int
	Name string

	Country Country
	Sport Sport

	Members []Athlete
}

type Site struct {
	ID int
	Name string
}

type Competition struct {
	ID int
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

func getAthletes() ([]Athlete, error) {
	var athletes []Athlete

	rows, err := db.Query("SELECT a.id, a.name, a.gender, a.birthday, c.name FROM athletes a JOIN countries c ON c.code = a.country_code;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		athlete := Athlete{}
		err := rows.Scan(&athlete.ID, &athlete.Name, &athlete.Gender, &athlete.Birthday, &athlete.CountryName)
		if err != nil {
			return nil, err
		}
		athletes = append(athletes, athlete)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return athletes, nil
}

func addAthlete(name string, isMale bool, birthday time.Time, countryCode string) error {
	var gender string
	if isMale {
		gender = "M"
	} else {
		gender = "F"
	}
	_, err := db.Exec("INSERT INTO athletes ( name, gender, birthday, country_code ) VALUES ( ?, ?, ?, ? );",
		name, gender, birthday.Unix(), countryCode)
	return err
}

func deleteAthlete(ID int) error {
	_, err := db.Exec("DELETE FROM athletes WHERE id = ?;", ID)
	return err
}

func getSites() ([]Site, error) {
	var sites []Site

	rows, err := db.Query("SELECT id, name FROM sites;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		site := Site{}
		err := rows.Scan(&site.ID, &site.Name)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

func addSite(name string) error {
	_, err := db.Exec("INSERT INTO sites ( name ) VALUES ( ? );", name)
	return err
}

func deleteSite(ID int) error {
	_, err := db.Exec("DELETE FROM sites WHERE id = ?;", ID)
	return err
}

func getTeams() ([]Team, error) {
	var teams []Team

	rows, err := db.Query(`
		SELECT t.id, t.name,
		       c.code, c.name,
		       s.code, s.name, s.is_team
		FROM teams t
		JOIN countries c ON c.code = t.country_code
		JOIN sports s ON s.code = t.sport_code
		ORDER BY t.id;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		team := Team{}
		country := Country{}
		sport := Sport{}

		err := rows.Scan(&team.ID, &team.Name, &country.Code, &country.Name, &sport.Code, &sport.Name, &sport.IsTeam)
		if err != nil {
			return nil, err
		}

		team.Country = country
		team.Sport = sport

		memberRows, err := db.Query(`
			SELECT a.id, a.name, a.gender, a.birthday, c2.name
			FROM team_members tm
			JOIN athletes a ON a.id = tm.athlete_id
			JOIN countries c2 ON c2.code = a.country_code
			WHERE tm.team_id = ?;
		`, team.ID)
		if err != nil {
			return nil, err
		}

		var members []Athlete
		for memberRows.Next() {
			athlete := Athlete{}
			err := memberRows.Scan(&athlete.ID, &athlete.Name, &athlete.Gender,
				&athlete.Birthday, &athlete.CountryName)
			if err != nil {
				memberRows.Close()
				return nil, err
			}
			members = append(members, athlete)
		}
		memberRows.Close()

		if err = memberRows.Err(); err != nil {
			return nil, err
		}

		team.Members = members
		teams = append(teams, team)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

func addTeam(name string, countryCode string, sportCode string) error {
	_, err := db.Exec("INSERT INTO teams (name, country_code, sport_code) VALUES (?, ?, ?);",
		name, countryCode, sportCode)
	return err
}

func deleteTeam(ID int) error {
	_, err := db.Exec("DELETE FROM teams WHERE id = ?;", ID)
	return err
}

func addAthleteToTeam(teamID int, athleteID int) error {
	_, err := db.Exec("INSERT INTO team_members (team_id, athlete_id) VALUES (?, ?);", teamID, athleteID)
	return err
}

func deleteAthleteFromTeam(teamID int, athleteID int) error {
	_, err := db.Exec("DELETE FROM team_members WHERE team_id = ? AND athlete_id = ?;", teamID, athleteID)
	return err
}

