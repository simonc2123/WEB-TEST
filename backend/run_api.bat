@echo off
echo Running main.go to update the database...
go run main.go
if %errorlevel% neq 0 (
    echo Error running main.go. Exiting.
    exit /b %errorlevel%
)
echo Running api_in.go to start the API...
go run api_in.go