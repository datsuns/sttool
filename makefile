default: cui

cui:
	$(MAKE) -C ./cui

release:
	$(MAKE) -C ./cui release

test:
	go test -v ./backend

gen:
	AppClientID=$(AppClientID) AppClientSecret=$(AppClientSecret) go run ./tool/gen.go

.PHONY: default cui release
