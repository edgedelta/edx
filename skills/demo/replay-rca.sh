#!/usr/bin/env bash
# Replays a real investigation session (Claude Code + ed-investigate skill).
# The edx outputs below were captured from live runs and are replayed with
# scripted timing so the clip stays short. No credentials appear.

DIM=$'\e[2m'; BOLD=$'\e[1m'; RESET=$'\e[0m'; GREEN=$'\e[32m'; MAGENTA=$'\e[35m'

type_user() {
  printf '%s> ' "$BOLD"
  local s="$1"
  for ((i = 0; i < ${#s}; i++)); do
    printf '%s' "${s:i:1}"
    sleep 0.028
  done
  printf '%s\n\n' "$RESET"
}

say() { printf '%s⏺%s %s\n\n' "$MAGENTA" "$RESET" "$1"; }
tool() { printf '%s⏺%s %sBash%s%s(%s)%s\n' "$GREEN" "$RESET" "$BOLD" "$RESET" "$DIM" "$1" "$RESET"; }
out() { printf '  %s⎿  %s%s\n' "$DIM" "$1" "$RESET"; }
outc() { printf '  %s   %s%s\n' "$DIM" "$1" "$RESET"; }

clear
type_user "admin-api error volume looks elevated — find the root cause"
sleep 0.6

say "Investigating with the ed-investigate skill: surging patterns first, then quantify the onset."
sleep 0.9

tool "edx patterns list --negative -q 'service.name:\"admin-api-deployment\"' --offset 24h --output table"
sleep 0.9
out "COUNT  DELTA   SENTIMENT  PATTERN"
outc "12     1100    -6         * ERROR selfprofiler* Failed to save cpu profile to S3 err * failed to upload"
outc "12     300     -2         * ERROR middleware* Failed to validate OAuth access token err no ID token found"
outc "19     0       -1         * WARN telemetry* Failed to save recent search query invalid queryScope"
printf '\n'
sleep 1.6

tool "edx logs graph -q 'service.name:\"admin-api-deployment\" AND severity_text:\"ERROR\"' --lookback 6h"
sleep 0.9
out "13,804 errors / 6h - steady ~5 per bucket until 14:00 UTC,"
outc "then a sustained step to ~50 from 14:27 UTC onward (10x)."
printf '\n'
sleep 1.6

say "Root cause: error rate stepped up 10x at ~14:27 UTC, driven by two new signatures:"
printf '  selfprofiler S3 uploads failing (%s+1100%%%s) and OAuth token validation failures (+300%%).\n' "$BOLD" "$RESET"
printf '  The onset is a clean step, not a ramp - consistent with a deploy or permission change\n'
printf '  at ~14:27 UTC. Check what shipped to admin-api then, starting with profiler S3 access.\n'
sleep 5
