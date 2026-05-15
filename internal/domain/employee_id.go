package domain

type EmployeeID int64

func NewEmployeeID(n int64) EmployeeID {
	return EmployeeID(n)
}

func (id EmployeeID) Int64() int64 {
	return int64(id)
}
