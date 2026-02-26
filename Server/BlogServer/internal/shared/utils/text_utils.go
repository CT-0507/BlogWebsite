package utils

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func TextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func TimePointer(timeP *pgtype.Timestamptz) *time.Time {
	if timeP.Valid {
		return &timeP.Time
	}
	return nil
}

func UUIDPtr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	if s == "" {
		return nil
	}
	return &s
}

func UUIDFromStringPtr(s *string) (pgtype.UUID, error) {
	if s == nil {
		return pgtype.UUID{Valid: false}, nil
	}

	u, err := uuid.Parse(*s)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}, nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
