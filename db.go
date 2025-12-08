package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

const trigger_ensure_individual_sport = `
	CREATE TRIGGER IF NOT EXISTS ensure_individual_sport
	BEFORE INSERT ON competition_athletes FOR EACH ROW BEGIN
	    SELECT
	    CASE
	        WHEN (
	            SELECT s.is_team
	            FROM competitions c
	            JOIN sports s ON c.sport_code = s.code
	            WHERE c.id = NEW.competition_id
	        ) = TRUE
	        THEN RAISE(ABORT, 'Командный спорт — нельзя добавлять одиночного атлета')
	    END;
	END;
`

const trigger_ensure_team_sport = `
	CREATE TRIGGER IF NOT EXISTS ensure_team_sport
	BEFORE INSERT ON competition_teams FOR EACH ROW BEGIN
	    SELECT
	    CASE
	        WHEN (
	            SELECT s.is_team
	            FROM competitions c
	            JOIN sports s ON c.sport_code = s.code
	            WHERE c.id = NEW.competition_id
	        ) = FALSE
	        THEN RAISE(ABORT, 'Индивидуальный спорт — нельзя добавлять команду')
	    END;
	END;
`

const trigger_ensure_team_members_country = `
	CREATE TRIGGER IF NOT EXISTS ensure_team_members_country
	BEFORE INSERT ON team_members FOR EACH ROW BEGIN
	    SELECT
	    CASE
	        WHEN (
	            SELECT a.country_code != t.country_code
	            FROM athletes a
	            JOIN teams t ON t.id = NEW.team_id
	            WHERE a.id = NEW.athlete_id
	        ) = TRUE
	        THEN RAISE(ABORT, 'Спортсмен не соответствует команде по стране')
	    END;
	END;
`

var dbState struct {
	db *gorm.DB
	ctx context.Context
}

func dbOpen(ctx context.Context, path string) (error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	tables := []any{
		&Country{},
		&Sport{},
		&Athlete{},
		&Team{},
		&Site{},
		&Competition{},
	}
	if err := db.Migrator().AutoMigrate(tables...); err != nil {
		return fmt.Errorf("failed to create schema: %s", err)
	}

	triggers := []string{
		trigger_ensure_individual_sport,
		trigger_ensure_team_sport,
		trigger_ensure_team_members_country,
	}
	for _, trigger := range triggers {
		if err := db.Exec(trigger).Error; err != nil {
			return fmt.Errorf("failed to create trigger: %s", err)
		}
	}

	dbState.ctx = ctx
	dbState.db = db
	return nil
}

func getCountries() ([]Country, error) {
	return gorm.G[Country](dbState.db).Find(dbState.ctx)
}

func addCountry(code string, name string) error {
	return gorm.G[Country](dbState.db).Create(dbState.ctx, &Country{
		Code: code,
		Name: name,
	})
}

func deleteCountry(code string) error {
	_, err := gorm.G[Country](dbState.db).Where("code = ?", code).Delete(dbState.ctx)
	return err
}

