@echo off
for /l %%i in (1001,1,1100) do (
	start student.exe U202%%i 202%%i SEAT%%i
	timeout /T 1
)