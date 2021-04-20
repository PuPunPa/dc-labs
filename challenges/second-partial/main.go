package main

import (
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	//If it says that Go could not found the mrepository just type
	//go get github.com/gorilla/mux
	//to fix it
	"github.com/gorilla/mux"
)

//Login
//curl -u username:password http://localhost:8080/login
//Logout
//curl -H "Authorization: Bearer ONRhfKsUOH" http://localhost:8080/logout
//Upload
//curl -F 'data=dc-labs/challenges/second-partial/test.jpg' -H "Authorization: Bearer ONRhfKsUOH" http://localhost:8080/upload
//curl -F 'data=dc-labs/challenges/second-partial/architecture.png' -H "Authorization: Bearer ONRhfKsUOH" http://localhost:8080/upload
//Status
//curl -H "Authorization: Bearer ONRhfKsUOH" http://localhost:8080/status

// Logged: Structure to store logged in users info, key is Token, value is User
var online = make(map[string]string)

//Users: stores all users, key Username, value is Password
var users = make(map[string]string)

type Image struct {
	Token, Path, Name string
	size              int
}

// Global variable user. All functions are able to access to it
func main() {
	router := mux.NewRouter()
	users["username"] = "password"
	//All routes for API
	//Need functions
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/status", status)
	router.HandleFunc("/upload", upload)
	http.ListenAndServe(":8080", router)
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	usr, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Username and Password are required"
}
`
		w.Write([]byte(msg))
		return
	}
	if !checkPass(usr, pwd) {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Invalid Username and/or Password"
}
`
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	token := GenerateRandomString(10)
	msg := `
{
	"message": "Hi ` + usr + ` welcome to the DPIP System"
	"token" "` + token + `"
}
`
	w.Write([]byte(msg))
	online[token] = usr
	return
}

func checkPass(usr string, pwd string) bool {
	if p, ok := users[usr]; ok {
		return p == pwd
	}
	return false
}

func logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a token" ` + token + `smiles
}
`
		w.Write([]byte(msg))
		return
	}
	token = token[7:]
	if !loggedIn(token) {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a valid token"
}
`
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	msg := `
{
	"message": "Bye ` + online[token] + `, your token has been revoked"
}
`
	w.Write([]byte(msg))
	delete(online, token)
	return
}

func loggedIn(token string) bool {
	if _, ok := online[token]; ok {
		return true
	}
	return false
}

func status(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if len(token) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a token" ` + token + `
}
`
		w.Write([]byte(msg))
		return
	}
	token = token[7:]
	if !loggedIn(token) {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a valid token"
}
`
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	msg := `
{
	"message": "Hi ` + online[token] + `, the DPIP System is Up and Running "
	"time": "` + time.Now().Format("2006-01-02 15:04:05") + `"
}
`
	w.Write([]byte(msg))
	return
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generate a string of random characters of given length
func GenerateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx := rand.Int63() % int64(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}
	return sb.String()
}

// Receive image and return name and size of it
func upload(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	msg := ""
	if len(token) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a token" ` + token + `
}
`
		w.Write([]byte(msg))
		return
	}
	token = token[7:]
	if !loggedIn(token) {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a valid token"
}
`
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, imgData, err := r.FormFile("data")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"There was an error while parsing the image"
	"Error:" "` + err.Error() + `"
}
`
		w.Write([]byte(msg))
		return
	}
	msg = `
{
	"message": "An image has been successfully uploaded",
	"filename": "` + imgData.Filename + `",
	"size": "` + strconv.FormatInt(imgData.Size, 10) + `kb"
}
`
	w.Write([]byte(msg))
	return
}
