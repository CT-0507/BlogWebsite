package utils

import (
	"fmt"
	"slices"
)

func CanSubscribe(userID string, roles []string, topic string) bool {

	switch topic {

	case "prices":
		return true

	case fmt.Sprintf("user:%s", userID):
		return true

	case "admin", "blog_created_admin":
		return hasRole(roles, "admin")

	}

	return false
}

func hasRole(roles []string, role string) bool {
	return slices.Contains(roles, role)
}
