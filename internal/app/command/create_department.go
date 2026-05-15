package command

import (
	"context"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type CreateDepartment struct {
	Name     domain.DepartmentName
	ParentID *domain.DepartmentID
}

type CreateDepartmentHandler struct {
	departmentRepo interface {
		CreateDepartment(context.Context, *domain.Department) (*domain.Department, error)
	}
}

func NewCreateDepartmentHandler(
	departmentRepo interface {
		CreateDepartment(context.Context, *domain.Department) (*domain.Department, error)
	},
) CreateDepartmentHandler {
	return CreateDepartmentHandler{
		departmentRepo: departmentRepo,
	}
}

func (h CreateDepartmentHandler) Handle(ctx context.Context, cmd *CreateDepartment) (*domain.Department, error) {
	newDepartment := domain.NewDepartment(cmd.Name, cmd.ParentID)

	return h.departmentRepo.CreateDepartment(ctx, newDepartment)
}
