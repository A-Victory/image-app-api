package users

import (
	"net/http"
	"time"

	"app/check"
	"app/db"

	"github.com/google/uuid"
)

// SetCookie set the cookie for the user.
func SetCookie(w http.ResponseWriter, r *http.Request, u string) *http.Cookie {
	Db := db.New()
	cId, _ := uuid.NewRandom()
	value := cId.String()
	expiresAt := time.Now().Add(time.Hour * 24)

	c := &http.Cookie{
		Name:    "session",
		Value:   value,
		Expires: expiresAt,
	}

	stmt, err := Db.Prepare("INSERT INTO session(sID, username) VALUES ( ?, ? )")
	check.CheckDbErr(w, err)
	defer stmt.Close()

	_, err = stmt.Exec(value, u)
	check.CheckErr(w, err)

	return c
}

// DeleteCookie removes a cookie from the user's machine, including from the session database.
func DeleteCookie(w http.ResponseWriter, r *http.Request) *http.Cookie {
	Db := db.New()
	t, _ := r.Cookie("session")

	c := &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}

	stmt, err := Db.Prepare("DELETE FROM session WHERE sId = ?")
	check.CheckDbErr(w, err)
	defer stmt.Close()

	_, err = stmt.Exec(t.Value)
	check.CheckErr(w, err)

	return c

}
