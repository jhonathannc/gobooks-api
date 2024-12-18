package main

import (
	"database/sql"
	"fmt"
	"gobooks/internal/cli"
	"gobooks/internal/service"
	"gobooks/internal/web"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	bookService := service.NewBookService(db)
	bookHandler := web.NewBookHandlers(bookService)

	if len(os.Args) > 1 && (os.Args[1] == "search" || os.Args[1] == "simulate") {
		bookCli := cli.NewBookCLI(bookService)
		bookCli.Run()
		return
	}

	router := http.NewServeMux()
	web.SetupBookHandlers(router, bookHandler)

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("Books Api"))
	// })

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
