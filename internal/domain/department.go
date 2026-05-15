package domain

import (
	"time"
)

type Department struct {
	ID        DepartmentID
	Name      DepartmentName
	ParentID  *DepartmentID
	CreatedAt time.Time
}

func NewDepartment(name DepartmentName, parrentID *DepartmentID) *Department {
	return &Department{
		Name:      name,
		ParentID:  parrentID,
		CreatedAt: time.Now(),
	}
}
