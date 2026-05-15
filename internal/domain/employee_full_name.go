package domain

type EmployeeFullName string

func NewEmployeeFullName(s string) (EmployeeFullName, error) {
	switch {
	case len(s) == 0:
		return "", ErrEmployeeFullNameEmpty
	case len(s) > 200:
		return "", ErrEmployeeFullNameTooLong
	default:
		return EmployeeFullName(s), nil
	}
}

func (n EmployeeFullName) String() string {
	return string(n)
}
