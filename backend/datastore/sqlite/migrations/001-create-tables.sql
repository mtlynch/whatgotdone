CREATE TABLE IF NOT EXISTS user_preferences (
  username TEXT PRIMARY KEY,
  entry_template TEXT
  );

CREATE TABLE IF NOT EXISTS user_profiles (
  username TEXT PRIMARY KEY,
  about_markdown TEXT,
  email TEXT,
  twitter TEXT,
  mastodon TEXT
  );

CREATE TABLE IF NOT EXISTS journal_entries(
  username TEXT,
  date TEXT,
  last_modified TEXT,
  markdown TEXT,
  is_draft INTEGER,
  PRIMARY KEY (username, date, is_draft)
  );

CREATE TABLE IF NOT EXISTS follows(
  follower TEXT,
  leader TEXT,
  created TEXT,
  PRIMARY KEY (leader, follower)
  );

CREATE TABLE IF NOT EXISTS entry_reactions(
  entry_author TEXT,
  entry_date TEXT,
  reacting_user TEXT,
  reaction TEXT,
  timestamp TEXT,
  PRIMARY KEY (entry_author, entry_date, reacting_user)
  );

CREATE TABLE IF NOT EXISTS pageviews(
  path TEXT PRIMARY KEY,
  views INTEGER,
  last_updated TEXT
  );
