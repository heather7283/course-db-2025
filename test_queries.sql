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

