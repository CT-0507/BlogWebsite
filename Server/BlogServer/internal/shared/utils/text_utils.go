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

func StringToUUID(s string) *uuid.UUID {
	uuid, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &uuid
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Truncate(s string, max int, withEllipsis bool) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	result := string(runes[:max])
	if withEllipsis {
		result += "..."
	}
	return result
}

func DerefTextNullable(text pgtype.Text) *string {
	if text.Valid {
		return &text.String
	}
	return nil
}

func Deref[T any](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}

func StringPtr(s string) *string {
	return &s
}
