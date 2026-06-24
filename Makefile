VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/edgedelta/edx/internal/cli.Version=$(VERSION)

# Where the agent-skills source repo is checked out. The skills are vendored
# into internal/skills/data (committed, so `go install` and plain clones work);
# `make sync-skills` regenerates that copy. agent-skills stays the source of truth.
SKILLS_SRC ?= ../agent-skills
SKILLS_DST := internal/skills/data

.PHONY: build install test vet lint clean sync-skills

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
