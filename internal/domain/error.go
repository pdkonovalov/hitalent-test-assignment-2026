package domain

type Error struct {
	Code    string
	Message string
}

func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrDepartmentNotFound       = NewError("department_not_found", "department not found")
	ErrDepartmentParentNotFound = NewError("department_parent_not_found", "department parent not found")
	ErrDepartmentParentConflict = NewError("department_parent_conflict", "department parent conflict")
	ErrDepartmentNameConflict   = NewError("department_name_conflict", "department name conflict")
	ErrDepartmentNameEmpty      = NewError("department_name_empty", "department name empty")
	ErrDepartmentNameTooLong    = NewError("department_name_too_long", "department name too long")

	ErrEmployeeDepartmentNotFound = NewError("employee_department_not_found", "employee department not found")
	ErrEmployeeFullNameEmpty      = NewError("employee_full_name_empty", "employee full name empty")
	ErrEmployeeFullNameTooLong    = NewError("employee_full_name_too_long", "employee full name too long")
	ErrEmployeePositionEmpty      = NewError("employee_position_empty", "employee position empty")
	ErrEmployeePositionTooLong    = NewError("employee_position_too_long", "employee position too long")
)
