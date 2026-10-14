VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/edgedelta/edx/internal/cli.Version=$(VERSION)

# The agent skills are authored in ./skills and mirrored to a checkout of the
# public edgedelta/agent-skills repo (what `npx skills add` installs from). The
# release workflow runs the mirror on every tag; see mirror-skills below.
SKILLS_MIRROR ?= ../agent-skills

# Where the edgedelta monorepo is checked out. The dashboard JSON Schema is generated
# from the frontend's TypeScript types (the source of truth) and vendored into
# internal/dashboards, committed so `go install` and plain clones work.
MONOREPO_SRC ?= ../edgedelta
SCHEMA_SRC := $(MONOREPO_SRC)/web/src/modules/dashboards/versions/v4/schema/dashboard-v4.schema.json
SCHEMA_DST := internal/dashboards/dashboard-v4.schema.json
FIXTURES_DST := internal/dashboards/testdata/definitions

# The ANTLR query parsers, vendored from the same monorepo. These are generated Go files
# that depend only on the ANTLR runtime, so they compile here unmodified — the directory
# layout below matches the monorepo's so the sync stays a plain copy with no rewriting.
PARSERS_SRC := $(MONOREPO_SRC)/pkg/antlrcql
PARSERS_DST := internal/cql
PARSER_DIRS := logparser metric/parser formula/parser

# resource_accesses is derived from a definition by both the UI and edx, so edx's Go
# implementation is diffed against the frontend's actual output. This is that output.
ORACLE_DST := internal/dashboards/testdata/resource-accesses.json

.PHONY: build install test vet lint clean mirror-skills sync-dashboard-schema sync-cql-parsers sync-resource-accesses-oracle

build:
	go build -ldflags '$(LDFLAGS)' -o bin/edx .

install:
	go install -ldflags '$(LDFLAGS)' .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf bin

# Make $(SKILLS_MIRROR) an exact copy of ./skills, minus the Go embed file.
# The mirror's own .git and .github (its PR guard) are left alone. Committing
# and pushing the result is the caller's job (the release workflow does both).
mirror-skills:
	@test -d "$(SKILLS_MIRROR)/.git" || { echo "no agent-skills checkout at $(SKILLS_MIRROR); set SKILLS_MIRROR=/path/to/agent-skills"; exit 1; }
	rsync -a --delete --exclude=.git --exclude=.github --exclude='*.go' skills/ $(SKILLS_MIRROR)/

# Refresh the embedded dashboard schema after the frontend types change. Run
# `bun gen:dashboard-schema` in $(MONOREPO_SRC)/web first, then commit $(SCHEMA_DST).
# The fixtures are the dashboards the frontend ships, typed against the real
# definition, so `go test ./internal/dashboards` proves the schema has no false
# positives against real definitions.
sync-dashboard-schema:
	@test -f "$(SCHEMA_SRC)" || { echo "schema not found at $(SCHEMA_SRC); run 'bun gen:dashboard-schema' in the monorepo's web/, or set MONOREPO_SRC=/path/to/edgedelta"; exit 1; }
	cp $(SCHEMA_SRC) $(SCHEMA_DST)
	@echo "synced $(SCHEMA_DST)"
	@echo "reminder: refresh $(FIXTURES_DST) too if the shipped default dashboards changed"

# Refresh the vendored query parsers after the .g4 grammars change. Regenerate the Go
# files in the monorepo first (see pkg/antlrcql/README.md), then run this and commit.
# Only the .go files are copied; the .interp and .tokens files next to them are ANTLR
# tooling artifacts that nothing here compiles against.
sync-cql-parsers:
	@test -d "$(PARSERS_SRC)" || { echo "parsers not found at $(PARSERS_SRC); set MONOREPO_SRC=/path/to/edgedelta"; exit 1; }
	@for dir in $(PARSER_DIRS); do \
		rm -rf $(PARSERS_DST)/$$dir; \
		mkdir -p $(PARSERS_DST)/$$dir; \
		cp $(PARSERS_SRC)/$$dir/*.go $(PARSERS_DST)/$$dir/ || exit 1; \
		echo "synced $(PARSERS_DST)/$$dir"; \
	done
	@go build ./... && echo "parsers compile unmodified"

# Refresh the resource_accesses oracle from the frontend's own implementation, then run the
# differential test. Do this after any change to a widget's or variable's resolveResources
# in the monorepo: a failure here means edx and the UI would grant different access to a
# shared dashboard.
sync-resource-accesses-oracle:
	@test -d "$(MONOREPO_SRC)/web" || { echo "monorepo web/ not found at $(MONOREPO_SRC)/web; set MONOREPO_SRC=/path/to/edgedelta"; exit 1; }
	cd $(MONOREPO_SRC)/web && bun gen:resource-accesses > $(CURDIR)/$(ORACLE_DST)
	@echo "synced $(ORACLE_DST)"
	go test ./internal/dashboards -run TestResourceAccessesMatchesTheFrontend
