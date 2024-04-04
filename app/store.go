package app

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

type Store interface {
	Add(task string, priority string) (Todo, error)
	Remove(todoID int) error
	Done(todoID int) error
	GetAll() ([]Todo, error)
	Update(taskID int, status Status) error
}

type DefaultStore struct {
	db *sqlx.DB
}

func NewDefaultStore() (DefaultStore, error) {
	homeDir, _ := os.UserHomeDir()
	dbPath := fmt.Sprintf("%s/.todos.db", homeDir)

	_, err := os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		file, createErr := os.Create(dbPath)
		if createErr != nil {
			return DefaultStore{}, fmt.Errorf("could not create db in home dir %v", createErr)
		}
		file.Close()
	}

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("error opening db %v", err)
		return DefaultStore{}, err
	}

	ds := DefaultStore{db}

	initErr := ds.init()
	if initErr != nil {
		fmt.Printf("error creating schema %v", err)
		return DefaultStore{}, err
	}

	return ds, nil
}

func (d DefaultStore) Add(task string, priority string) (Todo, error) {
	result, err := d.db.Exec(
		`INSERT INTO todos (task, status, priority) VALUES (?, ?, ?)`,
		task, StatusPending, priority)
	if err != nil {
		fmt.Printf("error updating status %v", err)
		return Todo{}, err
	}

	taskId, _ := result.LastInsertId()

	return Todo{
		ID:       int(taskId),
		Task:     task,
		Status:   StatusPending,
		Priority: priority,
	}, nil
}

func (d DefaultStore) Remove(todoID int) error {
	_, err := d.db.Query(
		`DELETE FROM todos WHERE id = ?`, todoID)
	if err != nil {
		fmt.Printf("error deleting task %v", err)
		return err
	}

	return nil
}

func (d DefaultStore) Done(todoID int) error {
	_, err := d.db.Query(
		`UPDATE todos SET status = ? , updated_at = ? WHERE id = ?`,
		StatusDone, time.Now(), todoID)
	if err != nil {
		fmt.Printf("error updating status %v", err)
		return err
	}

	return nil
}

func (d DefaultStore) Update(todoID int, status Status) error {
	_, err := d.db.Exec(
		`UPDATE todos SET status = ? , updated_at = ? WHERE id = ?`,
		string(status), time.Now(), todoID)
	if err != nil {
		fmt.Printf("error updating status %v", err)
		return err
	}

	return nil
}

func (d DefaultStore) GetAll() ([]Todo, error) {
	todos := make([]Todo, 0)

	dbErr := d.db.Select(&todos, "SELECT * FROM todos ORDER BY created_at DESC")
	if dbErr != nil {
		fmt.Printf("error fetching todos from db %v", dbErr)
		return []Todo{}, dbErr
	}

	return todos, nil
}

func (d DefaultStore) init() error {
	exists, dbErr := d.tableExists()
	if dbErr != nil {
		fmt.Printf("error checking if table exists %v", dbErr)
		return dbErr
	}

	if exists {
		return nil
	}

	return d.createTable()
}

func (d DefaultStore) tableExists() (bool, error) {
	rows, err := d.db.Query(
		`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`,
		"todos")
	if err != nil {
		fmt.Printf("error updating status %v", err)
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		return true, nil
	}

	return false, nil
}

func (d DefaultStore) createTable() error {
	schema :=
		`CREATE TABLE todos (
    	 id INTEGER PRIMARY KEY ,
    	 task VARCHAR NOT NULL ,
    	 status VARCHAR NOT NULL DEFAULT 'Pending',
    	 priority VARCHAR NOT NULL DEFAULT 'High',
    	 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	_, err := d.db.Exec(schema)
	if err != nil {
		fmt.Printf("error creating todos schema %v", err)
		return err
	}

	return nil
}
