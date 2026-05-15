package domain

type DepartmentID int64

func NewDepartmentID(n int64) DepartmentID {
	return DepartmentID(n)
}

func (id DepartmentID) Int64() int64 {
	return int64(id)
}
