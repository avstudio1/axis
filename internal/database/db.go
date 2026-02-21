/*
MIT License

Copyright (c) 2026 Justin Andrew Wood

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package database

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

// DB wraps the sql.DB connection and provides state-specific methods.
type DB struct {
	db *sql.DB
	mu sync.RWMutex
}

// NewDB initializes a new SQLite database connection and runs migrations.
func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	d := &DB{db: db}
	if err := d.init(); err != nil {
		db.Close()
		return nil, err
	}

	return d, nil
}

// init creates the necessary tables if they don't exist.
func (d *DB) init() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS app_state (
			key TEXT PRIMARY KEY,
			value TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS item_statuses (
			id TEXT PRIMARY KEY,
			status TEXT
		);`,
	}

	for _, q := range queries {
		if _, err := d.db.Exec(q); err != nil {
			return fmt.Errorf("failed to initialize schema: %w", err)
		}
	}

	return nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.db.Close()
}

// SetMode updates the operational mode in the database.
func (d *DB) SetMode(mode string) error {
	_, err := d.db.Exec(`INSERT INTO app_state (key, value) VALUES ('mode', ?) 
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`, mode)
	return err
}

// GetMode retrieves the operational mode from the database.
func (d *DB) GetMode() (string, error) {
	var mode string
	err := d.db.QueryRow(`SELECT value FROM app_state WHERE key = 'mode'`).Scan(&mode)
	if err == sql.ErrNoRows {
		return "AUTO", nil
	}
	return mode, err
}

// SetStatus updates the status for a given item ID.
func (d *DB) SetStatus(id, status string) error {
	_, err := d.db.Exec(`INSERT INTO item_statuses (id, status) VALUES (?, ?) 
		ON CONFLICT(id) DO UPDATE SET status = excluded.status`, id, status)
	return err
}

// GetStatuses retrieves all item statuses as a map.
func (d *DB) GetStatuses() (map[string]string, error) {
	rows, err := d.db.Query(`SELECT id, status FROM item_statuses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	statuses := make(map[string]string)
	for rows.Next() {
		var id, status string
		if err := rows.Scan(&id, &status); err != nil {
			return nil, err
		}
		statuses[id] = status
	}
	return statuses, nil
}

// DeleteStatus removes a status entry for a given ID.
func (d *DB) DeleteStatus(id string) error {
	_, err := d.db.Exec(`DELETE FROM item_statuses WHERE id = ?`, id)
	return err
}
