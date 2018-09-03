GO_BACKEND_DIR=$(CURDIR)/golang-backend
DIST_DIR=$(CURDIR)/dist/
STATIC_FILES_DIR=$(DIST_DIR)/static
GO_BUILD=go build -tags embed
INDEX_PAGE=example-gauges.html

.PHONY: clean dist

dist/backend:
	cd $(GO_BACKEND_DIR) && $(GO_BUILD) -o $(DIST_DIR)/speedtest

dist/static:
	mkdir -p $(STATIC_FILES_DIR)
	cp speedtest_worker.js $(STATIC_FILES_DIR)
	cp speedtest_worker.min.js $(STATIC_FILES_DIR)
	cp *.html $(STATIC_FILES_DIR)
	cp $(STATIC_FILES_DIR)/$(INDEX_PAGE) $(STATIC_FILES_DIR)/index.html
	cd $(STATIC_FILES_DIR) && statik -tags embed -src ./ -dest $(GO_BACKEND_DIR)

dist: dist/static dist/backend

clean:
	rm -rf $(DIST_DIR)
	rm -rf $(GO_BACKEND_DIR)/statik
