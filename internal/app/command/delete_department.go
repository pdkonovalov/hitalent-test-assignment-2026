package command

import (
	"context"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type DeleteDepartment struct {
	ID                    domain.DepartmentID
	ReasignToDepartmentID *domain.DepartmentID
}

type DeleteDepartmentHandler struct {
	txManager interface {
		Do(context.Context, func(context.Context) error) error
	}
	departmentRepo interface {
		DeleteDepartment(context.Context, domain.DepartmentID) error
		DeleteDepartmentsParentID(context.Context, domain.DepartmentID) error
	}
	employeeRepo interface {
		UpdateEmployeesDepartmentID(context.Context, domain.DepartmentID, domain.DepartmentID) error
	}
}

func NewDeleteDepartmentHandler(
	txManager interface {
		Do(context.Context, func(context.Context) error) error
	},
	departmentRepo interface {
		DeleteDepartment(context.Context, domain.DepartmentID) error
		DeleteDepartmentsParentID(context.Context, domain.DepartmentID) error
	},
	employeeRepo interface {
		UpdateEmployeesDepartmentID(context.Context, domain.DepartmentID, domain.DepartmentID) error
	},
) DeleteDepartmentHandler {
	return DeleteDepartmentHandler{
		txManager:      txManager,
		departmentRepo: departmentRepo,
		employeeRepo:   employeeRepo,
	}
}

func (h DeleteDepartmentHandler) Handle(ctx context.Context, cmd *DeleteDepartment) error {
	if cmd.ReasignToDepartmentID == nil {
		return h.departmentRepo.DeleteDepartment(ctx, cmd.ID)
	}

	return h.txManager.Do(ctx, func(ctx context.Context) error {
		if err := h.departmentRepo.DeleteDepartmentsParentID(ctx, cmd.ID); err != nil {
			return err
		}

		if err := h.employeeRepo.UpdateEmployeesDepartmentID(ctx, cmd.ID, *cmd.ReasignToDepartmentID); err != nil {
			return err
		}

		return h.departmentRepo.DeleteDepartment(ctx, cmd.ID)
	})
}
