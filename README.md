# Description

Test assignment to implement rest api in golang. All requirements in [pdf file](https://github.com/pdkonovalov/hitalent-test-assignment-2026/blob/main/requirements.ru.pdf).

# Quickstart

Run with docker compose:

```bash
git clone https://github.com/pdkonovalov/hitalent-test-assignment-2026
cd hitalent-test-assignment-2026
cp docker.env.example .env
docker compose up -d
docker compose logs departments_api -f
```

Visit swagger web interface http://localhost:8080/swagger. Make some requests.

```bash
departments_api-1  | time=2026-05-15T14:32:57.604Z level=INFO msg="starting server" host=0.0.0.0 port=8080 base_path=""
departments_api-1  | time=2026-05-15T14:37:52.261Z level=INFO msg="request processed" request.method=POST request.path=/departments/ request.query="" request.status=Created request.duration_ms=5
departments_api-1  | time=2026-05-15T14:40:02.226Z level=INFO msg="request processed" request.method=GET request.path=/departments/1 request.query="depth=1&include_employees=true" request.status=OK request.duration_ms=11
departments_api-1  | time=2026-05-15T14:40:15.608Z level=INFO msg="request processed" request.method=GET request.path=/departments/2 request.query="depth=1&include_employees=true" request.status="Not Found" request.duration_ms=2 request.errors="[department not found]"
```
All requests log into console. For log only internal server set `LOG_LEVEL=error` in `.env` file.

# Tests

Run all end to end tests:

```bash
make test
```

This command run more than 40 tests of api endpoints via testcontainers.

```bash
--- PASS: TestE2ESuite (7.62s)
    --- PASS: TestE2ESuite/TestCreateDepartment_DuplicateNameSameParent_Returns409 (0.01s)
    --- PASS: TestE2ESuite/TestCreateDepartment_EmptyName_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestCreateDepartment_NameTooLong_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestCreateDepartment_NonExistentParent_Returns404 (0.00s)
    --- PASS: TestE2ESuite/TestCreateDepartment_RootSuccess (0.00s)
    --- PASS: TestE2ESuite/TestCreateDepartment_SameNameDifferentParents_IsAllowed (0.01s)
    --- PASS: TestE2ESuite/TestCreateDepartment_WhitespaceName_IsTrimmed (0.00s)
    --- PASS: TestE2ESuite/TestCreateDepartment_WithParent (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_EmptyFullName_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_EmptyPosition_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_FullNameTooLong_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_NonExistentDepartment_Returns404 (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_PositionTooLong_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_Success (0.00s)
    --- PASS: TestE2ESuite/TestCreateEmployee_WithoutHiredAt (0.01s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Cascade_RemovesChildrenRecursively (0.01s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Cascade_RemovesDepartmentItself (0.00s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Cascade_RemovesDepartmentWithEmployees (0.01s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_InvalidMode_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_MissingMode_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_NotFound (0.00s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Reassign_MissingTargetID_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Reassign_MovesAllEmployees (0.02s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Reassign_MovesEmployeesToTarget (0.01s)
    --- PASS: TestE2ESuite/TestDeleteDepartment_Reassign_NonExistentTargetID_Returns4xx (0.01s)
    --- PASS: TestE2ESuite/TestGetDepartment_Children_DefaultDepthIs1 (0.01s)
    --- PASS: TestE2ESuite/TestGetDepartment_Children_Depth2 (0.01s)
    --- PASS: TestE2ESuite/TestGetDepartment_Children_MaxDepth5 (0.02s)
    --- PASS: TestE2ESuite/TestGetDepartment_DepthExceedsMax_Returns4xx (0.00s)
    --- PASS: TestE2ESuite/TestGetDepartment_Employees_SortedByCreatedAt (0.01s)
    --- PASS: TestE2ESuite/TestGetDepartment_ExcludesEmployees_WhenFalse (0.01s)
    --- PASS: TestE2ESuite/TestGetDepartment_IncludesEmployeesByDefault (0.01s)
    --- PASS: TestE2ESuite/TestGetDepartment_NotFound (0.00s)
    --- PASS: TestE2ESuite/TestGetDepartment_Success (0.00s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_CycleDetection_Returns4xx (0.01s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_DuplicateNameInSameParent_Returns409 (0.01s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_MoveCreatesNameConflict_Returns409 (0.01s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_MoveToNewParent (0.01s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_NotFound (0.00s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_RenameSuccess (0.00s)
    --- PASS: TestE2ESuite/TestUpdateDepartment_SelfParent_Returns4xx (0.00s)
PASS
ok  	github.com/pdkonovalov/hitalent-test-assignment-2026/tests/e2e	7.753s
```

# Migrations

To add new migration you need just add migration file to `./migrations` directory using `goose`. And it will be automatically applied on startup. Migration folder load to go binary using `embed`.

# Swagger

After change of code do:

```bash
make swagger
```
This command will generate new swagger docs:

```
2026/05/15 18:56:44 Generate swagger docs....
2026/05/15 18:56:44 Generate general API Info, search dir:./
2026/05/15 18:56:44 Generating http.CreateDepartmentRequest
2026/05/15 18:56:44 Generating http.CreateDepartmentResponse
2026/05/15 18:56:44 Generating http.ErrorResponse
2026/05/15 18:56:44 Generating http.Error
2026/05/15 18:56:44 Generating http.CreateEmployeeRequest
2026/05/15 18:56:44 Generating http.CreateEmployeeResponse
2026/05/15 18:56:44 Generating query.GetDepartment
2026/05/15 18:56:44 Generating query.Department
2026/05/15 18:56:44 Generating query.Employee
2026/05/15 18:56:44 Generating http.UpdateDepartmentRequest
2026/05/15 18:56:44 Generating http.UpdateDepartmentResponse
2026/05/15 18:56:44 create docs.go at docs/docs.go
2026/05/15 18:56:44 create swagger.json at docs/swagger.json
2026/05/15 18:56:44 create swagger.yaml at docs/swagger.yaml
```
