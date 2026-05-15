package postgres

import (
	"context"
	"errors"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db txContextGetter
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: newTxContextGetter(db),
	}
}

func (r *Repository) CreateDepartment(ctx context.Context, newDepartment *domain.Department) (*domain.Department, error) {
	res := r.db(ctx).Create(newDepartment)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, domain.ErrDepartmentParentNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrDepartmentNameConflict
		}

		return nil, err
	}

	return newDepartment, nil
}

func (r *Repository) GetDepartment(ctx context.Context, id domain.DepartmentID) (*domain.Department, error) {
	department := &domain.Department{}

	res := r.db(ctx).First(department, id)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDepartmentNotFound
		}

		return nil, err
	}

	return department, nil
}

func (r *Repository) UpdateDepartmentName(ctx context.Context, id domain.DepartmentID, newName domain.DepartmentName) (*domain.Department, error) {
	department := &domain.Department{}

	res := r.db(ctx).Model(department).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("name", newName)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrDepartmentNameConflict
		}

		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, domain.ErrDepartmentNotFound
	}

	return department, nil
}

func (r *Repository) UpdateDepartmentParentID(ctx context.Context, id domain.DepartmentID, newParentID domain.DepartmentID) (*domain.Department, error) {
	department := &domain.Department{}

	res := r.db(ctx).Model(department).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Update("parent_id", newParentID)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, domain.ErrDepartmentParentNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrDepartmentNameConflict
		}

		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, domain.ErrDepartmentNotFound
	}

	return department, nil
}

func (r *Repository) UpdateDepartmentNameAndParentID(ctx context.Context, id domain.DepartmentID, newName domain.DepartmentName, newParentID domain.DepartmentID) (*domain.Department, error) {
	department := &domain.Department{}

	res := r.db(ctx).Model(department).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"name": newName, "parent_id": newParentID})
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, domain.ErrDepartmentParentNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, domain.ErrDepartmentNameConflict
		}

		return nil, err
	}
	if res.RowsAffected == 0 {
		return nil, domain.ErrDepartmentNotFound
	}

	return department, nil
}

func (r *Repository) DeleteDepartment(ctx context.Context, id domain.DepartmentID) error {
	res := r.db(ctx).Delete(&domain.Department{}, id)
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return domain.ErrDepartmentNotFound
	}

	return nil
}

func (r *Repository) CheckDepartmentsRelashionship(
	ctx context.Context,
	parentID domain.DepartmentID,
	childID domain.DepartmentID,
) (bool, error) {
	query := `
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_id
			FROM departments
			WHERE id = ?

			UNION ALL

			SELECT d.id, d.parent_id
			FROM departments d
			JOIN ancestors a ON d.id = a.parent_id
		)
		SELECT EXISTS (
			SELECT 1 FROM ancestors WHERE id = ?
		)
	`

	var exists bool
	res := r.db(ctx).Raw(query, childID, parentID).Scan(&exists)
	if err := res.Error; err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) DeleteDepartmentsParentID(ctx context.Context, parentID domain.DepartmentID) error {
	res := r.db(ctx).Model(&domain.Department{}).
		Where("parent_id = ?", parentID).
		Updates(map[string]interface{}{"parent_id": nil})
	if err := res.Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateEmployee(ctx context.Context, newEmployee *domain.Employee) (*domain.Employee, error) {
	res := r.db(ctx).Create(newEmployee)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return nil, domain.ErrEmployeeDepartmentNotFound
		}

		return nil, err
	}

	return newEmployee, nil
}

func (r *Repository) UpdateEmployeesDepartmentID(ctx context.Context, oldDepartmentID, newDepartmentID domain.DepartmentID) error {
	res := r.db(ctx).Model(&domain.Employee{}).
		Where("department_id = ?", oldDepartmentID).
		Update("department_id", newDepartmentID)
	if err := res.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return domain.ErrEmployeeDepartmentNotFound
		}

		return err
	}

	return nil
}
