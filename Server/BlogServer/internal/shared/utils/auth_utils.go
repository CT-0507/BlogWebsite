package utils

import (
	"fmt"
	"slices"
)

func CanSubscribe(userID string, roles []string, topic string) bool {

	switch {

	case topic == "prices":
		return true

	case topic == fmt.Sprintf("user:%s", userID):
		return true

	case topic == "admin":
	case topic == "blog_created_admin":
		return hasRole(roles, "admin")

	}

	return false
}

func hasRole(roles []string, role string) bool {
	return slices.Contains(roles, role)
}
