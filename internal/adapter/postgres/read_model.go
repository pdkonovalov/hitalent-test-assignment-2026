package postgres

import (
	"context"
	"time"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app/query"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReadModel struct {
	db *gorm.DB
}

func NewReadModel(db *gorm.DB) *ReadModel {
	return &ReadModel{
		db: db,
	}
}

func (m *ReadModel) GetDepartment(ctx context.Context, params *query.GetDepartmentParams) (*query.GetDepartment, error) {
	type departmentTreeRow struct {
		ID        int64
		Name      string
		ParentID  *int64
		CreatedAt time.Time
		Depth     int
	}

	const departmentTreeQuery = `
		WITH RECURSIVE departments_tree AS (
			SELECT id, name, parent_id, created_at, 0 AS depth
			FROM departments
			WHERE id = ?

			UNION ALL

			SELECT d.id, d.name, d.parent_id, d.created_at, dt.depth + 1
			FROM departments d
			JOIN departments_tree dt ON d.parent_id = dt.id
			WHERE dt.depth < ?
		)
		SELECT * FROM departments_tree
	`

	departmentTreeRows := make([]*departmentTreeRow, 0)

	res := m.db.WithContext(ctx).
		Raw(departmentTreeQuery, params.ID, params.Depth).
		Scan(&departmentTreeRows)
	if err := res.Error; err != nil {
		return nil, err
	}
	if len(departmentTreeRows) == 0 {
		return nil, domain.ErrDepartmentNotFound
	}

	getDepartment := query.GetDepartment{
		Employees: make([]*query.Employee, 0),
		Children:  make([]*query.Department, 0),
	}

	departmentIDs := make([]int64, 0, len(departmentTreeRows))

	for _, row := range departmentTreeRows {
		department := query.Department{
			ID:        row.ID,
			Name:      row.Name,
			ParentID:  row.ParentID,
			CreatedAt: row.CreatedAt,
		}

		if row.Depth == 0 {
			getDepartment.Department = &department
		} else {
			getDepartment.Children = append(getDepartment.Children, &department)
		}

		departmentIDs = append(departmentIDs, department.ID)
	}

	if params.IncludeEmployees {
		res := m.db.WithContext(ctx).
			Where("department_id IN ?", departmentIDs).
			Find(&getDepartment.Employees).
			Order(clause.OrderBy{Columns: []clause.OrderByColumn{
				{Column: clause.Column{Name: "created_at"}},
				{Column: clause.Column{Name: "full_name"}},
			}})
		if err := res.Error; err != nil {
			return nil, err
		}
	}

	return &getDepartment, nil
}
