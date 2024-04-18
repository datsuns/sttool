default: cui

cui:
	$(MAKE) -C ./cui

release:
	$(MAKE) -C ./cui release

gen:
	AppClientID=$(AppClientID) AppClientSecret=$(AppClientSecret) go run ./tool/gen.go

.PHONY: default cui release
