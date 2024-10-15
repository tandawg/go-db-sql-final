package main

import (
	"database/sql"
	"errors"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	query := `INSERT INTO parcels (client, status, address, created_at) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := `SELECT id, client, status, address, created_at FROM parcels WHERE id = ?`
	row := s.db.QueryRow(query, number)

	var parcel Parcel
	err := row.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return parcel, errors.New("посылка не найдена")
		}
		return parcel, err
	}

	return parcel, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT id, client, status, address, created_at FROM parcels WHERE client = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var parcel Parcel
		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, parcel)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := `UPDATE parcels SET status = ? WHERE id = ?`
	_, err := s.db.Exec(query, status, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return errors.New("адрес можно изменить только для посылок со статусом 'registered'")
	}

	query := `UPDATE parcels SET address = ? WHERE id = ?`
	_, err = s.db.Exec(query, address, number)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}

	if parcel.Status != ParcelStatusRegistered {
		return errors.New("удалить можно только посылку со статусом 'registered'")
	}

	query := `DELETE FROM parcels WHERE id = ?`
	_, err = s.db.Exec(query, number)
	if err != nil {
		return err
	}

	return nil
}
