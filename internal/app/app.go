package app

import (
	"context"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app/command"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app/query"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type App struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateDepartment command.CreateDepartmentHandler
	UpdateDepartment command.UpdateDepartmentHandler
	DeleteDepartment command.DeleteDepartmentHandler

	CreateEmployee command.CreateEmployeeHandler
}

type Queries struct {
	GetDepartment query.GetDepartmentHandler
}

func New(
	txManager interface {
		Do(context.Context, func(context.Context) error) error
		DoSerializable(context.Context, func(context.Context) error) error
	},
	departmentRepo interface {
		CreateDepartment(context.Context, *domain.Department) (*domain.Department, error)
		GetDepartment(context.Context, domain.DepartmentID) (*domain.Department, error)
		CheckDepartmentsRelashionship(context.Context, domain.DepartmentID, domain.DepartmentID) (bool, error)
		UpdateDepartmentName(context.Context, domain.DepartmentID, domain.DepartmentName) (*domain.Department, error)
		UpdateDepartmentParentID(context.Context, domain.DepartmentID, domain.DepartmentID) (*domain.Department, error)
		UpdateDepartmentNameAndParentID(context.Context, domain.DepartmentID, domain.DepartmentName, domain.DepartmentID) (*domain.Department, error)
		DeleteDepartment(context.Context, domain.DepartmentID) error
		DeleteDepartmentsParentID(context.Context, domain.DepartmentID) error
	},
	employeeRepo interface {
		CreateEmployee(context.Context, *domain.Employee) (*domain.Employee, error)
		UpdateEmployeesDepartmentID(context.Context, domain.DepartmentID, domain.DepartmentID) error
	},
	readModel interface {
		GetDepartment(context.Context, *query.GetDepartmentParams) (*query.GetDepartment, error)
	},
) *App {
	commands := Commands{
		CreateDepartment: command.NewCreateDepartmentHandler(departmentRepo),
		UpdateDepartment: command.NewUpdateDepartmentHandler(txManager, departmentRepo),
		DeleteDepartment: command.NewDeleteDepartmentHandler(txManager, departmentRepo, employeeRepo),

		CreateEmployee: command.NewCreateEmployeeHandler(employeeRepo),
	}

	queries := Queries{
		GetDepartment: query.NewGetDepartmentHandler(readModel),
	}

	return &App{
		Commands: commands,
		Queries:  queries,
	}
}
