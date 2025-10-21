-- Страны
CREATE TABLE countries (
    code CHAR(3) PRIMARY KEY,

    name TEXT NOT NULL UNIQUE
);

-- Виды спорта
CREATE TABLE sports (
    code CHAR(3) NOT NULL PRIMARY KEY,

    name TEXT NOT NULL UNIQUE,
    -- Групповой или нет (по умолчанию нет)
    is_team BOOLEAN NOT NULL DEFAULT FALSE
);

-- Спортсмены
CREATE TABLE athletes (
    id INTEGER PRIMARY KEY,

    name TEXT NOT NULL,
    gender CHAR(1) NOT NULL CHECK ( gender IN ( 'F', 'M' ) ),
    birthday TIMESTAMP NOT NULL CHECK ( birthday < CURRENT_TIMESTAMP ),

    country_code CHAR(3) NOT NULL,

    FOREIGN KEY ( country_code ) REFERENCES countries ( code )
);

CREATE INDEX idx_athletes_country_code ON athletes ( country_code );
CREATE INDEX idx_athletes_gender ON athletes ( gender );

-- Команды (для командных видов спорта)
CREATE TABLE teams (
    id INTEGER PRIMARY KEY,

    name TEXT NOT NULL,

    country_code CHAR(3) NOT NULL,
    sport_code CHAR(3) NOT NULL,

    FOREIGN KEY ( country_code ) REFERENCES countries ( code ),
    FOREIGN KEY ( sport_code ) REFERENCES sports ( code )
);

CREATE INDEX idx_teams_country_code ON teams ( country_code );
CREATE INDEX idx_teams_sport_code ON teams ( sport_code );

-- Таблица для связи между командами и их участниками
CREATE TABLE team_members (
    team_id INTEGER,
    athlete_id INTEGER,

    PRIMARY KEY ( team_id, athlete_id ),
    FOREIGN KEY ( team_id ) REFERENCES teams ( id ),
    FOREIGN KEY ( athlete_id ) REFERENCES athletes ( id )
);

-- Площадки
CREATE TABLE sites (
    id INTEGER PRIMARY KEY,

    name TEXT NOT NULL UNIQUE
);

-- Проведённые соревнования
CREATE TABLE competitions (
    id INTEGER PRIMARY KEY,

    time TIMESTAMP NOT NULL CHECK ( time < CURRENT_TIMESTAMP ),

    sport_code CHAR(3) NOT NULL,
    site_id INTEGER NOT NULL,

    FOREIGN KEY ( sport_code ) REFERENCES sports ( code ),
    FOREIGN KEY ( site_id ) REFERENCES sites ( id )
);

CREATE INDEX idx_competitions_sport_code ON competitions ( sport_code );
CREATE INDEX idx_competitions_site_id ON competitions ( site_id );

-- Таблица для связи атлетов и соревнований, в которых они участвовали
CREATE TABLE competition_athletes (
    competition_id INTEGER,
    athlete_id INTEGER,

    -- Какое место атлет занял в этом соревновании
    place INTEGER NOT NULL CHECK ( place > 0 ),

    -- Несколько атлетов не могут одновременно занять одно и то же место
    UNIQUE ( competition_id, place ),

    PRIMARY KEY ( competition_id, athlete_id ),
    FOREIGN KEY ( competition_id ) REFERENCES competitions ( id ),
    FOREIGN KEY ( athlete_id ) REFERENCES athletes ( id )
);

-- Таблица для связи команд и соревнований, в которых они участвовали
CREATE TABLE competition_teams (
    competition_id INTEGER,
    team_id INTEGER,

    -- Какое место команда заняла в этом соревновании
    place INTEGER NOT NULL CHECK ( place > 0 ),

    -- Несколько команд не могут одновременно занять одно и то же место
    UNIQUE ( competition_id, place ),

    PRIMARY KEY ( competition_id, team_id ),
    FOREIGN KEY ( competition_id ) REFERENCES competitions ( id ),
    FOREIGN KEY ( team_id ) REFERENCES teams ( id )
);

-- вьюха которая объединяет индивидуальные и коммандные результаты
CREATE VIEW competition_results AS
    SELECT
        competition_id,
        'athlete' AS participant_type,
        athlete_id AS participant_id,
        place
    FROM competition_athletes
    UNION ALL
    SELECT
        competition_id,
        'team' AS participant_type,
        team_id AS participant_id,
        place
    FROM competition_teams
;

-- вьюха предоставляющая данные о кол-ве медалей у каждой страны
CREATE VIEW country_medals AS
    SELECT
        c.name AS country,
        COUNT(CASE WHEN cr.place = 1 THEN 1 END) AS gold,
        COUNT(CASE WHEN cr.place = 2 THEN 1 END) AS silver,
        COUNT(CASE WHEN cr.place = 3 THEN 1 END) AS bronze,
        COUNT(*) AS total
    FROM competition_results cr
    JOIN competitions comp ON cr.competition_id = comp.id
    LEFT JOIN athletes a ON cr.participant_type = 'athlete' AND cr.participant_id = a.id
    LEFT JOIN teams t ON cr.participant_type = 'team' AND cr.participant_id = t.id
    JOIN countries c ON a.country_code = c.code OR t.country_code = c.code
    GROUP BY c.code, c.name
    HAVING cr.place IN (1, 2, 3)
    ORDER BY gold DESC, silver DESC, bronze DESC;
;

-- триггер который проверяет чтобы в таблицу с командными результатами
-- нельзя было пихать индивидуальные виды спорта
CREATE TRIGGER ensure_individual_sport BEFORE INSERT ON competition_athletes FOR EACH ROW BEGIN
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

-- то же что и выше но наоборт
CREATE TRIGGER ensure_team_sport BEFORE INSERT ON competition_teams FOR EACH ROW BEGIN
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

-- триггер который не даст добавить атлета в команду чужой страны
CREATE TRIGGER ensure_team_members_country BEFORE INSERT ON team_members FOR EACH ROW BEGIN
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

