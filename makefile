default: cui

cui:
	$(MAKE) -C ./cui

cur_release:
	$(MAKE) -C ./cui release

release:
	wails build
	cp ./build/bin/sttool.exe ../twichevent.exe

test:
	go test -v ./backend

auto:
	autocmd -v -t '.*\.go' -- make test

gen:
	AppClientID=$(AppClientID) AppClientSecret=$(AppClientSecret) go run ./tool/gen.go

.PHONY: default release cui cui_release test auto gen
