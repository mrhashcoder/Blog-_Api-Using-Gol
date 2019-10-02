package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/gorilla/mux"
)

type blog struct {
	Id          string `json : "Id"`
	Author      string `json : "Author"`
	Title       string `json : "Title"`
	Description string `json : "Description"`
	Content     string `json : "Content"`
}

//												//
// FUNCTIONS FOR DATABASES ACCESESS START HERE//
//												//

func insertBlog(myDb *db.DB, b blog) {
	fmt.Println("Insert blog is invoked")
	blogcol := myDb.Use("Blogs")
	_, err := blogcol.Insert(map[string]interface{}{
		"Id":          b.Id,
		"Author":      b.Author,
		"Title":       b.Title,
		"Description": b.Description,
		"Content":     b.Content,
	})
	if err != nil {
		panic(err)
	}
}

func deleteBlog(myDb *db.DB, Id string) {
	fmt.Println("Delete blog is invoked")
	blogs := myDb.Use("Blogs")
	var docid int
	blogs.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		var raw map[string]interface{}
		json.Unmarshal(docContent, &raw)
		if raw["Id"].(string) == Id {
			docid = id
		}
		return true
		return false
	})
	err := blogs.Delete(docid)
	if err != nil {
		panic(err)
	}
}

func retriveAll(mydb *db.DB) []blog {
	blogs := mydb.Use("Blogs")
	var allblogs []blog
	blogs.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		var raw map[string]interface{}
		json.Unmarshal(docContent, &raw)
		var b blog
		b.Author = raw["Author"].(string)
		b.Content = raw["Content"].(string)
		b.Description = raw["Description"].(string)
		b.Id = raw["Id"].(string)
		b.Title = raw["Title"].(string)
		allblogs = append(allblogs, b)
		return true
		return false
	})
	return allblogs
}

func retriveSingle(myDb *db.DB, Author string) []blog {
	blogs := myDb.Use("Blogs")
	var blogsByAuthor []blog
	blogs.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		var raw map[string]interface{}
		json.Unmarshal(docContent, &raw)
		if raw["Author"].(string) == Author {

			var b blog
			b.Author = raw["Author"].(string)
			b.Content = raw["Content"].(string)
			b.Description = raw["Description"].(string)
			b.Id = raw["Id"].(string)
			b.Title = raw["Title"].(string)
			blogsByAuthor = append(blogsByAuthor, b)

		}
		return true
		return false
	})
	return blogsByAuthor
}

func getDocId(Id string, myDb *db.DB) int {
	blogs := myDb.Use("Blogs")
	var docid int
	blogs.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		var raw map[string]interface{}
		json.Unmarshal(docContent, &raw)
		if raw["Id"].(string) == Id {
			docid = id
		}
		return true  // move on to the next document OR
		return false // do not move on to the next document
	})

	return docid
}

//											//
//FUNCTION FOR DATABASES ACCSESS ENDS HERE//
//											//

//											//
//FUNCTION FOR API ROUTES STARTS HERE//
//											//

func handleRequests(myDb *db.DB) {

	hashRouter := mux.NewRouter().StrictSlash(true)

	hashRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello world")
		fmt.Fprintln(w, "Welcome to Blog api ")
		fmt.Fprintln(w, "created during learing of go and api creation")
		fmt.Fprintln(w, "1. For viewing available blogs in database put /blog ")
		fmt.Fprintln(w, "2. For viewing blogs of an specific user put /blog/author_name ")
		fmt.Fprintln(w, "3. For deleting any blog from database put /deleteBlog/blog_id")
		fmt.Fprintln(w, "4. For Inserting a blog in database send a http post request using createBlogrequest.go ")
		fmt.Fprintln(w, "5. For updating available blog in databases send a http post request using updateBlogRequest.go")

	})

	hashRouter.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		allblogs := retriveAll(myDb)
		blogAddress := &allblogs
		e, err := json.Marshal(blogAddress)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(w, string(e))
	})
	hashRouter.HandleFunc("/blog/{author}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		authorName := vars["author"]
		var blogsByAuthor []blog
		blogsByAuthor = retriveSingle(myDb, authorName)
		blogsByAuthorAddress := &blogsByAuthor
		e, err := json.Marshal(blogsByAuthorAddress)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(w, string(e))

	})

	hashRouter.HandleFunc("/deleteBlog/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		delId := vars["id"]

		deleteBlog(myDb, delId)
		fmt.Fprintf(w, "deleted blog of id  %d ", delId)
	})

	hashRouter.HandleFunc("/createBlog", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit Endpoint of Create Blog")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var raw map[string]interface{}
		json.Unmarshal(reqBody, &raw)

		var newBlog blog
		newBlog.Id = raw["Id"].(string)
		newBlog.Author = raw["Author"].(string)
		newBlog.Content = raw["content"].(string)
		newBlog.Description = raw["Description"].(string)
		newBlog.Title = raw["Title"].(string)

		insertBlog(myDb, newBlog)

	}).Methods("Post")

	hashRouter.HandleFunc("/updateBlog", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the end point of update blog")
		reqBody, _ := ioutil.ReadAll(r.Body)
		var raw map[string]interface{}
		json.Unmarshal(reqBody, &raw)
		docid := getDocId(raw["Id"].(string), myDb)
		blogs := myDb.Use("Blogs")
		err := blogs.Update(docid, raw)
		if err != nil {
			panic(err)
		}
	}).Methods("Post")

	log.Fatal(http.ListenAndServe(":1234", hashRouter))
}

//											//
//FUNCTION FOR API ROUTES ENDS HERE//
//											//

func main() {
	myDbDir := "/MyDatabase"
	os.RemoveAll(myDbDir)
	defer os.RemoveAll(myDbDir)

	myDb, err := db.OpenDB(myDbDir)
	if err != nil {

		panic(err)
	}

	err1 := myDb.Create("Blogs")
	if err1 != nil {
		panic(err1)
	}

	blogs := myDb.Use("Blogs")
	err2 := blogs.Index([]string{"Id", "Author", "Title", "Description", "Content"})
	if err2 != nil {
		panic(err2)
	}
	var b1 = blog{"1", "Abhishek", "Blog api", "leraning go for blog api and gawds", "kuch samj nhi aa rha yaar"}

	var b2 = blog{"2", "Rekhansh", "help me", "Mere ko sab aata", "gyan mt do yaar"}

	var b3 = blog{"3", "Keshav", "Blog web", "leraning go for blog api and gawds", "kuch samj nhi aa rha yaar"}

	insertBlog(myDb, b1)
	insertBlog(myDb, b2)
	insertBlog(myDb, b3)
	insertBlog(myDb, b2)
	insertBlog(myDb, b1)

	handleRequests(myDb)

}
