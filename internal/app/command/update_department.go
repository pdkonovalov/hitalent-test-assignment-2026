package command

import (
	"context"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type UpdateDepartment struct {
	ID          domain.DepartmentID
	NewName     *domain.DepartmentName
	NewParentID *domain.DepartmentID
}

type UpdateDepartmentHandler struct {
	txManager interface {
		DoSerializable(context.Context, func(context.Context) error) error
	}
	departmentRepo interface {
		GetDepartment(context.Context, domain.DepartmentID) (*domain.Department, error)
		CheckDepartmentsRelashionship(context.Context, domain.DepartmentID, domain.DepartmentID) (bool, error)
		UpdateDepartmentName(context.Context, domain.DepartmentID, domain.DepartmentName) (*domain.Department, error)
		UpdateDepartmentParentID(context.Context, domain.DepartmentID, domain.DepartmentID) (*domain.Department, error)
		UpdateDepartmentNameAndParentID(context.Context, domain.DepartmentID, domain.DepartmentName, domain.DepartmentID) (*domain.Department, error)
	}
}

func NewUpdateDepartmentHandler(
	txManager interface {
		DoSerializable(context.Context, func(context.Context) error) error
	},
	departmentRepo interface {
		GetDepartment(context.Context, domain.DepartmentID) (*domain.Department, error)
		CheckDepartmentsRelashionship(context.Context, domain.DepartmentID, domain.DepartmentID) (bool, error)
		UpdateDepartmentName(context.Context, domain.DepartmentID, domain.DepartmentName) (*domain.Department, error)
		UpdateDepartmentParentID(context.Context, domain.DepartmentID, domain.DepartmentID) (*domain.Department, error)
		UpdateDepartmentNameAndParentID(context.Context, domain.DepartmentID, domain.DepartmentName, domain.DepartmentID) (*domain.Department, error)
	},
) UpdateDepartmentHandler {
	return UpdateDepartmentHandler{
		txManager:      txManager,
		departmentRepo: departmentRepo,
	}
}

func (h UpdateDepartmentHandler) Handle(ctx context.Context, cmd *UpdateDepartment) (*domain.Department, error) {
	if cmd.NewName == nil && cmd.NewParentID == nil {
		return h.departmentRepo.GetDepartment(ctx, cmd.ID)
	}

	if cmd.NewName != nil && cmd.NewParentID == nil {
		return h.departmentRepo.UpdateDepartmentName(ctx, cmd.ID, *cmd.NewName)
	}

	if cmd.ID == *cmd.NewParentID {
		return nil, domain.ErrDepartmentParentConflict
	}

	var department *domain.Department

	if err := h.txManager.DoSerializable(ctx, func(ctx context.Context) error {
		newParentIsChild, err := h.departmentRepo.CheckDepartmentsRelashionship(ctx, cmd.ID, *cmd.NewParentID)
		if err != nil {
			return err
		}
		if newParentIsChild {
			return domain.ErrDepartmentParentConflict
		}

		if cmd.NewName == nil {
			department, err = h.departmentRepo.UpdateDepartmentParentID(ctx, cmd.ID, *cmd.NewParentID)
		} else {
			department, err = h.departmentRepo.UpdateDepartmentNameAndParentID(ctx, cmd.ID, *cmd.NewName, *cmd.NewParentID)
		}

		return err
	}); err != nil {
		return nil, err
	}

	return department, nil
}
