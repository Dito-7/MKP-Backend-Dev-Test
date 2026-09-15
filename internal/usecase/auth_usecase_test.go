package usecase_test

import (
	"context"
	"testing"
	"time"

	"mkp-cinema-ticketing/internal/config"
	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/usecase"
	"mkp-cinema-ticketing/pkg/hash"
)

type mockUserRepo struct {
	users map[string]*entity.User
}

func (m *mockUserRepo) Create(ctx context.Context, u *entity.User) error {
	u.ID = "mock-user-id-123"
	m.users[u.Email] = u
	return nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func TestAuthUsecase_Login(t *testing.T) {
	pwHash, _ := hash.HashPassword("password123")
	mockRepo := &mockUserRepo{
		users: map[string]*entity.User{
			"admin@mkp.com": {
				ID:           "u1",
				FullName:     "Super Admin",
				Email:        "admin@mkp.com",
				PasswordHash: pwHash,
				Role:         entity.RoleAdmin,
				IsActive:     true,
			},
		},
	}

	cfg := &config.Config{
		JWTSecret:      "test-secret",
		JWTExpiresHour: 1 * time.Hour,
	}

	authUC := usecase.NewAuthUsecase(mockRepo, cfg)

	res, err := authUC.Login(context.Background(), entity.LoginRequest{
		Email:    "admin@mkp.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected success, got err: %v", err)
	}
	if res.Token == "" {
		t.Errorf("expected valid token, got empty")
	}
	if res.User.Email != "admin@mkp.com" {
		t.Errorf("expected email admin@mkp.com, got %s", res.User.Email)
	}

	_, err = authUC.Login(context.Background(), entity.LoginRequest{
		Email:    "admin@mkp.com",
		Password: "wrongpassword",
	})
	if err == nil || err != usecase.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}

	_, err = authUC.Login(context.Background(), entity.LoginRequest{
		Email:    "unknown@mkp.com",
		Password: "password123",
	})
	if err == nil || err != usecase.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials for unknown user, got %v", err)
	}
}
