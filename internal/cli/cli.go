package cli

import (
	"fmt"
	"gobooks/internal/service"
	"os"
	"strconv"
	"time"
)

type BookCLI struct {
	service *service.BookService
}

func NewBookCLI(service *service.BookService) *BookCLI {
	return &BookCLI{service: service}
}

func (cli *BookCLI) Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: books <command> [arguments]")
		return
	}

	command := os.Args[1]
	switch command {
	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: books search <book title>")
			return
		}
		bookTitle := os.Args[2]
		cli.searchBooks(bookTitle)
	case "simulate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: books simulate <book_id> <book_id> <book_id> ...")
			return
		}
		bookIDS := os.Args[2:]
		cli.simuleReading(bookIDS)
	default:
		fmt.Println("Invalid command")
	}
}

func (cli *BookCLI) searchBooks(bookTitle string) {
	books, err := cli.service.SearchBooksByTitle(bookTitle)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, book := range books {
		fmt.Printf("ID: %d, Title: %s, Author: %s, Genre: %s\n", book.ID, book.Title, book.Author, book.Genre)
	}
}

func (cli *BookCLI) simuleReading(bookIDs []string) {
	var bookIDS []int
	for _, idString := range bookIDs {
		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println("Invalid book ID:", idString)
			continue
		}
		bookIDS = append(bookIDS, id)
	}
	responses := cli.service.SimulateMultipleReadings(bookIDS, 5*time.Second)
	for _, response := range responses {
		fmt.Println(response)
	}
}
