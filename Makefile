define run
	@GOGC=off go build -o ./bin/$(1) ./cmd/$(1).go
	@./bin/$(1) --config=$(2)
endef

run:
	$(call run,app,)

run-conf:
	$(call run,app,./etc/app.conf)

define build
	@go build -o ./bin/$(1) ./cmd/$(1).go
endef

build:
	$(call build,app)

test:
	@go test -v ./...
