// Package Check checks for errors
package check

import (
	"fmt"
	"net/http"
)

// CheckErr returns an internal server error if errors are encountered.
func CheckErr(w http.ResponseWriter, e error) {
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// CheckDbErr returns an error if there was an error communicating with database.
func CheckDbErr(w http.ResponseWriter, e error) {
	if e != nil {
		fmt.Fprintln(w, "Error comunicating with database")
		return
	}
}
