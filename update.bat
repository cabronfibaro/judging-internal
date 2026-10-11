@echo off
cd /d "%~dp0"

git pull --ff-only origin main

if errorlevel 1 (
    echo Update failed.
    pause
    exit /b 1
)

call npm install
call npm run build

pause