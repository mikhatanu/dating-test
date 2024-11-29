package db

import (
	"database/sql"
)

// for storing sql open

var DB *sql.DB
var DB_error error
