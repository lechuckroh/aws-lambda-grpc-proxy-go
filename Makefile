.PHONY: vendor build-linux zip create-func

EXE=app
ZIP_FILE=app.zip
AWS_PROFILE ?= default
FUNC_NAME ?= proxy-lambda
IAM ?= IAM-is-not-defined
ROLE_NAME ?= lambda_basic_execution
ROLE_ARN ?= arn:aws:iam::$(IAM):role/$(ROLE_NAME)

dep:
	@go get && go mod vendor

build-linux:
	@GO111MODULE=on GOOS=linux go build -mod=vendor -o $(EXE)

clean:
	@rm -f $(EXE) $(ZIP_FILE)

zip:
	@zip $(ZIP_FILE) $(EXE)

create-function: build-linux zip
	@aws lambda create-function \
	--profile $(AWS_PROFILE) \
	--function-name $(FUNC_NAME) \
	--runtime go1.x \
	--zip-file fileb://$(ZIP_FILE) \
	--handler $(EXE) \
	--role $(ROLE_ARN)

invoke:
	@aws lambda invoke \
	--profile $(AWS_PROFILE) \
	--function-name $(FUNC_NAME) \
	--payload '{"key1":"v1", "key2":"v2", "key3":"v3"}' \
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
