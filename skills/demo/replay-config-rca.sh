#!/usr/bin/env bash
# Replays a real config-change investigation (Claude Code + ed-pipelines skill).
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
type_user "the dogfood fleet started misbehaving on the evening of Mar 4 - was it a config change?"
sleep 0.6

say "Checking pipeline deploy history with the ed-pipelines skill."
sleep 0.9

tool "edx pipelines list --keyword dogfood --output table --columns tag,fleet_type,environment"
sleep 0.9
out "TAG              FLEET_TYPE  ENVIRONMENT"
outc "dogfood-confv3   Edge        Kubernetes"
printf '\n'
sleep 1.2

tool "edx pipelines history 70669202-2ca9-49e2-8e7a-862cbf2eb093"
sleep 0.9
out "WHEN (UTC)       AUTHOR                 STATUS    DESCRIPTION"
outc "Mar 04 22:58     fatih@edgedelta.com    Deployed  Connect cloud trail mp to cloud trail pack."
outc "Mar 04 21:34     fatih@edgedelta.com    Deployed"
outc "Mar 03 10:46     orhan@edgedelta.com    Deployed"
printf '\n'
sleep 1.6

say "Top suspect: the Mar 04 22:58 UTC deploy (\"Connect cloud trail mp to cloud trail pack\")."
printf '  It is the last change before onset; the prior deploy (21:34) ran clean for over an hour.\n'
printf '  To mitigate, roll back to the previous version:\n'
printf '    %sedx pipelines deploy 70669202-... <previous-version> --yes%s\n' "$BOLD" "$RESET"
printf '  then verify recovery with  %sedx fleet deployments%s  and  %sedx health problems%s.\n' "$BOLD" "$RESET" "$BOLD" "$RESET"
sleep 5
