BIN_Main := twitch_event_cmdline.exe
BIN_Cmd  := twitch_command.exe

default: build

build: $(BIN_Main) $(BIN_Cmd)

$(BIN_Main): ../backend/*.go cui.go
	go build -o $@ cui.go

$(BIN_Cmd): ../backend/*.go cmd.go
	go build -o $@ cmd.go

run: $(BIN_Main)
	./$(BIN_Main)

release: $(BIN_Main) $(BIN_Cmd)
	cp $^ ../../../test/

.PHONY: default build release
