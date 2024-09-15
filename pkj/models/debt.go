package models

import "database/sql"

type Debt struct {
	User_id int `json:"user_id"`
	Book_id int `json:"book_id"`
	Amount  int `json:"amount"`
}

func GetAllDebts() (*[]Debt, error) {
	var debts []Debt
	query, err := db.Query("select * from debts")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	for query.Next() {
		var debt Debt
		err := query.Scan(&debt.User_id, &debt.Book_id, &debt.Amount)
		if err != nil {
			return nil, err
		}
		debts = append(debts, debt)
	}
	return &debts, nil
}
func GetDebt(user *User, NewDebt *Debt) (*Debt, error) {
	var debt Debt
	row := db.QueryRow("select * from debts where user_id=? and book_id=?", user.Id, NewDebt.Book_id)
	err := row.Scan(&debt.User_id, &debt.Book_id, &debt.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}

	}
	return &debt, nil

}
func GetAllUsersDebts(user *User) (*[]Debt, error) {
	var debts []Debt
	query, err := db.Query("select * from debts where user_id=?", user.Id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	for query.Next() {
		var debt Debt
		err := query.Scan(&debt.User_id, &debt.Book_id, &debt.Amount)
		if err != nil {
			return nil, err
		}
		debts = append(debts, debt)
	}
	return &debts, nil
}
func AlreadyExist(user *User, NewDebt *Debt) (bool, error) {
	var debt Debt
	row := db.QueryRow("Select * from debts where user_id=? and book_id=?", user.Id, NewDebt.Book_id)
	err := row.Scan(&debt.User_id, &debt.Book_id, &debt.Amount)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil

		} else {
			return true, err
		}
	}

	NewDebt.User_id = user.Id
	if NewDebt.Book_id == debt.Book_id && NewDebt.User_id == debt.User_id {
		return true, nil
	}
	return false, nil

}
func CreateDebt(user_id int, debt *Debt) error {
	_, err := db.Exec("insert into debts values(?,?,?)", user_id, debt.Book_id, debt.Amount)
	if err != nil {
		return err
	}
	return nil

}

func AddAmountForDebt(user *User, debt *Debt) error {

	_, err := db.Exec("Update debts set amount=amount+? where book_id=? and user_id=?",
		debt.Amount, debt.Book_id, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func SubAmountFromDebt(user *User, debt *Debt) error {
	_, err := db.Exec("update debts set amount=amount-? where book_id=? and user_id=?",
		debt.Amount, debt.Book_id, debt.User_id)
	if err != nil {
		return err
	}
	return nil
}
func ResetDebt(user *User, debt *Debt) error {
	_, err := db.Exec("delete from debts where user_id=? and book_id=?", user.Id, debt.Book_id)
	if err != nil {

		return err
	}

	return nil

}
