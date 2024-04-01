default: cui

cui:
	$(MAKE) -C ./cui

release:
	$(MAKE) -C ./cui release

.PHONY: default cui release
