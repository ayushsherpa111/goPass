package schema

const (
	Password_DB_Schema = `
CREATE TABLE IF NOT EXISTS "passwords" (
	"pid" INTEGER PRIMARY KEY,
	"Username" TEXT,
	"Email" TEXT,
	"Password" BLOB,
	"Nonce" BLOB,
	"Site" TEXT
);
	`
)
