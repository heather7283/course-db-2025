-- проверить вьюхи
SELECT * FROM competition_results;
SELECT * FROM country_medals;

-- получить инфу о спортсмене
SELECT
    a.name AS athlete_name,
    c.name AS country,
    s.name AS sport,
    sites.name AS venue,
    comp.time AS date,
    ca.place AS place
FROM athletes a
JOIN competition_athletes ca ON a.id = ca.athlete_id
JOIN competitions comp ON ca.competition_id = comp.id
JOIN sports s ON comp.sport_code = s.code
JOIN sites ON comp.site_id = sites.id
JOIN countries c ON a.country_code = c.code
WHERE a.name = 'Алексей Козлов'
ORDER BY comp.time;

-- инфа про команды и кол-во участников в команде
SELECT
    t.name AS team,
    c.name AS country,
    s.name AS sport,
    COUNT(*) AS member_count
FROM teams t
JOIN countries c ON t.country_code = c.code
JOIN sports s ON t.sport_code = s.code
JOIN team_members tm ON tm.athlete_id = t.id
GROUP BY tm.team_id;

-- узнать какие соревнования и где проводятся в определённый день
SELECT
    s.name AS sport,
    st.name AS site,
    c.time AS time
FROM competitions c
JOIN sports s ON c.sport_code = s.code
JOIN sites st ON c.site_id = st.id
WHERE date(c.time) = '2024-01-15';


