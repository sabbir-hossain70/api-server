package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

var Secret = []byte("secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type mpCred map[string]string

var CredList mpCred

func Login(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	fmt.Println("Cred ", cred, "   ", err)
	if err != nil {
		http.Error(w, "Can not decode the Credentials... ", http.StatusBadRequest)
		return
	}
	password, ok := CredList[cred.Username]
	if !ok {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if password != cred.Password {
		http.Error(w, "wrong Password", http.StatusBadRequest)
		return
	}
	EndTime := time.Now().Add(10 * time.Minute)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"aud": "Sabbir Hossain",
		"exp": EndTime.Unix(),
	})
	fmt.Println("tokenString : ", tokenString)
	if err != nil {
		http.Error(w, "Can not generate jwt", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: EndTime,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully..."))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged out"))
}

type mpBook map[int]Book

var BookList mpBook

type mpAuthor map[string]Author

var AuthorList mpAuthor

type Author struct {
	AuthorFirstName string   `json:"fname"`
	AuthorLastName  string   `json:"lname"`
	AuthorHomeTown  string   `json:"hometown"`
	AuthorAge       int      `json:"age"`
	AuthorPhone     string   `json:"phone"`
	AuthorAllBooks  []string `json:"publishedbooks"`
}

type Book struct {
	Title           string `json:"title"`
	PublishDate     string `json:"publishdate"`
	BookNo          int    `json:"bookno"`
	PublicationName string `json:"publication"`
	Author          `json:"author"`
}

func Init() {
	//tokenAuth = jwtauth.New("HS256", []byte("secret"), nil, jwt.WithAcceptableSkew(30*time.Second))
	//fmt.Println(tokenAuth)
	a1 := Author{
		AuthorFirstName: "Sabbir",
		AuthorLastName:  "Hossain",
		AuthorHomeTown:  "Saidpur",
		AuthorAge:       24,
		AuthorAllBooks:  []string{"book1", "book2"},
	}
	a2 := Author{
		AuthorFirstName: "Shaikat",
		AuthorLastName:  "Shaikat",
		AuthorHomeTown:  "Nilphamari",
		AuthorAge:       44,
		AuthorAllBooks:  []string{"book3"},
	}
	b1 := Book{
		Title:           "book1",
		Author:          a1,
		PublishDate:     "Dec 30, 2011",
		BookNo:          1,
		PublicationName: "PublicationA",
	}
	b2 := Book{
		Title:           "book2",
		Author:          a1,
		PublishDate:     "Mar 23, 2013",
		BookNo:          2,
		PublicationName: "PublicationB",
	}
	b3 := Book{
		Title:           "book3",
		Author:          a2,
		PublishDate:     "Feb 05, 2020",
		BookNo:          3,
		PublicationName: "PublicationA",
	}
	BookList = make(mpBook)
	BookList[b1.BookNo] = b1
	BookList[b2.BookNo] = b2
	BookList[b3.BookNo] = b3

	AuthorList = make(mpAuthor)
	AuthorList[strings.ToLower(a1.AuthorFirstName+a1.AuthorLastName)] = a1
	AuthorList[strings.ToLower(a2.AuthorFirstName+a2.AuthorLastName)] = a2

	User := Credentials{
		Username: "sabbir",
		Password: "pass",
	}
	CredList = make(mpCred)
	CredList[User.Username] = User.Password
	tokenAuth = jwtauth.New(string(jwa.HS256), Secret, nil)

}
func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	BooksJson, err := json.Marshal(BookList)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(BooksJson))
}

func GetSingleBook(w http.ResponseWriter, r *http.Request) {
	bookNoString := chi.URLParam(r, "bookno")
	bookNo, err := strconv.Atoi(bookNoString)
	if err != nil {
		http.Error(w, "Invalid conversion from string to integer... ", http.StatusBadRequest)
	}
	_, exist := BookList[bookNo]
	if exist == false {
		w.Write([]byte("Invalid Book Number "))
		http.Error(w, "Book Not found ", http.StatusBadRequest)
		//println("Invalid Book Number")
		return
	}

	JsonData, _ := json.Marshal(BookList[bookNo])
	w.Write([]byte(JsonData))

}

func GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	AuthorsJson, err := json.Marshal(AuthorList)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(AuthorsJson))
}

func GetSingleAuthor(w http.ResponseWriter, r *http.Request) {
	AuthorName := chi.URLParam(r, "authorname")
	_, exists := AuthorList[AuthorName]
	if exists == false {
		//println("Invalid author name")
		w.Write([]byte("Invalid author name"))
		return
	}
	JsonData, _ := json.Marshal(AuthorList[AuthorName])
	w.Write([]byte(JsonData))

}

func EqualAuthor(a1, a2 Author) bool {
	if a1.AuthorFirstName != a2.AuthorFirstName {
		return false
	}
	if a1.AuthorLastName != a2.AuthorLastName {
		return false
	}
	if a1.AuthorHomeTown != a2.AuthorHomeTown {
		return false
	}
	if a1.AuthorAge != a2.AuthorAge {
		return false
	}
	return true
}

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	AuthorName := chi.URLParam(r, "authorname")
	NeedtoDel := []int{}
	DelAuthor := AuthorList[AuthorName]

	for x, _ := range BookList {
		//fmt.Printf(" THe ans: %T %T", DelAuthor, BookList[x].Author)
		if EqualAuthor(DelAuthor, (BookList[x].Author)) {
			NeedtoDel = append(NeedtoDel, BookList[x].BookNo)
		}
	}
	for x, _ := range NeedtoDel {
		delete(BookList, NeedtoDel[x])
	}
	delete(AuthorList, AuthorName)
	fmt.Println(AuthorList)
	fmt.Println(BookList)
	w.Write([]byte(AuthorName + " was deleted successfully from the author list"))
}
func UpdateAuthorInfo(w http.ResponseWriter, r *http.Request) {
	AuthorName := chi.URLParam(r, "authorname")
	Phone := chi.URLParam(r, "phone")
	var ToUpdate Author
	ToUpdate, exists := AuthorList[AuthorName]
	if exists == false {
		w.Write([]byte("No such author exists"))
		return
	}
	delete(AuthorList, AuthorName)
	ToUpdate.AuthorPhone = Phone
	AuthorList[AuthorName] = ToUpdate

	w.Write([]byte("Phone number of " + AuthorName + " was successfully updated"))
	NeedtoUpdate := []int{}
	UpdateAuthor := AuthorList[AuthorName]
	UpdateAuthor.AuthorPhone = Phone
	for x, _ := range BookList {
		if EqualAuthor(BookList[x].Author, UpdateAuthor) {
			NeedtoUpdate = append(NeedtoUpdate, BookList[x].BookNo)
		}
	}
	for x, _ := range NeedtoUpdate {
		var book Book
		book = BookList[NeedtoUpdate[x]]
		delete(BookList, NeedtoUpdate[x])
		book.Author = UpdateAuthor
		BookList[book.BookNo] = book

	}
	fmt.Println("Author: ", AuthorList)
	fmt.Println(BookList)

}
func AddAuthor(w http.ResponseWriter, r *http.Request) {
	FirstName := chi.URLParam(r, "fn")
	LastName := chi.URLParam(r, "ln")
	Phone := chi.URLParam(r, "phone")
	AuthorName := FirstName + LastName
	_, exists := AuthorList[AuthorName]
	if exists == true {
		w.Write([]byte("Author already exists with this name"))
		return
	}
	var ToAdd Author
	ToAdd.AuthorFirstName = FirstName
	ToAdd.AuthorLastName = LastName
	ToAdd.AuthorPhone = Phone
	AuthorList[FirstName+LastName] = ToAdd
	fmt.Println("AuthorList: ", AuthorList)
	w.Write([]byte(AuthorName + " was successfully added to the author list"))
}

func SearchKey(w http.ResponseWriter, r *http.Request) {
	keyword := chi.URLParam(r, "keyword")
	Books := []Book{}
	Authors := []Author{}
	for x := range BookList {
		if strings.Contains(BookList[x].Title, keyword) {
			Books = append(Books, BookList[x])
		}
	}
	for x := range AuthorList {
		if strings.Contains(strings.ToLower(AuthorList[x].AuthorFirstName+AuthorList[x].AuthorLastName), keyword) {
			Authors = append(Authors, AuthorList[x])
		}
	}
	fmt.Println(Books)
	fmt.Println(Authors)
	AuthorsJson, err := json.Marshal(Authors)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(AuthorsJson))
	BooksJson, err := json.Marshal(Books)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(BooksJson))
}

func main() {
	Init()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Post("/login", Login)
	r.Post("/logout", Logout)
	r.Get("/allbooks", GetAllBooks)
	r.Get("/book/{bookno}", GetSingleBook)
	r.Get("/allauthors", GetAllAuthors)
	r.Get("/author/{authorname}", GetSingleAuthor)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Delete("/author/{authorname}", DeleteAuthor)
		r.Put("/author/{authorname}/{phone}", UpdateAuthorInfo)
		r.Post("/author/{fn}/{ln}/{phone}", AddAuthor)
		r.Get("/search/{keyword}", SearchKey)
	})

	//fmt.Println(BookList)
	log.Fatal(http.ListenAndServe(":8080", r))
}
