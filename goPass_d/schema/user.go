package schema

const (
	User_DB_SCHEMA = `
		CREATE TABLE IF NOT EXISTS "user" (
			"HomePath" TEXT PRIMARY KEY,
			"Salt" BLOB,
			"Username" TEXT,
			"Email" TEXT,
			"Hash" BLOB
		);
	`
)
