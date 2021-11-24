.PHONY: build

build:
	sam build

start: build
	sam local start-api