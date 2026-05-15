package domain

import "strings"

type DepartmentName string

func NewDepartmentName(s string) (DepartmentName, error) {
	s = strings.Trim(s, " ")

	switch {
	case len(s) == 0:
		return "", ErrDepartmentNameEmpty
	case len(s) > 200:
		return "", ErrDepartmentNameTooLong
	default:
		return DepartmentName(s), nil
	}
}

func (n DepartmentName) String() string {
	return string(n)
}
