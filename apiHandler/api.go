package apiHandler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/sabbir-hossain70/api-server/authHandler"
	"github.com/sabbir-hossain70/api-server/dataHandler"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	BooksJson, err := json.Marshal(dataHandler.BookList)
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
	_, exist := dataHandler.BookList[bookNo]
	if exist == false {
		w.Write([]byte("Invalid Book Number "))
		http.Error(w, "Book Not found ", http.StatusBadRequest)
		//println("Invalid Book Number")
		return
	}

	JsonData, _ := json.Marshal(dataHandler.BookList[bookNo])
	w.Write([]byte(JsonData))

}

func GetAllAuthors(w http.ResponseWriter, r *http.Request) {
	AuthorsJson, err := json.Marshal(dataHandler.AuthorList)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(AuthorsJson))
}

func GetSingleAuthor(w http.ResponseWriter, r *http.Request) {
	AuthorName := chi.URLParam(r, "authorname")
	_, exists := dataHandler.AuthorList[AuthorName]
	if exists == false {
		//println("Invalid author name")
		w.Write([]byte("Invalid author name"))
		return
	}
	JsonData, _ := json.Marshal(dataHandler.AuthorList[AuthorName])
	w.Write([]byte(JsonData))

}

func EqualAuthor(a1, a2 dataHandler.Author) bool {
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
	DelAuthor := dataHandler.AuthorList[AuthorName]

	for x, _ := range dataHandler.BookList {
		//fmt.Printf(" THe ans: %T %T", DelAuthor, BookList[x].Author)
		if EqualAuthor(DelAuthor, (dataHandler.BookList[x].Author)) {
			NeedtoDel = append(NeedtoDel, dataHandler.BookList[x].BookNo)
		}
	}
	for x, _ := range NeedtoDel {
		delete(dataHandler.BookList, NeedtoDel[x])
	}
	delete(dataHandler.AuthorList, AuthorName)
	fmt.Println(dataHandler.AuthorList)
	fmt.Println(dataHandler.BookList)
	w.Write([]byte(AuthorName + " was deleted successfully from the author list"))
}
func UpdateAuthorInfo(w http.ResponseWriter, r *http.Request) {
	AuthorName := chi.URLParam(r, "authorname")
	Phone := chi.URLParam(r, "phone")
	var ToUpdate dataHandler.Author
	ToUpdate, exists := dataHandler.AuthorList[AuthorName]
	if exists == false {
		w.Write([]byte("No such author exists"))
		return
	}
	delete(dataHandler.AuthorList, AuthorName)
	ToUpdate.AuthorPhone = Phone
	dataHandler.AuthorList[AuthorName] = ToUpdate

	w.Write([]byte("Phone number of " + AuthorName + " was successfully updated"))
	NeedtoUpdate := []int{}
	UpdateAuthor := dataHandler.AuthorList[AuthorName]
	UpdateAuthor.AuthorPhone = Phone
	for x, _ := range dataHandler.BookList {
		if EqualAuthor(dataHandler.BookList[x].Author, UpdateAuthor) {
			NeedtoUpdate = append(NeedtoUpdate, dataHandler.BookList[x].BookNo)
		}
	}
	for x, _ := range NeedtoUpdate {
		var book dataHandler.Book
		book = dataHandler.BookList[NeedtoUpdate[x]]
		delete(dataHandler.BookList, NeedtoUpdate[x])
		book.Author = UpdateAuthor
		dataHandler.BookList[book.BookNo] = book

	}
	fmt.Println("Author: ", dataHandler.AuthorList)
	fmt.Println(dataHandler.BookList)

}
func AddAuthor(w http.ResponseWriter, r *http.Request) {
	FirstName := chi.URLParam(r, "fn")
	LastName := chi.URLParam(r, "ln")
	Phone := chi.URLParam(r, "phone")
	AuthorName := FirstName + LastName
	_, exists := dataHandler.AuthorList[AuthorName]
	if exists == true {
		w.Write([]byte("Author already exists with this name"))
		return
	}
	var ToAdd dataHandler.Author
	ToAdd.AuthorFirstName = FirstName
	ToAdd.AuthorLastName = LastName
	ToAdd.AuthorPhone = Phone
	dataHandler.AuthorList[FirstName+LastName] = ToAdd
	fmt.Println("AuthorList: ", dataHandler.AuthorList)
	w.Write([]byte(AuthorName + " was successfully added to the author list"))
}

func SearchKey(w http.ResponseWriter, r *http.Request) {
	keyword := chi.URLParam(r, "keyword")
	Books := []dataHandler.Book{}
	Authors := []dataHandler.Author{}
	for x := range dataHandler.BookList {
		if strings.Contains(dataHandler.BookList[x].Title, keyword) {
			Books = append(Books, dataHandler.BookList[x])
		}
	}
	for x := range dataHandler.AuthorList {
		if strings.Contains(strings.ToLower(dataHandler.AuthorList[x].AuthorFirstName+dataHandler.AuthorList[x].AuthorLastName), keyword) {
			Authors = append(Authors, dataHandler.AuthorList[x])
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

func RunServer(Port string) {
	dataHandler.Init()
	authHandler.InitToken()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)
	r.Get("/allbooks", GetAllBooks)
	r.Get("/book/{bookno}", GetSingleBook)
	r.Get("/allauthors", GetAllAuthors)
	r.Get("/author/{authorname}", GetSingleAuthor)
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(authHandler.TokenAuth))
		r.Use(jwtauth.Authenticator(authHandler.TokenAuth))
		r.Delete("/author/{authorname}", DeleteAuthor)
		r.Put("/author/{authorname}/{phone}", UpdateAuthorInfo)
		r.Post("/author/{fn}/{ln}/{phone}", AddAuthor)
		r.Get("/search/{keyword}", SearchKey)
	})

	fmt.Println("Port :::: ", Port)
	log.Fatal(http.ListenAndServe(":"+Port, r))
}
