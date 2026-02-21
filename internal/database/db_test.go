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
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	// Test Mode
	if err := db.SetMode("MANUAL"); err != nil {
		t.Errorf("failed to set mode: %v", err)
	}
	mode, err := db.GetMode()
	if err != nil {
		t.Errorf("failed to get mode: %v", err)
	}
	if mode != "MANUAL" {
		t.Errorf("expected mode MANUAL, got %s", mode)
	}

	// Test Status
	if err := db.SetStatus("note-1", "Blocked"); err != nil {
		t.Errorf("failed to set status: %v", err)
	}
	statuses, err := db.GetStatuses()
	if err != nil {
		t.Errorf("failed to get statuses: %v", err)
	}
	if statuses["note-1"] != "Blocked" {
		t.Errorf("expected status Blocked, got %s", statuses["note-1"])
	}

	// Test Update Status
	if err := db.SetStatus("note-1", "Complete"); err != nil {
		t.Errorf("failed to update status: %v", err)
	}
	statuses, _ = db.GetStatuses()
	if statuses["note-1"] != "Complete" {
		t.Errorf("expected status Complete, got %s", statuses["note-1"])
	}

	// Test Delete Status
	if err := db.DeleteStatus("note-1"); err != nil {
		t.Errorf("failed to delete status: %v", err)
	}
	statuses, _ = db.GetStatuses()
	if _, exists := statuses["note-1"]; exists {
		t.Errorf("expected note-1 to be deleted")
	}
}
