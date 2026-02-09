package workspace

import "errors"

// MockService is a manual mock for the Service interface.
type MockService struct {
	OnGetUser func(email string) (*User, error)
}

// GetUser calls the function assigned to OnGetUser.
func (m *MockService) GetUser(email string) (*User, error) {
	if m.OnGetUser != nil {
		return m.OnGetUser(email)
	}
	return nil, errors.New("OnGetUser not implemented")
}
