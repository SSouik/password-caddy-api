.PHONY: build

build:
	sam build

start: build
	sam local start-api

controller:
	C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe .\scripts\CreateController.ps1