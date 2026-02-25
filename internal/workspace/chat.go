// Copyright (c) 2026 Justin Andrew Wood. All rights reserved.
// This software is licensed under the AGPL-3.0.
// Commercial licensing is available at echosh-labs.com.
/*
File: internal/workspace/chat.go
Description: Google Chat API integration for Axis Mundi. Handles sending direct
messages for telemetry and notifications.
*/
package workspace

import (
	"fmt"

	chat "google.golang.org/api/chat/v1"
)

// SendDirectMessage sends a direct message to the specified email address.
// Resolves the space or creates a DM and posts the message text.
func (s *Service) SendDirectMessage(email string, text string) error {
	if s.chatService == nil {
		return fmt.Errorf("chat service is not initialized")
	}

	// Set up the Direct Message space.
	// If the DM already exists, this method returns the existing space.
	req := &chat.SetUpSpaceRequest{
		Space: &chat.Space{
			SpaceType: "DIRECT_MESSAGE",
		},
		Memberships: []*chat.Membership{
			{
				Member: &chat.User{
					Name: "users/" + email,
					Type: "HUMAN",
				},
			},
		},
	}

	space, err := s.chatService.Spaces.Setup(req).Do()
	if err != nil {
		return fmt.Errorf("failed to setup chat space for %s: %w", email, err)
	}

	// Send the message
	msg := &chat.Message{
		Text: text,
	}

	_, err = s.chatService.Spaces.Messages.Create(space.Name, msg).Do()
	if err != nil {
		return fmt.Errorf("failed to send chat message to %s: %w", email, err)
	}

	return nil
}
