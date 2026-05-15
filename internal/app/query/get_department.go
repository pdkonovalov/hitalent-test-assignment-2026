package query

import (
	"context"
	"fmt"
	"time"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type GetDepartmentParams struct {
	ID               domain.DepartmentID
	Depth            int
	IncludeEmployees bool
}

type Department struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Employee struct {
	ID           int64      `json:"id"`
	DepartmentID int64      `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type GetDepartment struct {
	Department *Department   `json:"department"`
	Employees  []*Employee   `json:"employees"`
	Children   []*Department `json:"children"`
}

type GetDepartmentHandler struct {
	readModel interface {
		GetDepartment(context.Context, *GetDepartmentParams) (*GetDepartment, error)
	}
}

func NewGetDepartmentHandler(
	readModel interface {
		GetDepartment(context.Context, *GetDepartmentParams) (*GetDepartment, error)
	},
) GetDepartmentHandler {
	return GetDepartmentHandler{
		readModel: readModel,
	}
}

func (h GetDepartmentHandler) Handle(ctx context.Context, params *GetDepartmentParams) (*GetDepartment, error) {
	return h.readModel.GetDepartment(ctx, params)
}

func NewGetDepartmentParams(id domain.DepartmentID, depth *int, includeEmployees *bool) (*GetDepartmentParams, error) {
	params := GetDepartmentParams{
		ID:               id,
		Depth:            1,
		IncludeEmployees: true,
	}

	if depth != nil {
		if *depth < 1 || *depth > 5 {
			return nil, fmt.Errorf("departments depth invalid")
		}

		params.Depth = *depth
	}

	if includeEmployees != nil {
		params.IncludeEmployees = *includeEmployees
	}

	return &params, nil
}
