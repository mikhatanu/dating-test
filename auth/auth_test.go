package auth

import (
	"database/sql"
	"os"
	"testing"

	"github.com/mikhatanu/dating-test/db"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth(t *testing.T) {
	db.DB, db.DB_error = sql.Open("sqlite3", "./app.db")
	defer db.DB.Close()
	defer os.Remove("./app.db")
	createUserTable := `
		CREATE TABLE IF NOT EXISTS user (
			username text unique not null,
			password nvarchar(255) not null
		);
	`
	db.DB.Exec(createUserTable)
	t.Run("test signup", func(t *testing.T) {
		got := Signup("test123@email.com", "superstrong")
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})
	t.Run("test signup no username", func(t *testing.T) {
		got := Signup("", "superstrong")
		if got == nil {
			t.Errorf("got %v, want nil", got)
		}
	})
	t.Run("test signup no password", func(t *testing.T) {
		got := Signup("tes!@sdf.com", "")
		if got == nil {
			t.Errorf("got %v, want nil", got)
		}
	})
	t.Run("Check if user is signed up", func(t *testing.T) {
		got, _ := GetUser("test123@email.com")
		want := got.Username
		if want != "test123@email.com" {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("Check if user not found", func(t *testing.T) {
		got, _ := GetUser("sd@123dwd4.com")
		want := ""
		if got.Id != want {
			t.Errorf("got %v, want %v", got.Id, want)
		}
	})
	t.Run("Check password", func(t *testing.T) {
		password := "secret"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		got := CheckPassword(password, string(hashed))
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})
	t.Run("Check wrong password", func(t *testing.T) {
		password := "secret"
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		got := CheckPassword("newpas", string(hashed))
		if got == nil {
			t.Errorf("got %v, want nil", got)
		}
	})

}
