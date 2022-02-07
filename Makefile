.PHONY: validate-local

validate-local:
	sam validate --template template.local.yml

.PHONY: build-local

build-local: validate-local
	sam build --template template.local.yml

.PHONY: test

test:
	go test ./... -v

.PHONY: coverage

coverage:
	go test ./... -cover -coverprofile=coverage
	go tool cover -html=coverage

.PHONY:

mod:
	go mod vendor

.PHONY: validate

validate:
	sam validate --template template.yml

.PHONY: build

build: validate
	sam build --template template.yml

.PHONY: api

api: build-local
	sam local start-api --env-vars env.json

.PHONY: package

package: build
	sam package \
		--s3-bucket password-caddy-cloudformation-artifacts-dev \
		--s3-prefix v1 \
		--region us-east-2 \
		--output-template-file packaged-dev.yml

	C:\"Program Files"\Git\usr\bin\bash.exe ./scripts/set_env.sh

.PHONY: deploy

deploy: package
	sam deploy \
		--stack-name password-caddy-api-dev-v1 \
		--template packaged-dev.yml \
		--capabilities CAPABILITY_IAM \
		--region us-east-2 \
		--s3-bucket password-caddy-cloudformation-artifacts-dev \
		--s3-prefix v1 \
		--no-fail-on-empty-changeset \
		--role-arn arn:aws:iam::480277082058:role/password-caddy-cloudformation-execution-role \
		--parameter-overrides 'ENV=dev ACCOUNTID=480277082058'

.PHONY: delete

delete:
	sam delete \
		--stack-name password-caddy-api-dev \
		--no-prompts \
		--region us-east-2

.PHONY: login

invoke-login: build-local
	sam local invoke LoginFunction --env-vars env.json

.PHONY: controller

controller:
	C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe .\scripts\CreateController.ps1
