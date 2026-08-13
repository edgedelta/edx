VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/edgedelta/edx/internal/cli.Version=$(VERSION)

# Where the agent-skills source repo is checked out. The skills are vendored
# into internal/skills/data (committed, so `go install` and plain clones work);
# `make sync-skills` regenerates that copy. agent-skills stays the source of truth.
SKILLS_SRC ?= ../agent-skills
SKILLS_DST := internal/skills/data

# Where the edgedelta monorepo is checked out. The dashboard JSON Schema is generated
# from the frontend's TypeScript types (the source of truth) and vendored into
# internal/dashboards, committed for the same reason as the skills above.
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

.PHONY: build install test vet lint clean sync-skills sync-dashboard-schema sync-cql-parsers sync-resource-accesses-oracle

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

# Regenerate the vendored copy of the agent skills embedded in the binary.
# Run after skills change in $(SKILLS_SRC), then commit internal/skills/data.
sync-skills:
	@test -d "$(SKILLS_SRC)" || { echo "skills source not found at $(SKILLS_SRC); set SKILLS_SRC=/path/to/agent-skills"; exit 1; }
	rm -rf $(SKILLS_DST)
	mkdir -p $(SKILLS_DST)
	cp -R $(SKILLS_SRC)/ed-* $(SKILLS_DST)/
	@echo "synced skills from $(SKILLS_SRC) into $(SKILLS_DST):"
	@ls -1 $(SKILLS_DST)

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
