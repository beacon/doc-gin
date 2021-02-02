package main

import (
	"io/ioutil"
	"log"
	"net/http"

	gin "github.com/beacon/doc-gin"
	"github.com/beacon/doc-gin/openapi"
)

// Book book
type Book struct {
	Name      string `json:"name" binding:"required"`
	Author    string `json:"author" binding:"min=1,max=40"`
	ISDN      string `json:"isdn" binding:"required,min=40,max=41"`
	Publisher string `json:"publisher"`
}

// Error resp error
type Error struct {
	Message string `json:"message"`
}

var (
	errConflict   = &Error{Message: "Conflict"}
	errBadRequest = &Error{Message: "Bad request"}
	errNotFound   = &Error{Message: "Not found"}
)

type bookHandler struct {
	books map[string]Book
}

func (b *bookHandler) post(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, errBadRequest)
		return
	}
	if _, exists := b.books[book.Name]; exists {
		c.JSON(http.StatusConflict, errConflict)
		return
	}
	b.books[book.Name] = book
	c.JSON(http.StatusOK, &book)
}

func (b *bookHandler) get(c *gin.Context) {
	name := c.Param("name")
	book, ok := b.books[name]
	if !ok {
		c.JSON(http.StatusNotFound, errNotFound)
		return
	}
	c.JSON(http.StatusOK, &book)
}

func (b *bookHandler) update(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, errBadRequest)
		return
	}
	if _, exists := b.books[book.Name]; !exists {
		c.JSON(http.StatusNotFound, errNotFound)
		return
	}
	b.books[book.Name] = book
	c.JSON(http.StatusOK, &book)
}

func (b *bookHandler) delete(c *gin.Context) {
	name := c.Param("name")
	if _, exists := b.books[name]; !exists {
		c.JSON(http.StatusNotFound, errNotFound)
		return
	}
	delete(b.books, name)
	c.Status(http.StatusNoContent)
}

func main() {
	g := gin.NewEngine(true)
	g.Doc(func(o *openapi.OpenAPI) {
		o.Servers = append(o.Servers, openapi.Server{
			URL: "http://localhost:8080",
		})
	})
	{
		handler := &bookHandler{books: make(map[string]Book)}
		g := g.Group("/books", func(r openapi.Router) {

		})
		g.GET(":name", func(o *openapi.Operation) {

		}, handler.get)
		g.POST("", func(o *openapi.Operation) {
		}, handler.post)
		g.PUT("", func(o *openapi.Operation) {

		}, handler.update)
		g.DELETE(":name", func(o *openapi.Operation) {
			o.Metadata("deleteBook", "Delete a book", "Delete book by name")
		}, handler.delete)
	}
	if doc := g.Doc(nil); doc != nil {
		docYAML, err := doc.YAML()
		if err != nil {
			log.Fatalln(err)
		}
		if err := ioutil.WriteFile("openapi.yaml", docYAML, 0644); err != nil {
			log.Fatalln(err)
		}
	}
}
