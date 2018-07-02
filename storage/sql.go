package storage

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type SqlConfig struct {
	Endpoint string `json:"endpoint"`
	UserName string `json:"username"`
	DataBase string `json:"database"`
	Password string `json:"password"`
}

type SqlClient struct {
	client *sql.DB
}

type ShareData struct {
	ID        string
	Address   string
	Nonce     string
	HashNonce string
	Score     string
}

func (s *SqlClient) Ping() error {
	return s.client.Ping()
}

func (s *SqlClient) CloseConnection() error {
	return s.client.Close()
}

func (s *SqlClient) GetAllShares(pageStart, pageLength int) ([]ShareData, error) {
	var rows *sql.Rows
	var err error
	//TODO: poor way of doing defaults.. redo this
	if pageStart <= 0 && pageLength <= 0 {
		rows, err = s.client.Query("select id, address, nonce, hash_nonce, score from shares")
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = s.client.Query("select id, address, nonce, hash_nonce, score from shares order by id desc limit ?,?", pageStart, pageLength)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	var result []ShareData
	for rows.Next() {
		var s ShareData
		err := rows.Scan(&s.ID, &s.Address, &s.Nonce, &s.HashNonce, &s.Score)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	return result, nil
}

func (s *SqlClient) InsertShare(address, nonce, hashNonce, score string) (sql.Result, error) {
	stmt, err := s.client.Prepare("Insert into shares(address, nonce, hash_nonce, score) VALUES(?,?,?,?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(address, nonce, hashNonce, score)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewSqlClient(cfg *SqlConfig) (*SqlClient, error) {
	// debug
	// log.Printf("Connection string %s", cfg.UserName+":"+cfg.Password+"@tcp("+cfg.Endpoint+")/"+cfg.DataBase)
	db, err := sql.Open("mysql", cfg.UserName+":"+cfg.Password+"@tcp("+cfg.Endpoint+")/"+cfg.DataBase)
	return &SqlClient{client: db}, err
}
