package models

import "fmt"

// Role type definition
type Role int

// Enum values for Role using iota
const (
	ADMIN Role = iota
	STUDENT
)

// String method to get the string representation of the Role
func (r Role) String() string {
	return [...]string{"ADMIN", "STUDENT"}[r]
}

// ParseRole function to convert a string to a Role type
func ParseRole(role string) (Role, error) {
	switch role {
	case "ADMIN":
		return ADMIN, nil
	case "STUDENT":
		return STUDENT, nil
	default:
		return -1, fmt.Errorf("invalid role: %s", role)
	}
}
