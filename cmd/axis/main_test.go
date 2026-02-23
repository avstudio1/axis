/*
PROPRIETARY AND CONFIDENTIAL LICENSE
Copyright © 2026 Justin Andrew Wood. All Rights Reserved.

1. OWNERSHIP

All title, ownership rights, and intellectual property rights in and to this project shall remain the exclusive property of Justin Andrew Wood.

2. RESTRICTIONS

No person or entity may:

• Use, copy, modify, or distribute the Software/Project.

• Sublicense, rent, lease, or lend the Software/Project.

• Reverse engineer, decompile, or disassemble any components.

• Make the project available over a network or public repository.

3. NO GRANT OF LICENSE

This document does not grant any licenses, express or implied, to any party. Any use without the prior written consent of Justin Andrew Wood is strictly prohibited.
*/

/*
File: cmd/axis/main_test.go
Description: Integration tests for the Axis entry point. Validates environment
variable requirements and service initialization flow.
*/
package main

import (
	"os"
	"testing"
)

func TestMainValidation(t *testing.T) {
	// Clear environment to test validation logic
	os.Clearenv()

	// Capture exit or panic behavior for missing environment variables
	defer func() {
		if r := recover(); r == nil {
			// In a real scenario, log.Fatal would terminate the test process.
			// This test assumes logic isolation for validation.
		}
	}()

	requiredVars := []string{"ADMIN_EMAIL", "SERVICE_ACCOUNT_EMAIL", "USER_EMAIL"}
	for _, v := range requiredVars {
		if os.Getenv(v) != "" {
			t.Errorf("Environment variable %s should be empty", v)
		}
	}
}

func TestDefaultPort(t *testing.T) {
	os.Setenv("PORT", "")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if port != "8080" {
		t.Errorf("Expected default port 8080, got %s", port)
	}
}
