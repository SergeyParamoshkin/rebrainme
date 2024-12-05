package domain

import (
	"fmt"
	"strings"
)

type Role string

const (
	RolePrefix      = "r_"
	RoleAdmin  Role = RolePrefix + "admin"
	RoleUser   Role = RolePrefix + "user"
)

func (r Role) Validate() error {
	switch r {
	case RoleAdmin, RoleUser:
		return nil
	default:
		return fmt.Errorf("role is invalid. Valid values are: %v", strings.Join([]string{
			string(RoleAdmin),
			string(RoleUser),
		}, ", "))
	}
}
