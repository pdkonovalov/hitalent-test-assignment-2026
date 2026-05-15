package domain

type EmployeePosition string

func NewEmployeePosition(s string) (EmployeePosition, error) {
	switch {
	case len(s) == 0:
		return "", ErrEmployeePositionEmpty
	case len(s) > 200:
		return "", ErrEmployeePositionTooLong
	default:
		return EmployeePosition(s), nil
	}
}

func (p EmployeePosition) String() string {
	return string(p)
}
