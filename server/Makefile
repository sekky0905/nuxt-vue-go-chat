.PHONY: deps
deps: tools
	go get github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: tools
tools:
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
	go get github.com/kisielk/errcheck

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	find . -name '*.go' | grep -v 'vendor' | xargs gofmt -l
	find . -name '*.go' | grep -v 'vendor' | xargs goimports -l
	go list ./... | grep -v 'vendor' | xargs golint
	go list ./... | grep -v 'vendor' | xargs go vet
	go list ./... | grep -v 'vendor' | grep -v 'mock' | xargs errcheck -asserts -blank -exclude errcheck_excludes.txt

.PHONY:run
run:
	cd ../client/nuxt-vue-go-chat && yarn lint && yarn build && cd ../.. && docker-compose up
