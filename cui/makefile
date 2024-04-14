BIN := csttool.exe

default: build

build: $(BIN)

$(BIN): ../backend/*.go
	go build -o $@

run: $(BIN)
	./$(BIN)

release: $(BIN)
	cp ./$(BIN) ../../../../test/

.PHONY: default build release
