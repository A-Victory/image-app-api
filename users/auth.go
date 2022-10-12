package users

import (
	"fmt"
	"net/http"

	"app/db"
)

// LoggedIn validates the user session to see if they are logged in.
func LoggedIn(w http.ResponseWriter, r *http.Request) error {
	Db := db.New()
	c, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return err
	}

	row := Db.QueryRow("SELECT username FROM session WHERE sID = ?", c.Value)
	if err := row.Err(); err != nil {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return err
	}
	return nil
}

// GetUser returns the user associated with the session.
func GetUser(w http.ResponseWriter, r *http.Request) (string, error) {
	Db := db.New()
	var name string
	c, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		fmt.Fprintln(w, "Authentication not found. Please login.")
		return "", nil
	}
	row := Db.QueryRow("SELECT username FROM session WHERE sId = ?", c.Value)
	row.Scan(&name)

	return name, row.Err()
}
