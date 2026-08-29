INSERT INTO academic_year (label)
SELECT 'TEST TA / SEMESTER 2'
WHERE NOT EXISTS (
    SELECT 1 FROM academic_year WHERE label = 'TEST TA / SEMESTER 2'
);
