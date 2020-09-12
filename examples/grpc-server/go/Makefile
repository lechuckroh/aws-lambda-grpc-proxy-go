EXE=app

.PHONY: dep
dep:
	@go get && go mod vendor

.PHONY: build
build:
	@GO111MODULE=on go build -mod=vendor -o $(EXE)

.PHONY: clean
clean:
	@rm -f $(EXE)

