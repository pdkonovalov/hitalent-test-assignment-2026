package domain

import "time"

type Employee struct {
	ID           EmployeeID
	DepartmentID DepartmentID
	FullName     EmployeeFullName
	Position     EmployeePosition
	HiredAt      *time.Time
	CreatedAt    time.Time
}

func NewEmployee(
	departmentID DepartmentID,
	fullName EmployeeFullName,
	position EmployeePosition,
	hiredAt *time.Time,
) *Employee {
	return &Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
		CreatedAt:    time.Now(),
	}
}
