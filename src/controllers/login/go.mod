module password-caddy/login

go 1.17

require github.com/aws/aws-lambda-go v1.27.0

require password-caddy/response v1.0.0

require password-caddy/util v1.0.0

replace password-caddy/response => ../../lib/response

replace password-caddy/util => ../../lib/util
