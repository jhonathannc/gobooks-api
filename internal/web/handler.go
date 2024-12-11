package web

import (
	"encoding/json"
	"gobooks/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type BookHandlers struct {
	service *service.BookService
}

func NewBookHandlers(service *service.BookService) *BookHandlers {
	return &BookHandlers{service: service}
}

func (h *BookHandlers) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book service.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.service.CreateBook(&book)
	if err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}
	book, err := h.service.GetBookByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}
	var book service.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	book.ID = id

	if err = h.service.UpdateBook(&book); err != nil {
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteBook(id); err != nil {
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func SetupBookHandlers(router *http.ServeMux, bookHandler *BookHandlers) {
	// Rota para /books (GET e POST)
	router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetBooks(w, r)
		case http.MethodPost:
			bookHandler.CreateBook(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Rota para /books/{id} (GET, PUT, DELETE)
	router.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		// Extrair o {id} da URL
		// Ex: /books/123 -> partes = ["books", "123"]
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		idStr := parts[1]
		_, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid book ID", http.StatusBadRequest)
			return
		}

		// Ajusta o Query String para simular ?id={id}, caso handlers usem isso
		q := r.URL.Query()
		q.Set("id", idStr)
		r.URL.RawQuery = q.Encode()

		switch r.Method {
		case http.MethodGet:
			bookHandler.GetBookByID(w, r)
		case http.MethodPut:
			bookHandler.UpdateBook(w, r)
		case http.MethodDelete:
			bookHandler.DeleteBook(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
