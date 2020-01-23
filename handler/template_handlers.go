package handler

import (
	"database/sql"
	"html/template"
	"log"
	validate "loginregistration/validation"
	"net/http"
)


var tplUpdate = template.Must(template.ParseFiles("./templates/update.html"))

//User struct to pass data to the html templates
type User struct {
	Username string
}

//LoginPageHandler for rendering Login page
func LoginPageHandler(w http.ResponseWriter, r *http.Request)  {
	tplLogin, err := template.ParseFiles("./templates/login.html","./templates/base.html")
	if err !=nil {
		log.Println(err)
	}
	tplLogin.Execute(w,nil)
}

//LoginHandler for handling post data from login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	result := db.QueryRow("select password from users where username=$1", username)
	var obtainedPassword string
	err := result.Scan(&obtainedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("<script>alert('No user exist!')</script>"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if obtainedPassword != password {
		w.Write([]byte("<script>alert('Login Failed!')</script>"))
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		http.Redirect(w,r,"/",302)
	}

}

//SignUpPageHandler for rendering sign up page
func SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	tplRegister, err := template.ParseFiles("./templates/signup.html","./templates/base.html")
	if err !=nil {
		log.Println(err)
	}
	tplRegister.Execute(w, nil)
}

//SignUpHandler for getting post request and handle them
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	_username, _password1, _password2 := false, false, false
	_username = !validate.IsEmpty(username)
	_password1 = !validate.IsEmpty(password1)
	_password2 = !validate.IsEmpty(password2)
	if _username && _password1 && _password2 {
		if string(password1) != string(password2) {
			http.Redirect(w, r, "/signup" , 302)
		} else {
			if _, err := db.Query("insert into users values ($1, $2)", username, password1); err != nil {
				w.Write([]byte("<script>alert('Error occurred!')</script>"))
			} else {
				w.Write([]byte("<script>alert('Success! Please login')</script>"))
			}
		}
	}else {
		w.Write([]byte("<script>alert('Sorry! Fields can not be empty')</script>"))
	}

}

//DeleteHandler for handling request of deleting a user
func DeleteHandler(w http.ResponseWriter, r *http.Request){
	username := r.URL.Query().Get("username")
	_, err := db.Query("DELETE FROM users WHERE username=$1",username)
	if err != nil {
		panic(err.Error())
	}
	log.Println("DELETE")
	http.Redirect(w, r, "/", 301)
}

//UpdatePage to render user information update page
func UpdatePage(w http.ResponseWriter, r *http.Request)  {
	tplUpdate.Execute(w,nil)
}

//UpdateHandler for handling submitted update data
func UpdateHandler(w http.ResponseWriter, r *http.Request){
	userToBeUpdated := r.URL.Query().Get("username")
	username := r.FormValue("username")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")
	_username, _password1, _password2 := false, false, false
	_username = !validate.IsEmpty(username)
	_password1 = !validate.IsEmpty(password1)
	_password2 = !validate.IsEmpty(password2)
	if _username && _password1 && _password2{
		if string(password1) != string(password2) {
			http.Redirect(w, r, "/signup" , 302)
		} else {
			if _, err := db.Query("update users set username=$1,password=$2 where username =$3", username, password1,userToBeUpdated); err != nil {
				w.Write([]byte("<script>alert('Error occurred!')</script>"))
			} else {
				http.Redirect(w,r,"/",302)
			}
		}
	}else {
		w.Write([]byte("<script>alert('Sorry! Fields can not be empty')</script>"))
	}

}

//IndexHandler for rendering and handling Index page
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var Users []User
	for rows.Next() {
		var username string
		var user User
		err = rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		user.Username = username
		Users = append(Users, user)

	}
	log.Println(len(Users))
	tm := template.Must(template.ParseFiles("./templates/index.html","./templates/base.html"))
	errortm := tm.Execute(w ,Users)
	log.Println(errortm)
}
