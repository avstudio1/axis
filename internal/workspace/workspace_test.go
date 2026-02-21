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

/*
File: internal/workspace/workspace_test.go
Description: Unit tests for the Workspace service. Validates service initialization,
registry item consolidation, and Google API service wrapping logic.
*/
package workspace

import (
	"testing"

	admin "google.golang.org/api/admin/directory/v1"
	docs "google.golang.org/api/docs/v1"
	drive "google.golang.org/api/drive/v3"
	keep "google.golang.org/api/keep/v1"
	sheets "google.golang.org/api/sheets/v4"
)

func TestNewService(t *testing.T) {
	adminSvc := &admin.Service{}
	keepSvc := &keep.Service{}
	docsSvc := &docs.Service{}
	sheetsSvc := &sheets.Service{}
	driveSvc := &drive.Service{}

	ws := NewService(adminSvc, keepSvc, docsSvc, sheetsSvc, driveSvc)

	if ws.adminService != adminSvc {
		t.Error("Admin service not correctly assigned")
	}
	if ws.keepService != keepSvc {
		t.Error("Keep service not correctly assigned")
	}
	if ws.docsService != docsSvc {
		t.Error("Docs service not correctly assigned")
	}
	if ws.sheetsService != sheetsSvc {
		t.Error("Sheets service not correctly assigned")
	}
	if ws.driveService != driveSvc {
		t.Error("Drive service not correctly assigned")
	}
}
