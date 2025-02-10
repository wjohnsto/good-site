MAKEFLAGS += --no-print-directory

## ---------------------------------------------------------------------------
## | The purpose of this Makefile is to provide all the functionality needed |
## | to install, build, run, and deploy static sites.                        |
## ---------------------------------------------------------------------------

help:              ## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

install:           ## Install all dependencies only if they don't exist
	@$(MAKE) install-deps

install-deps: install-inotifywait install-lightningcss install-minhtml install-purgecss install-templ
	@go mod tidy

install-inotifywait:
	@type inotifywait >/dev/null 2>&1 || sudo apt install inotify-tools

install-lightningcss:
	@type lightningcss >/dev/null 2>&1 || npm i -g lightningcss-cli@latest

install-minhtml:
	@type minhtml >/dev/null 2>&1 || cargo install minhtml

install-purgecss:
	@type purgecss >/dev/null 2>&1 || npm i -g purgecss@latest

install-templ:
	@type templ >/dev/null 2>&1 || go install github.com/a-h/templ/cmd/templ@latest

install-f:         ## Force install dependencies for building the site
	@$(MAKE) install-deps
	@npm i -g lightningcss-cli@latest purgecss@latest
	@cargo install minhtml
	@sudo apt install inotify-tools
	@go install github.com/a-h/templ/cmd/templ@latest

copy-assets:
	@mkdir -p public
	@cp -r assets/* public
	@mkdir -p nonsense
	@cp -r assets/* nonsense
	@cp -r assets/favicon/favicon.ico public
	@cp -r assets/favicon/favicon.ico nonsense

format:            ## Format .go and .templ files
	@gofmt -s -w .
	@templ fmt .

build:             ## Build the static site
	@$(MAKE) build-all

build-all: clean install copy-assets build-templ build-static build-purgecss build-lightningcss build-mergecss build-minhtml build-nonsense

build-static:
	@echo "Building static pages"
	@go run cmd/static/main.go

build-purgecss:
	@echo "Purging unused CSS"
	@find ./public -type d -exec test -e '{}'/styles.css -a -e '{}'/index.html \; | xargs -P 10 -I % purgecss -css "%/styles.css" --content "%/index.html" --output "%/styles.css"

build-lightningcss:
	@echo "Minifying all CSS"
	@find ./public -type f -name "*.css" | xargs -P 10 -I % lightningcss --minify "%" -o "%"

build-mergecss:
	@echo "Merging CSS into HTML files"
	@go run cmd/css/main.go

build-minhtml:
	@echo "Minifying HTML"
	@find ./public -type f -name "*.html" | xargs -P 10 -I % minhtml \
				--keep-closing-tags \
				--do-not-minify-doctype \
				--ensure-spec-compliant-unquoted-attribute-values \
				--keep-html-and-head-opening-tags \
				--minify-css \
				--minify-js \
				-o % \
				%

build-nonsense:
	@echo "Generating nonsense"
	@go run cmd/nonsense/main.go --percent 20

build-templ:
	@templ generate

prune:             ## Prune unused docker images, volumes, and builder cache
	@docker image prune -a -f
	@docker volume prune -a -f
	@docker builder prune -a -f

deploy:            ## Deploy to fly.io using a local build
	@$(MAKE) deploy-local

deploy-local: prune build
	@docker build -t registry.fly.io/good-site:latest .
	@docker push registry.fly.io/good-site:latest
	@fly deploy -i registry.fly.io/good-site:latest

deploy-remote:     ## Deploy to fly.io using fly.io's cloud builder
	@$(MAKE) deploy-remote-f

deploy-remote-f: prune static
	@fly deploy

dev:               ## Watch files for changes and serve the site on :8080
	@$(MAKE) watch | $(MAKE) serve

serve:             ## Run the local file server on :8080
	@go run cmd/serve/main.go

serve-nonsense:    ## Run the local file server for nonsense files on :8080
	@go run cmd/serve/main.go nonsense

watch:             ## Watch directories for changes and rebuild the site
	@while true; do \
		$(MAKE) build; \
		inotifywait --quiet -qre close_write ./pkg ./markdown; \
	done

list-updates:      ## List updates to go dependencies
	@go list -m -u all

update:            ## Update all go dependencies and install
	@go get -u ./...
	@$(MAKE) install

docker:            ## Rebuild and run docker container
	@docker compose down
	@$(MAKE) prune
	@$(MAKE) build
	@docker compose up -d
	
clean:             ## Clean all public files and templ build files
	@rm -rf public/*
	@rm -rf nonsense/*
	@find . -type f -name '*_templ.go' -delete
	@find . -type f -name '*_templ.txt' -delete
