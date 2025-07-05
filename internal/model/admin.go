package model

import "slices"

var Admins []string

func IsAdmin(user string) bool {
	return slices.Contains(Admins, user)
}
