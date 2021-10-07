@echo off
for /l %%i in (1001,1,1100) do (
	start student.exe U%%i %%i SEAT%%i
)