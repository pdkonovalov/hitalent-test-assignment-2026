package e2e

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (s *E2ESuite) mustCreateDepartment(name string, parentID ...int64) CreateDepartmentResponse {
	req := CreateDepartmentRequest{Name: name}
	if len(parentID) > 0 {
		req.ParentID = &parentID[0]
	}
	rec := s.do(http.MethodPost, "/departments/", req)
	s.Require().Equal(http.StatusCreated, rec.Code,
		"mustCreateDepartment(%q) unexpected status", name)
	return decode[CreateDepartmentResponse](s.T(), rec)
}

func (s *E2ESuite) mustCreateEmployee(deptID int64, fullName, position string) CreateEmployeeResponse {
	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", deptID),
		CreateEmployeeRequest{FullName: fullName, Position: position},
	)
	s.Require().Equal(http.StatusCreated, rec.Code,
		"mustCreateEmployee(%q) unexpected status", fullName)
	return decode[CreateEmployeeResponse](s.T(), rec)
}

func childIDs(departments []*Department) []int64 {
	ids := make([]int64, len(departments))
	for i, d := range departments {
		ids[i] = d.ID
	}
	return ids
}

// ============================================================
// POST /departments/
// ============================================================

func (s *E2ESuite) TestCreateDepartment_RootSuccess() {
	rec := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{Name: "Engineering"})
	s.Require().Equal(http.StatusCreated, rec.Code)

	resp := decode[CreateDepartmentResponse](s.T(), rec)
	s.Equal("Engineering", resp.Name)
	s.Nil(resp.ParentID)
	s.NotZero(resp.ID)
	s.False(resp.CreatedAt.IsZero())
}

func (s *E2ESuite) TestCreateDepartment_WithParent() {
	root := s.mustCreateDepartment("RootDept")
	child := s.mustCreateDepartment("BackendDept", root.ID)

	s.Require().NotNil(child.ParentID)
	s.Equal(root.ID, *child.ParentID)
	s.NotZero(child.ID)
}

func (s *E2ESuite) TestCreateDepartment_EmptyName_Returns4xx() {
	rec := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{Name: ""})
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestCreateDepartment_NameTooLong_Returns4xx() {
	rec := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{
		Name: strings.Repeat("x", 201),
	})
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestCreateDepartment_WhitespaceName_IsTrimmed() {
	rec := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{Name: "  Trimmed  "})
	s.Require().Equal(http.StatusCreated, rec.Code)

	resp := decode[CreateDepartmentResponse](s.T(), rec)
	s.Equal("Trimmed", resp.Name)
}

func (s *E2ESuite) TestCreateDepartment_DuplicateNameSameParent_Returns409() {
	parent := s.mustCreateDepartment("DupParent")

	first := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{
		Name:     "DupChild",
		ParentID: &parent.ID,
	})
	s.Require().Equal(http.StatusCreated, first.Code)

	second := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{
		Name:     "DupChild",
		ParentID: &parent.ID,
	})
	s.Equal(http.StatusConflict, second.Code)
}

func (s *E2ESuite) TestCreateDepartment_SameNameDifferentParents_IsAllowed() {
	p1 := s.mustCreateDepartment("ParentAlpha")
	p2 := s.mustCreateDepartment("ParentBeta")

	rec1 := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{
		Name:     "SharedName",
		ParentID: &p1.ID,
	})
	s.Require().Equal(http.StatusCreated, rec1.Code)

	rec2 := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{
		Name:     "SharedName",
		ParentID: &p2.ID,
	})
	s.Require().Equal(http.StatusCreated, rec2.Code)
}

func (s *E2ESuite) TestCreateDepartment_NonExistentParent_Returns404() {
	nonExistentID := int64(999_999)
	rec := s.do(http.MethodPost, "/departments/", CreateDepartmentRequest{
		Name:     "Orphan",
		ParentID: &nonExistentID,
	})
	s.Equal(http.StatusNotFound, rec.Code)
}

// ============================================================
// POST /departments/{id}/employees/
// ============================================================

func (s *E2ESuite) TestCreateEmployee_Success() {
	dept := s.mustCreateDepartment("HRDept")
	hiredAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)

	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", dept.ID),
		CreateEmployeeRequest{
			FullName: "Alice Smith",
			Position: "Senior Engineer",
			HiredAt:  &hiredAt,
		},
	)
	s.Require().Equal(http.StatusCreated, rec.Code)

	emp := decode[CreateEmployeeResponse](s.T(), rec)
	s.Equal("Alice Smith", emp.FullName)
	s.Equal("Senior Engineer", emp.Position)
	s.Equal(dept.ID, emp.DepartmentID)
	s.NotZero(emp.ID)
	s.False(emp.CreatedAt.IsZero())
	s.Require().NotNil(emp.HiredAt)
	s.Equal(hiredAt.UTC(), emp.HiredAt.UTC())
}

func (s *E2ESuite) TestCreateEmployee_WithoutHiredAt() {
	dept := s.mustCreateDepartment("FinanceDept")

	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", dept.ID),
		CreateEmployeeRequest{FullName: "Bob Jones", Position: "Analyst"},
	)
	s.Require().Equal(http.StatusCreated, rec.Code)

	emp := decode[CreateEmployeeResponse](s.T(), rec)
	s.Nil(emp.HiredAt)
}

func (s *E2ESuite) TestCreateEmployee_NonExistentDepartment_Returns404() {
	rec := s.do(
		http.MethodPost,
		"/departments/999999/employees/",
		CreateEmployeeRequest{FullName: "Ghost", Position: "Unknown"},
	)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestCreateEmployee_EmptyFullName_Returns4xx() {
	dept := s.mustCreateDepartment("LegalDept")

	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", dept.ID),
		CreateEmployeeRequest{FullName: "", Position: "Lawyer"},
	)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestCreateEmployee_EmptyPosition_Returns4xx() {
	dept := s.mustCreateDepartment("MarketingDept")

	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", dept.ID),
		CreateEmployeeRequest{FullName: "Carol White", Position: ""},
	)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestCreateEmployee_FullNameTooLong_Returns4xx() {
	dept := s.mustCreateDepartment("SalesDept")

	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", dept.ID),
		CreateEmployeeRequest{FullName: strings.Repeat("x", 201), Position: "Salesperson"},
	)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestCreateEmployee_PositionTooLong_Returns4xx() {
	dept := s.mustCreateDepartment("DevDept")

	rec := s.do(
		http.MethodPost,
		fmt.Sprintf("/departments/%d/employees/", dept.ID),
		CreateEmployeeRequest{FullName: "Dan Brown", Position: strings.Repeat("x", 201)},
	)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422, got %d", rec.Code,
	)
}

// ============================================================
// GET /departments/{id}
// ============================================================

func (s *E2ESuite) TestGetDepartment_Success() {
	dept := s.mustCreateDepartment("GetTestDept")

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", dept.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	s.Require().NotNil(resp.Department)
	s.Equal(dept.ID, resp.Department.ID)
	s.Equal("GetTestDept", resp.Department.Name)
	s.False(resp.Department.CreatedAt.IsZero())
}

func (s *E2ESuite) TestGetDepartment_NotFound() {
	rec := s.do(http.MethodGet, "/departments/999999", nil)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestGetDepartment_IncludesEmployeesByDefault() {
	dept := s.mustCreateDepartment("WithEmpsDept")
	s.mustCreateEmployee(dept.ID, "Emp One", "Dev")
	s.mustCreateEmployee(dept.ID, "Emp Two", "QA")

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", dept.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	s.Len(resp.Employees, 2)
}

func (s *E2ESuite) TestGetDepartment_ExcludesEmployees_WhenFalse() {
	dept := s.mustCreateDepartment("HiddenEmpsDept")
	s.mustCreateEmployee(dept.ID, "Hidden Emp", "Dev")

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d?include_employees=false", dept.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	s.Empty(resp.Employees)
}

func (s *E2ESuite) TestGetDepartment_Children_DefaultDepthIs1() {
	root := s.mustCreateDepartment("DepthRoot")
	l1 := s.mustCreateDepartment("DepthL1", root.ID)
	l2 := s.mustCreateDepartment("DepthL2", l1.ID)

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", root.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	ids := childIDs(resp.Children)
	s.Contains(ids, l1.ID, "immediate child must appear at depth=1")
	s.NotContains(ids, l2.ID, "grandchild must not appear at depth=1")
}

func (s *E2ESuite) TestGetDepartment_Children_Depth2() {
	root := s.mustCreateDepartment("D2Root")
	l1 := s.mustCreateDepartment("D2L1", root.ID)
	l2 := s.mustCreateDepartment("D2L2", l1.ID)
	l3 := s.mustCreateDepartment("D2L3", l2.ID)

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d?depth=2", root.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	ids := childIDs(resp.Children)
	s.Contains(ids, l1.ID, "L1 must appear at depth=2")
	s.Contains(ids, l2.ID, "L2 must appear at depth=2")
	s.NotContains(ids, l3.ID, "L3 must not appear at depth=2")
}

func (s *E2ESuite) TestGetDepartment_Children_MaxDepth5() {
	root := s.mustCreateDepartment("MaxRoot")
	l1 := s.mustCreateDepartment("MaxL1", root.ID)
	l2 := s.mustCreateDepartment("MaxL2", l1.ID)
	l3 := s.mustCreateDepartment("MaxL3", l2.ID)
	l4 := s.mustCreateDepartment("MaxL4", l3.ID)
	l5 := s.mustCreateDepartment("MaxL5", l4.ID)

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d?depth=5", root.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	ids := childIDs(resp.Children)
	for _, expected := range []int64{l1.ID, l2.ID, l3.ID, l4.ID, l5.ID} {
		s.Contains(ids, expected, "department %d must appear at depth=5", expected)
	}
}

func (s *E2ESuite) TestGetDepartment_DepthExceedsMax_Returns4xx() {
	dept := s.mustCreateDepartment("DepthMaxDept")

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d?depth=6", dept.ID), nil)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422 for depth > 5, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestGetDepartment_Employees_SortedByCreatedAt() {
	dept := s.mustCreateDepartment("SortedDept")
	emp1 := s.mustCreateEmployee(dept.ID, "First Employee", "Dev")
	emp2 := s.mustCreateEmployee(dept.ID, "Second Employee", "QA")

	rec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", dept.ID), nil)
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[GetDepartmentResponse](s.T(), rec)
	s.Require().Len(resp.Employees, 2)
	s.Equal(emp1.ID, resp.Employees[0].ID)
	s.Equal(emp2.ID, resp.Employees[1].ID)
	s.False(resp.Employees[0].CreatedAt.After(resp.Employees[1].CreatedAt),
		"employees must be sorted by created_at ascending")
}

// ============================================================
// PATCH /departments/{id}
// ============================================================

func (s *E2ESuite) TestUpdateDepartment_RenameSuccess() {
	dept := s.mustCreateDepartment("OldName")

	newName := "NewName"
	rec := s.do(http.MethodPatch, fmt.Sprintf("/departments/%d", dept.ID), UpdateDepartmentRequest{
		Name: &newName,
	})
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[UpdateDepartmentResponse](s.T(), rec)
	s.Equal(dept.ID, resp.ID)
	s.Equal("NewName", resp.Name)
}

func (s *E2ESuite) TestUpdateDepartment_MoveToNewParent() {
	oldParent := s.mustCreateDepartment("OldParent")
	newParent := s.mustCreateDepartment("NewParent")
	child := s.mustCreateDepartment("MovableChild", oldParent.ID)

	rec := s.do(http.MethodPatch, fmt.Sprintf("/departments/%d", child.ID), UpdateDepartmentRequest{
		ParentID: &newParent.ID,
	})
	s.Require().Equal(http.StatusOK, rec.Code)

	resp := decode[UpdateDepartmentResponse](s.T(), rec)
	s.Require().NotNil(resp.ParentID)
	s.Equal(newParent.ID, *resp.ParentID)

	oldRec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", oldParent.ID), nil)
	s.Require().Equal(http.StatusOK, oldRec.Code)
	oldResp := decode[GetDepartmentResponse](s.T(), oldRec)
	s.NotContains(childIDs(oldResp.Children), child.ID, "child must be removed from old parent")

	newRec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", newParent.ID), nil)
	s.Require().Equal(http.StatusOK, newRec.Code)
	newResp := decode[GetDepartmentResponse](s.T(), newRec)
	s.Contains(childIDs(newResp.Children), child.ID, "child must appear under new parent")
}

func (s *E2ESuite) TestUpdateDepartment_NotFound() {
	name := "Nope"
	rec := s.do(http.MethodPatch, "/departments/999999", UpdateDepartmentRequest{Name: &name})
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestUpdateDepartment_SelfParent_Returns4xx() {
	dept := s.mustCreateDepartment("SelfRefDept")

	rec := s.do(http.MethodPatch, fmt.Sprintf("/departments/%d", dept.ID), UpdateDepartmentRequest{
		ParentID: &dept.ID,
	})
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusConflict,
		"expected 400 or 409, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestUpdateDepartment_CycleDetection_Returns4xx() {
	a := s.mustCreateDepartment("CycleA")
	b := s.mustCreateDepartment("CycleB", a.ID)
	c := s.mustCreateDepartment("CycleC", b.ID)

	rec := s.do(http.MethodPatch, fmt.Sprintf("/departments/%d", a.ID), UpdateDepartmentRequest{
		ParentID: &c.ID,
	})
	s.True(
		rec.Code == http.StatusConflict || rec.Code == http.StatusBadRequest,
		"expected 409 or 400 for cycle, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestUpdateDepartment_DuplicateNameInSameParent_Returns409() {
	parent := s.mustCreateDepartment("UpdateDupParent")
	_ = s.mustCreateDepartment("AlreadyExists", parent.ID)
	target := s.mustCreateDepartment("TargetToRename", parent.ID)

	existingName := "AlreadyExists"
	rec := s.do(http.MethodPatch, fmt.Sprintf("/departments/%d", target.ID), UpdateDepartmentRequest{
		Name: &existingName,
	})
	s.Equal(http.StatusConflict, rec.Code)
}

func (s *E2ESuite) TestUpdateDepartment_MoveCreatesNameConflict_Returns409() {
	srcParent := s.mustCreateDepartment("SrcParentForConflict")
	dstParent := s.mustCreateDepartment("DstParentForConflict")
	_ = s.mustCreateDepartment("ConflictName", dstParent.ID) // blocker already in dst
	mover := s.mustCreateDepartment("ConflictName", srcParent.ID)

	rec := s.do(http.MethodPatch, fmt.Sprintf("/departments/%d", mover.ID), UpdateDepartmentRequest{
		ParentID: &dstParent.ID,
	})
	s.Equal(http.StatusConflict, rec.Code)
}

// ============================================================
// DELETE /departments/{id}
// ============================================================

func (s *E2ESuite) TestDeleteDepartment_Cascade_RemovesDepartmentItself() {
	dept := s.mustCreateDepartment("CascadeMe")

	rec := s.do(http.MethodDelete, fmt.Sprintf("/departments/%d?mode=cascade", dept.ID), nil)
	s.Require().Equal(http.StatusNoContent, rec.Code)

	getRec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d", dept.ID), nil)
	s.Equal(http.StatusNotFound, getRec.Code)
}

func (s *E2ESuite) TestDeleteDepartment_Cascade_RemovesChildrenRecursively() {
	root := s.mustCreateDepartment("CascadeRoot")
	child := s.mustCreateDepartment("CascadeChild", root.ID)
	grandchild := s.mustCreateDepartment("CascadeGrandchild", child.ID)

	rec := s.do(http.MethodDelete, fmt.Sprintf("/departments/%d?mode=cascade", root.ID), nil)
	s.Require().Equal(http.StatusNoContent, rec.Code)

	s.Equal(http.StatusNotFound,
		s.do(http.MethodGet, fmt.Sprintf("/departments/%d", child.ID), nil).Code,
		"child department should be cascade-deleted")
	s.Equal(http.StatusNotFound,
		s.do(http.MethodGet, fmt.Sprintf("/departments/%d", grandchild.ID), nil).Code,
		"grandchild department should be cascade-deleted")
}

func (s *E2ESuite) TestDeleteDepartment_Cascade_RemovesDepartmentWithEmployees() {
	dept := s.mustCreateDepartment("CascadeWithEmps")
	s.mustCreateEmployee(dept.ID, "Soon Gone", "Dev")

	rec := s.do(http.MethodDelete, fmt.Sprintf("/departments/%d?mode=cascade", dept.ID), nil)
	s.Require().Equal(http.StatusNoContent, rec.Code)

	s.Equal(http.StatusNotFound,
		s.do(http.MethodGet, fmt.Sprintf("/departments/%d", dept.ID), nil).Code)
}

func (s *E2ESuite) TestDeleteDepartment_Reassign_MovesEmployeesToTarget() {
	src := s.mustCreateDepartment("ReassignSrc")
	dst := s.mustCreateDepartment("ReassignDst")
	s.mustCreateEmployee(src.ID, "Migrant Employee", "Dev")

	path := fmt.Sprintf("/departments/%d?mode=reassign&reassign_to_department_id=%d", src.ID, dst.ID)
	rec := s.do(http.MethodDelete, path, nil)
	s.Require().Equal(http.StatusNoContent, rec.Code)

	s.Equal(http.StatusNotFound,
		s.do(http.MethodGet, fmt.Sprintf("/departments/%d", src.ID), nil).Code,
		"source department should be deleted after reassign")

	getRec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d?include_employees=true", dst.ID), nil)
	s.Require().Equal(http.StatusOK, getRec.Code)
	dstResp := decode[GetDepartmentResponse](s.T(), getRec)
	s.Require().Len(dstResp.Employees, 1)
	s.Equal("Migrant Employee", dstResp.Employees[0].FullName)
}

func (s *E2ESuite) TestDeleteDepartment_Reassign_MovesAllEmployees() {
	src := s.mustCreateDepartment("MultiReassignSrc")
	dst := s.mustCreateDepartment("MultiReassignDst")
	s.mustCreateEmployee(src.ID, "Worker A", "Dev")
	s.mustCreateEmployee(src.ID, "Worker B", "QA")
	s.mustCreateEmployee(src.ID, "Worker C", "PM")

	path := fmt.Sprintf("/departments/%d?mode=reassign&reassign_to_department_id=%d", src.ID, dst.ID)
	rec := s.do(http.MethodDelete, path, nil)
	s.Require().Equal(http.StatusNoContent, rec.Code)

	getRec := s.do(http.MethodGet, fmt.Sprintf("/departments/%d?include_employees=true", dst.ID), nil)
	s.Require().Equal(http.StatusOK, getRec.Code)
	dstResp := decode[GetDepartmentResponse](s.T(), getRec)
	s.Len(dstResp.Employees, 3)
}

func (s *E2ESuite) TestDeleteDepartment_Reassign_MissingTargetID_Returns4xx() {
	dept := s.mustCreateDepartment("ReassignNoTarget")

	rec := s.do(http.MethodDelete, fmt.Sprintf("/departments/%d?mode=reassign", dept.ID), nil)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422 when reassign_to_department_id is absent, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestDeleteDepartment_Reassign_NonExistentTargetID_Returns4xx() {
	dept := s.mustCreateDepartment("ReassignBadTarget")
	s.mustCreateEmployee(dept.ID, "Reassigned Emp", "Dev")

	path := fmt.Sprintf("/departments/%d?mode=reassign&reassign_to_department_id=999999", dept.ID)
	rec := s.do(http.MethodDelete, path, nil)
	s.True(
		rec.Code == http.StatusNotFound || rec.Code == http.StatusBadRequest,
		"expected 404 or 400 for non-existent target, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestDeleteDepartment_NotFound() {
	rec := s.do(http.MethodDelete, "/departments/999999?mode=cascade", nil)
	s.Equal(http.StatusNotFound, rec.Code)
}

func (s *E2ESuite) TestDeleteDepartment_InvalidMode_Returns4xx() {
	dept := s.mustCreateDepartment("BadModeDept")

	rec := s.do(http.MethodDelete, fmt.Sprintf("/departments/%d?mode=invalid", dept.ID), nil)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422 for unknown mode, got %d", rec.Code,
	)
}

func (s *E2ESuite) TestDeleteDepartment_MissingMode_Returns4xx() {
	dept := s.mustCreateDepartment("NoModeDept")

	rec := s.do(http.MethodDelete, fmt.Sprintf("/departments/%d", dept.ID), nil)
	s.True(
		rec.Code == http.StatusBadRequest || rec.Code == http.StatusUnprocessableEntity,
		"expected 400 or 422 when mode is absent, got %d", rec.Code,
	)
}
