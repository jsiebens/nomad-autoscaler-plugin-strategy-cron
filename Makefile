SHELL := bash
LDFLAGS := "-s -w"
.PHONY: all

.PHONY: %.zip
%.zip:
	touch $@

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o dist/cron

.PHONY: dist
dist:
	mkdir -p dist
	./scripts/dist.sh linux amd64
	./scripts/dist.sh linux arm64
	./scripts/dist.sh linux arm
	./scripts/dist.sh darwin amd64
	./scripts/dist.sh windows amd64
	cd dist && shasum -a 256 *.zip > cron_SHA256SUMS && cd ..