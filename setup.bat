@echo off
setlocal

node -v >nul 2>&1

if %ERRORLEVEL% neq 0 (
    echo Node.js is not installed. Installing via Winget...

    winget install -e --id OpenJS.NodeJS --accept-source-agreements --accept-package-agreements

    if %ERRORLEVEL% neq 0 (
        echo Failed to install Node.js.
        pause
        exit /b 1
    )

    echo Node.js installed. Refreshing PATH...

    set "PATH=%PATH%;C:\Program Files\nodejs"
)

node -v

echo.
echo Installing dependencies...
call npm install

echo.
echo Building application...
call npm run build

echo.
echo Running database migrations...
call npm run db:migrate

if %ERRORLEVEL% neq 0 (
    echo Database migration failed.
    pause
    exit /b 1
)

echo Database migration completed.

pause