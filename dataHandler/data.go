package dataHandler

import (
	"strings"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type mpCred map[string]string

var CredList mpCred

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
	return
}
