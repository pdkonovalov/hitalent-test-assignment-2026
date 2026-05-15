package command

import (
	"context"
	"time"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type CreateEmployee struct {
	DepartmentID domain.DepartmentID
	FullName     domain.EmployeeFullName
	Position     domain.EmployeePosition
	HiredAt      *time.Time
}

type CreateEmployeeHandler struct {
	employeeRepo interface {
		CreateEmployee(context.Context, *domain.Employee) (*domain.Employee, error)
	}
}

func NewCreateEmployeeHandler(
	employeeRepo interface {
		CreateEmployee(context.Context, *domain.Employee) (*domain.Employee, error)
	},
) CreateEmployeeHandler {
	return CreateEmployeeHandler{
		employeeRepo: employeeRepo,
	}
}

func (h CreateEmployeeHandler) Handle(ctx context.Context, cmd *CreateEmployee) (*domain.Employee, error) {
	newEmployee := domain.NewEmployee(cmd.DepartmentID, cmd.FullName, cmd.Position, cmd.HiredAt)

	return h.employeeRepo.CreateEmployee(ctx, newEmployee)
}
