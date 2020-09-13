EXE=app
ZIP_FILE=app.zip
AWS_PROFILE ?= default
FUNC_NAME ?= grpc-proxy-go
IAM ?= IAM-is-not-defined
ROLE_NAME ?= lambda_basic_execution
ROLE_ARN ?= arn:aws:iam::$(IAM):role/$(ROLE_NAME)
GRPC_SERVER_ADDR ?= example.com:9090

.PHONY: dep
dep:
	@go get && go mod vendor

build-linux:
	@GO111MODULE=on GOOS=linux go build -mod=vendor -o $(EXE)

.PHONY: clean
clean:
	@rm -f $(EXE) $(ZIP_FILE)

.PHONY: zip
zip:
	@zip $(ZIP_FILE) $(EXE)

create-function: build-linux zip
	@aws lambda create-function \
	--profile $(AWS_PROFILE) \
	--function-name $(FUNC_NAME) \
	--runtime go1.x \
	--zip-file fileb://$(ZIP_FILE) \
	--handler $(EXE) \
	--role $(ROLE_ARN) \
	--environment Variables={SERVER_ADDR=$(GRPC_SERVER_ADDR)}

.PHONY: invoke
invoke:
	@aws lambda invoke \
	--profile $(AWS_PROFILE) \
	--function-name $(FUNC_NAME) \
	--payload '{"service":"lechuckroh.service.hello", "method":"Call", "key3":"v3"}' \
	out \
	--log-type Tail \
	--query 'LogResult' \
	--output text \
	| base64 -d
	@rm -f out

delete-function:
	@aws lambda delete-function \
	--function-name $(FUNC_NAME) \
	--profile $(AWS_PROFILE)

export-plantuml:
	@-cat docs/diagram.iuml | plantuml -pipe -tsvg > docs/diagram.svg
