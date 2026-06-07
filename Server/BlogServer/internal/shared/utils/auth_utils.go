package utils

import (
	"fmt"
	"log"
	"slices"
)

func CanSubscribe(userID string, roles []string, topic string) bool {

	log.Println("userID: ", userID)
	log.Println("topic: ", topic)

	switch topic {

	case "prices":
		return true

	case fmt.Sprintf("user:%s", userID):
		return hasRole(roles, "user")

	case "admin", "blog_created_admin":
		return hasRole(roles, "admin")

	}

	return false
}

func hasRole(roles []string, role string) bool {
	return slices.Contains(roles, role)
}
