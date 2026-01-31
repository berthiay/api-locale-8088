CREATE TABLE IF NOT EXISTS reports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS decision_proposals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS validation_requests (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  content TEXT NOT NULL
);
