package workspace

import (
	"context"

	admin "google.golang.org/api/admin/directory/v1"
)

type googleService struct {
	admin *admin.Service
}

func NewService(adminSvc *admin.Service) Service {
	return &googleService{
		admin: adminSvc,
	}
}

func (s *googleService) GetUser(email string) (*User, error) {
	ctx := context.Background()

	u, err := s.admin.Users.Get(email).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &User{
		Email: u.PrimaryEmail,
		Name:  u.Name.FullName,
	}, nil
}

func (s *googleService) SuspendUser(email string) error {
	ctx := context.Background()

	u := &admin.User{Suspended: true}
	_, err := s.admin.Users.Update(email, u).Context(ctx).Do()
	return err
}
