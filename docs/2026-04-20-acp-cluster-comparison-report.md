---
title: ACP Cluster Comparison Report
date: 2026-04-20
status: final
---

# ACP Cluster Comparison Report

## Summary

This report compares two earlier manual `prtags` duplicate groups against the current `pr-search-cli` issue-cluster output for `openclaw/openclaw`.

The result is mixed:

- the earlier manual `sessions_spawn` duplicate group looks more correct than the current `pr-search-cli` clusters
- the earlier manual `/acp` duplicate group caught a real duplicate pair, but it was too narrow
- the current `pr-search-cli` `/acp` cluster also caught a real duplicate pair, but it was too narrow in a different way

## Manual `prtags` Groups

### `neutral-polliwog-g6xx`

Title:

`ACP duplicate: sessions_spawn ACP-only field stripping`

Members:

- `#56342` `Fix sessions_spawn for subagent runtime with ACP-only fields`
- `#56438` `fix: strip ACP-only fields silently when runtime=subagent`
- `#68397` `fix(sessions_spawn): silently strip ACP-only fields for runtime=subagent`

### `great-loon-t2te`

Title:

`ACP duplicate: bound-session /acp command dispatch`

Members:

- `#66407` `fix(acp): bypass ACP dispatch for /acp text commands in bound threads`
- `#68617` `fix(acp): keep /acp commands local in bound sessions`

## `pr-search-cli` Results

Commands used:

- `uvx pr-search-cli@latest issues for-pr 56342`
- `uvx pr-search-cli@latest issues for-pr 56438`
- `uvx pr-search-cli@latest issues for-pr 68397`
- `uvx pr-search-cli@latest issues for-pr 66407`
- `uvx pr-search-cli@latest issues for-pr 68617`
- `uvx pr-search-cli@latest issues show cluster-65248-2`

Observed results:

- `#56342` was not found in any issue cluster
- `#56438` was not found in any issue cluster
- `#68397` was not found in any issue cluster
- `#66407` was not found in any issue cluster
- `#68617` was found in `cluster-65248-2`

`cluster-65248-2` contains:

- canonical PR `#65248` `fix(acp): bypass bound slash commands to local handlers`
- member PR `#68617` `fix(acp): keep /acp commands local in bound sessions`

## Comparison

### `sessions_spawn` field-stripping group

The manual `prtags` group is more correct here.

All three PRs describe the same bug:

- `sessions_spawn`
- `runtime="subagent"`
- ACP-only fields like `streamTo` or `resumeSessionId`
- fix by silently stripping or ignoring those fields instead of failing

The titles, summaries, and PR bodies all point at the same root cause and the same fix shape.

The fact that `pr-search-cli` does not place these PRs into an issue cluster looks like a miss of the issue-clustering method, not evidence that the PRs are unrelated.

### `/acp` command-routing group

Both systems caught part of the truth.

The earlier manual `prtags` group paired:

- `#66407`
- `#68617`

That pair is clearly valid.

Both PRs describe the same routing bug:

- `/acp ...` commands
- inside ACP-bound conversations or threads
- were being forwarded into the ACP runtime instead of staying on the local OpenClaw command path

The `pr-search-cli` cluster paired:

- `#65248`
- `#68617`

That pair is also clearly valid.

`#65248` fixes the same family of bug:

- recognized slash or ACP control commands in ACP-bound sessions
- should bypass ACP dispatch
- and should stay on the local handler path

## Objective Judgment

The best current reading is:

- the manual `sessions_spawn` cluster is better than the current `pr-search-cli` result
- neither system has the full best `/acp` duplicate set yet

The most faithful `/acp` duplicate set appears to be:

- `#65248`
- `#66407`
- `#68617`

Why:

- `#65248` and `#68617` are clearly in the same family of local-command bypass fixes for ACP-bound sessions
- `#66407` and `#68617` are also clearly in the same family and even name the same bypass logic in `shouldBypassAcpDispatchForCommand()`
- the three together describe overlapping fixes to the same command-routing bug family

So the two narrow clusters are each missing one real related PR.

## Recommended Update

If these groupings are going to be curated in `prtags`, the current best correction is:

1. keep the `sessions_spawn` manual duplicate cluster as is
2. broaden the `/acp` duplicate cluster to include `#65248`

That would make the curated duplicate view more faithful than either narrow pair alone.
