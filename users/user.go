// Package users is responside for the User related interface.
package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"app/check"
	"app/utils"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type Info struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserController struct {
	Db *sql.DB
}

// NewUserController creates a new UserController instance for connection to the database.
func NewUserController(db *sql.DB) *UserController {
	return &UserController{db}
}

// Signup opens up new user account, sending the neccessary informations to the database.
func (uc *UserController) Signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user := User{}
	var user_id int
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashPswrd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if err != nil {
		log.Fatal("Error while encypting password")
		return
	}

	row := uc.Db.QueryRow("SELECT uId FROM info WHERE username = ?", user.Username)
	row.Scan(&user_id)
	if user_id != 0 {
		fmt.Fprintln(w, "Username already taken")
		/*
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
		*/
		return
	} else {
		if user.Email == "" {
			fmt.Fprintln(w, "Please enter your email address")
			return
		}
		stmt, err := uc.Db.Prepare("INSERT INTO info(username, password, first_name, last_name, email) VALUES ( ?, ?, ?, ?, ? )")
		check.CheckDbErr(w, err)
		defer stmt.Close()
		req, err := stmt.Exec(user.Username, string(hashPswrd), user.FirstName, user.LastName, user.Email)
		if err != nil {
			check.CheckDbErr(w, err)
		} else {
			fmt.Println(req.RowsAffected())
		}
		fmt.Fprintln(w, "Signup successful. Please login.")

	}

}

// Profile returns the user's home profile.
// Returning images the user has uploaded
func (uc *UserController) Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}
	user, err := GetUser(w, r)
	check.CheckErr(w, err)
	var imageName string

	rows, err := uc.Db.Query("SELECT image_name FROM image WHERE username = ?", user)
	check.CheckDbErr(w, err)
	for rows.Next() {
		err := rows.Scan(&imageName)
		if err != nil {
			fmt.Fprintln(w, "No image found for user.")
			return
		}

		public_id := user + "/" + imageName
		res, err := utils.GetImage(public_id)
		check.CheckErr(w, err)

		resp, err := http.Get(res)
		check.CheckErr(w, err)
		fmt.Fprintln(w, resp.Status)
	}
	fmt.Fprintln(w, "Images are displayed")

}

// Login lets users login, creating cookies that are stored in the clients machine.
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	c, er := r.Cookie("session")
	if er == http.ErrNoCookie {
		info := Info{}
		var hashPswrd string
		var email string
		var user_id int

		e := json.NewDecoder(r.Body).Decode(&info)
		check.CheckErr(w, e)

		user_row := uc.Db.QueryRow("SELECT uId FROM info WHERE username = ?", info.Username)
		user_row.Scan(&user_id)
		if user_id == 0 {
			fmt.Fprintln(w, "User not found. PLease Signup.")
			/*
				w.Header().Set("Location", "/signup")
				w.WriteHeader(http.StatusTemporaryRedirect)
			*/
			return
		}

		rw := uc.Db.QueryRow("SELECT email FROM info WHERE username = ?", info.Username)
		rw.Scan(&email)
		if email != info.Email {
			fmt.Fprintln(w, "Email not registered.")
			return
		}

		row := uc.Db.QueryRow("SELECT password FROM info WHERE username = ?", info.Username)
		row.Scan(&hashPswrd)

		err := bcrypt.CompareHashAndPassword([]byte(hashPswrd), []byte(info.Password))
		if err != nil {
			fmt.Fprintln(w, "Password is incorrect.")
			return
		}

		cookie := SetCookie(w, r, info.Username)
		http.SetCookie(w, cookie)
		fmt.Fprintln(w, "Login Successful.")

		return
	} else {
		row := uc.Db.QueryRow("SELECT username FROM session WHERE sID = ?", c.Value)
		if err := row.Err(); err == nil {
			fmt.Fprintln(w, "You are already logged in.")
			w.Header().Set("Location", "/public")
			w.WriteHeader(http.StatusTemporaryRedirect)

			return
		}
	}

}

// Public is the main field of the app.
// It returns images from multiple users.
func (uc *UserController) Public(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}

	var imageNames string
	var u_name string

	rows, err := uc.Db.Query("SELECT image_name FROM image")
	check.CheckDbErr(w, err)
	for rows.Next() {
		err := rows.Scan(&imageNames)
		if err != nil {
			fmt.Fprintln(w, "An error occurred! Please try again")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		rows, err := uc.Db.Query("SELECT username FROM image WHERE image_name = ?", imageNames)
		check.CheckDbErr(w, err)
		for rows.Next() {
			err := rows.Scan(&u_name)
			if err != nil {
				fmt.Fprintln(w, "An error occurred! Please try again")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			public_id := u_name + "/" + imageNames
			res, err := utils.GetImage(public_id)
			if err != nil {
				log.Fatal("Error getting image: ", err)
			}

			resp, err := http.Get(res)
			check.CheckErr(w, err)
			fmt.Fprintln(w, resp.StatusCode)
		}

	}

	fmt.Fprintln(w, "Images are displayed")

}

// Search allows user to search other users profile.
func (uc *UserController) Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}

	var user_id int
	var imageNames string
	name := ps.ByName("user")
	row := uc.Db.QueryRow("SELECT uId FROM info WHERE username = ?", name)
	row.Scan(&user_id)
	if user_id == 0 {
		fmt.Fprint(w, "User not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rows, err := uc.Db.Query("SELECT image_name FROM image WHERE username = ?", name)
	check.CheckErr(w, err)

	for rows.Next() {
		err := rows.Scan(&imageNames)
		if err != nil {
			fmt.Fprintln(w, "Nothing to display.", err)
			return
		}

		public_id := name + "/" + imageNames

		res, err := utils.GetImage(public_id)
		if err != nil {
			log.Fatal("Error getting image: ", err)
		}

		resp, err := http.Get(res)
		check.CheckErr(w, err)
		fmt.Fprintln(w, resp.Status)

		fmt.Fprintln(w, "Images sre displayed")
	}

}

// UploadImage allows users to upload images.
func (uc *UserController) UploadImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!\n ")
		return
	}
	u_name, err := GetUser(w, r)
	check.CheckErr(w, err)

	img, prt, err := r.FormFile("image")
	check.CheckErr(w, err)
	//defer img.Close()
	name, ext, found := strings.Cut(prt.Filename, ".")
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !(ext == "jpg" || ext == "jpeg" || ext == "png") {
		fmt.Fprintln(w, "Image format not supported.")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// Converting to string
	image_name := fmt.Sprintf("%x", name)
	public_id := u_name + "/" + image_name

	e := utils.UploadImage(public_id, img)
	if e != nil {
		log.Fatal("Failed to upload image: ", e)
	}

	stmt, err := uc.Db.Prepare("INSERT INTO image(username, image_name) VALUES ( ?, ? )")
	check.CheckDbErr(w, err)
	defer stmt.Close()

	_, err = stmt.Exec(u_name, image_name)
	if err != nil {
		fmt.Fprintln(w, "An error occurred.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintln(w, "Image upload successful.")
	}

}

// UpdateInfo allows users to update their informations
func (uc *UserController) UpdateInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}
	u_name, e := GetUser(w, r)
	check.CheckErr(w, e)
	details := User{}
	err := json.NewDecoder(r.Body).Decode(&details)
	check.CheckErr(w, err)

	hashPswrd, err := bcrypt.GenerateFromPassword([]byte(details.Password), 8)
	if err != nil {
		log.Fatal("Error while encypting password")
	}

	if details.FirstName != "" {
		stmt, err := uc.Db.Prepare("UPDATE info SET first_name = ? WHERE username = ? ")
		check.CheckDbErr(w, err)
		defer stmt.Close()
		_, err = stmt.Exec(details.FirstName, u_name)
		check.CheckErr(w, err)
	}

	if details.LastName != "" {
		stmt, err := uc.Db.Prepare("UPDATE info SET last_name = ? WHERE username = ? ")
		check.CheckDbErr(w, err)
		defer stmt.Close()
		_, err = stmt.Exec(details.LastName, u_name)
		check.CheckErr(w, err)
	}

	if details.Password != "" {
		stmt, err := uc.Db.Prepare("UPDATE info SET password = ? WHERE username = ? ")
		check.CheckDbErr(w, err)
		defer stmt.Close()
		_, err = stmt.Exec(string(hashPswrd), u_name)
		check.CheckErr(w, err)
	}

	if details.Email != "" {
		stmt, err := uc.Db.Prepare("UPDATE info SET email = ? WHERE username = ? ")
		check.CheckDbErr(w, err)
		defer stmt.Close()
		_, err = stmt.Exec(details.Email, u_name)
		check.CheckErr(w, err)
	}

	fmt.Fprintln(w, "Info updated successfully")

}

// DeleteImage allows a user to delete a previously uploaded image.
func (uc *UserController) DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}
	u_name, e := GetUser(w, r)
	var image_name string
	check.CheckErr(w, e)
	id, e := strconv.Atoi(ps.ByName("id"))
	check.CheckErr(w, e)
	err := uc.Db.QueryRow("SELECT image_name FROM image WHERE username = ? AND imageID = ?", u_name, id).Scan(&image_name)
	check.CheckErr(w, err)
	public_id := u_name + "/" + image_name
	resp, err := utils.DeleteImage(public_id)
	check.CheckErr(w, err)

	stmt, err := uc.Db.Prepare("DELETE FROM image WHERE username=? AND imageID = ?")
	check.CheckDbErr(w, err)
	defer stmt.Close()
	res, err := stmt.Exec(u_name, id)
	check.CheckDbErr(w, err)
	res_int, _ := res.RowsAffected()
	fmt.Fprintf(w, "%d, %s Image deleted.", res_int, resp)

}

// DeleteUser allows a user delete themselves and associated information from the server.
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}
	u_name, e := GetUser(w, r)
	check.CheckErr(w, e)
	c := DeleteCookie(w, r)
	http.SetCookie(w, c)
	var imageNames string
	rows, err := uc.Db.Query("SELECT image_name FROM image WHERE username = ?", u_name)
	check.CheckDbErr(w, err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&imageNames)
		if err == rows.Err() {
			fmt.Fprintln(w, err)
		}
		public_id := u_name + "/" + imageNames
		fmt.Fprintln(w, imageNames)

		res, err := utils.DeleteImage(public_id)
		if err != nil {
			log.Fatal("Error deleting image: ", err)
		}
		fmt.Fprintln(w, res)
	}

	stmt1, err := uc.Db.Prepare("DELETE FROM image WHERE username = ?")
	check.CheckErr(w, err)
	_, err = stmt1.Exec(u_name)
	check.CheckDbErr(w, err)
	stmt, err := uc.Db.Prepare("DELETE FROM info WHERE username = ?")
	check.CheckErr(w, err)
	_, err = stmt.Exec(u_name)
	check.CheckDbErr(w, err)

	fmt.Fprintln(w, "User deleted! To continue with our service, create a new account.")

}

// Logout allows a user logout from the server deleting cookies from their machines.
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := LoggedIn(w, r); err != nil {
		fmt.Fprintln(w, "You are not logged in!")
		return
	}

	c := DeleteCookie(w, r)
	http.SetCookie(w, c)

}
