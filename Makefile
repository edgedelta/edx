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

.PHONY: build install test vet lint clean sync-skills sync-dashboard-schema

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
