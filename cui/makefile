BIN := csttool.exe

default: build

build: $(BIN)

$(BIN):
	go build -o $@

run:
	go run cui.go

release: $(BIN)
	cp ./$(BIN) ../../../../test/

.PHONY: default build release
