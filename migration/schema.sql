CREATE TABLE IF NOT EXISTS country (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  iso2 TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS answer (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  selected_country INTEGER NOT NULL,
  correct_country INTEGER NOT NULL,
  is_correct BOOLEAN NOT NULL,
  FOREIGN KEY(selected_country) REFERENCES country(id),
  FOREIGN KEY(correct_country) REFERENCES country(id)
);

