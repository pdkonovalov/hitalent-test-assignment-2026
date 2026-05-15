package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app/command"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/app/query"
	"github.com/pdkonovalov/hitalent-test-assignment-2026/internal/domain"
)

type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *int64 `json:"parent_id,omitempty"`
}

type CreateDepartmentResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateDepartment godoc
//
//	@Summary		Create a department
//	@Description	Creates a new department, optionally nested under a parent
//	@Tags			departments
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateDepartmentRequest	true	"Department payload"
//	@Success		201		{object}	CreateDepartmentResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body"
//	@Failure		409		{object}	ErrorResponse	"Department name conflict"
//	@Failure		404		{object}	ErrorResponse	"Department parent not found"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/departments/ [post]
func CreateDepartment(
	commandHandler command.CreateDepartmentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := readRequest[CreateDepartmentRequest](r)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		cmd := command.CreateDepartment{}

		name, err := domain.NewDepartmentName(request.Name)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}
		cmd.Name = name

		if request.ParentID != nil {
			cmd.ParentID = new(domain.NewDepartmentID(*request.ParentID))
		}

		department, err := commandHandler.Handle(r.Context(), &cmd)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		response := CreateDepartmentResponse{
			ID:        department.ID.Int64(),
			Name:      department.Name.String(),
			CreatedAt: department.CreatedAt,
		}

		if department.ParentID != nil {
			response.ParentID = new(department.ParentID.Int64())
		}

		writeResponse(w, r, response, http.StatusCreated)
	}
}

type CreateEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at,omitempty"`
}

type CreateEmployeeResponse struct {
	ID           int64      `json:"id"`
	DepartmentID int64      `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// CreateEmployee godoc
//
//	@Summary		Create an employee in a department
//	@Description	Creates a new employee and assigns them to the specified department
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Department ID"
//	@Param			body	body		CreateEmployeeRequest	true	"Employee payload"
//	@Success		201		{object}	CreateEmployeeResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body or department ID"
//	@Failure		404		{object}	ErrorResponse	"Department not found"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/departments/{id}/employees/ [post]
func CreateEmployee(
	commandHandler command.CreateEmployeeHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		departmentIDStr := r.PathValue("id")
		departmentIDParsed, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		departmentID := domain.NewDepartmentID(departmentIDParsed)

		request, err := readRequest[CreateEmployeeRequest](r)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		fullName, err := domain.NewEmployeeFullName(request.FullName)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		position, err := domain.NewEmployeePosition(request.Position)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		cmd := command.CreateEmployee{
			DepartmentID: departmentID,
			FullName:     fullName,
			Position:     position,
			HiredAt:      request.HiredAt,
		}

		employee, err := commandHandler.Handle(r.Context(), &cmd)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		response := CreateEmployeeResponse{
			ID:           employee.ID.Int64(),
			DepartmentID: employee.DepartmentID.Int64(),
			FullName:     employee.FullName.String(),
			Position:     employee.Position.String(),
			HiredAt:      employee.HiredAt,
			CreatedAt:    employee.CreatedAt,
		}

		writeResponse(w, r, response, http.StatusCreated)
	}
}

// GetDepartment godoc
//
//	@Summary		Get department details with employees and subtree
//	@Description	Returns a department with its employees and nested sub-departments up to the specified depth
//	@Tags			departments
//	@Produce		json
//	@Param			id					path		int		true	"Department ID"
//	@Param			depth				query		int		false	"Depth of nested departments (default: 1, max: 5)"	minimum(1)	maximum(5)	default(1)
//	@Param			include_employees	query		bool	false	"Include employees in response (default: true)"		default(true)
//	@Success		200					{object}	query.GetDepartment
//	@Failure		400					{object}	ErrorResponse	"Invalid parameters"
//	@Failure		404					{object}	ErrorResponse	"Department not found"
//	@Failure		500					{object}	ErrorResponse	"Internal server error"
//	@Router			/departments/{id} [get]
func GetDepartment(
	queryHandler query.GetDepartmentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		departmentIDStr := r.PathValue("id")
		departmentIDParsed, err := strconv.ParseInt(departmentIDStr, 10, 64)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		departmentID := domain.NewDepartmentID(departmentIDParsed)

		var depth *int
		depthStr := r.URL.Query().Get("depth")
		if len(depthStr) != 0 {
			depthParsed, err := strconv.Atoi(depthStr)
			if err != nil {
				writeBadRequestError(w, r, err)
				return
			}

			depth = &depthParsed
		}

		var includeEmployees *bool
		includeEmployeesStr := r.URL.Query().Get("include_employees")
		if len(includeEmployeesStr) != 0 {
			includeEmployeesParsed, err := strconv.ParseBool(includeEmployeesStr)
			if err != nil {
				writeBadRequestError(w, r, err)
				return
			}

			includeEmployees = &includeEmployeesParsed
		}

		params, err := query.NewGetDepartmentParams(departmentID, depth, includeEmployees)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		response, err := queryHandler.Handle(r.Context(), params)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		writeResponse(w, r, response, http.StatusOK)
	}
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name,omitempty"`
	ParentID *int64  `json:"parent_id,omitempty"`
}

type UpdateDepartmentResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	ParentID  *int64    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateDepartment godoc
//
//	@Summary		Update a department
//	@Description	Updates department name and/or moves it under a new parent. All body fields are optional.
//	@Tags			departments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Department ID"
//	@Param			body	body		UpdateDepartmentRequest	true	"Fields to update"
//	@Success		200		{object}	UpdateDepartmentResponse
//	@Failure		400		{object}	ErrorResponse	"Invalid request body or department ID"
//	@Failure		404		{object}	ErrorResponse	"Department not found"
//	@Failure		409		{object}	ErrorResponse	"Circular parent reference detected"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/departments/{id} [patch]
func UpdateDepartment(
	commandHandler command.UpdateDepartmentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		idParsed, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		request, err := readRequest[UpdateDepartmentRequest](r)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		cmd := command.UpdateDepartment{
			ID: domain.NewDepartmentID(idParsed),
		}

		if request.Name != nil {
			name, err := domain.NewDepartmentName(*request.Name)
			if err != nil {
				writeDomainOrInternalError(w, r, err)
				return
			}

			cmd.NewName = &name
		}

		if request.ParentID != nil {
			cmd.NewParentID = new(domain.NewDepartmentID(*request.ParentID))
		}

		department, err := commandHandler.Handle(r.Context(), &cmd)
		if err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		response := UpdateDepartmentResponse{
			ID:        department.ID.Int64(),
			Name:      department.Name.String(),
			CreatedAt: department.CreatedAt,
		}

		if department.ParentID != nil {
			response.ParentID = new(department.ParentID.Int64())
		}

		writeResponse(w, r, response, http.StatusOK)
	}
}

// DeleteDepartment godoc
//
//	@Summary		Delete a department
//	@Description	Deletes a department using one of two modes:
//	@Description	- `cascade` — deletes the department along with all its employees and sub-departments recursively.
//	@Description	- `reassign` — deletes the department and moves its employees to the department specified by `reassign_to_department_id`.
//	@Description	Note: `reassign_to_department_id` is required when mode is `reassign`.
//	@Tags			departments
//	@Produce		json
//	@Param			id							path	int		true	"Department ID"
//	@Param			mode						query	string	true	"Deletion mode"	Enums(cascade, reassign)
//	@Param			reassign_to_department_id	query	int		false	"Target department ID for employee reassignment (required when mode=reassign)"
//	@Success		204							"No Content"
//	@Failure		400							{object}	ErrorResponse	"Invalid mode or missing reassign_to_department_id"
//	@Failure		404							{object}	ErrorResponse	"Department not found"
//	@Failure		500							{object}	ErrorResponse	"Internal server error"
//	@Router			/departments/{id} [delete]
func DeleteDepartment(
	commandHandler command.DeleteDepartmentHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		idParsed, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			writeBadRequestError(w, r, err)
			return
		}

		cmd := command.DeleteDepartment{
			ID: domain.NewDepartmentID(idParsed),
		}

		mode := r.URL.Query().Get("mode")

		if mode != "cascade" && mode != "reassign" {
			writeBadRequestError(w, r, nil)
			return
		} else if mode == "reassign" {
			reassignToStr := r.URL.Query().Get("reassign_to_department_id")
			reassignToParsed, err := strconv.ParseInt(reassignToStr, 10, 64)
			if err != nil {
				writeBadRequestError(w, r, err)
				return
			}

			cmd.ReasignToDepartmentID = new(domain.NewDepartmentID(reassignToParsed))
		}

		if err := commandHandler.Handle(r.Context(), &cmd); err != nil {
			writeDomainOrInternalError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
