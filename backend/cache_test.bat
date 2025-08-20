@echo off
REM Save as: cache_test.bat
REM This batch file runs the PowerShell cache test script

title Cache Performance Test

echo Starting Cache Performance Test...
echo.

powershell.exe -ExecutionPolicy Bypass -File "%~dp0cache_test.ps1"

pause