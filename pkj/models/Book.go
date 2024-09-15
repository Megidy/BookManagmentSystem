package models

import (
	"database/sql"
)

type Book struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn"`
	Avaibe_Copies int    `json:"avaible_copies"`
}

func CreateBook(book *Book) error {
	_, err := db.Exec("insert into books(title,author,isbn,avaible_copies) values(?,?,?,?)",
		book.Title, book.Author, book.Isbn, book.Avaibe_Copies)
	if err != nil {

		return err
	}
	return nil
}
func IsCreated(NewBook Book) (bool, error) {
	var book Book
	row := db.QueryRow("select (isbn) from books where isbn=?", NewBook.Isbn)
	err := row.Scan(&book.Isbn)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found, meaning the book is not created
			return false, nil
		} else {
			return true, err
		}
	}
	if book.Isbn == NewBook.Isbn {
		return true, nil
	}
	return false, nil
}
func GetAllBooks() ([]Book, error) {
	var books []Book
	query, err := db.Query("select * from books")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	for query.Next() {
		var book Book
		err := query.Scan(&book.Id, &book.Title, &book.Author, &book.Isbn, &book.Avaibe_Copies)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
func FindBook(Id int) (*Book, error) {
	var book Book
	row := db.QueryRow("select title,author,isbn,avaible_copies from books where id=?", Id)
	err := row.Scan(&book.Title, &book.Author, &book.Isbn, &book.Avaibe_Copies)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}
func DeleteBook(Id int) error {
	_, err := db.Exec("delete from books where id=?", Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
func UpdateBook(book *Book) error {
	_, err := db.Exec("update books set title=?,author=?,isbn=?,avaible_copies=? where id=?",
		book.Title, book.Author, book.Isbn, book.Avaibe_Copies, book.Id)
	if err != nil {
		return err
	}
	return nil
}

func TakeBook(id int, amount int) error {
	_, err := db.Exec("update books set avaible_copies=avaible_copies+? where id =?", amount, id)
	if err != nil {

		return err
	}
	return nil

}
func GiveBook(id int, amount int) error {
	_, err := db.Exec("update books set avaible_copies=avaible_copies-? where avaible_copies>=? and id=?",
		amount, amount, id)
	if err != nil {

		return err
	}
	return nil
}
