INDEX_PAGE=example-gauges.html

GO_BACKEND_DIR=$(CURDIR)/golang-backend
DIST_DIR=$(CURDIR)/dist/
STATIC_FILES_DIR=$(DIST_DIR)/static
GO_BUILD=CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-s -w" -tags embed

.PHONY: clean dist dist-linux-amd64
.DEFAULT_GOAL := dist

dist/backend:
	cd $(GO_BACKEND_DIR) && $(GO_BUILD) -o $(DIST_DIR)/speedtest

dist/backend-linux-amd64:
	cd $(GO_BACKEND_DIR) && GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(DIST_DIR)/speedtest

dist/backend-linux-mipsle:
	cd $(GO_BACKEND_DIR) && GOOS=linux GOARCH=mipsle $(GO_BUILD) -o $(DIST_DIR)/speedtest

dist/static:
	@mkdir -p $(STATIC_FILES_DIR)
	@cp speedtest_worker.js $(STATIC_FILES_DIR)
	@cp speedtest_worker.min.js $(STATIC_FILES_DIR)
	@cp *.html $(STATIC_FILES_DIR)
	@cp $(STATIC_FILES_DIR)/$(INDEX_PAGE) $(STATIC_FILES_DIR)/index.html
	@cd $(STATIC_FILES_DIR) && statik -tags embed -src ./ -dest $(GO_BACKEND_DIR)

dist: dist/static dist/backend

dist-linux-amd64: dist/static dist/backend-linux-amd64

dist-linux-mipsle: dist/static dist/backend-linux-mipsle

clean:
	@rm -rf $(DIST_DIR)
	@rm -rf $(GO_BACKEND_DIR)/statik
