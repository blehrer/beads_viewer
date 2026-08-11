package agents

import (
	"strings"

	"github.com/Dicklesworthstone/beads_viewer/pkg/beadscli"
)

// READMEAgentBlurb is the copy-paste agent sidecar blurb shown in README.md.
// It includes README-only detail (count semantics, liveness, --brief) that the
// shorter AGENTS.md injection blurb omits.
const READMEAgentBlurb = `### Using bv as an AI sidecar

bv is a graph-aware triage engine for Beads projects (` + "`.beads/issues.jsonl`" + ` in current ` + "`br`" + ` workspaces, with ` + "`.beads/beads.jsonl`" + ` supported for legacy/` + "`bd`" + ` workspaces). Instead of parsing JSONL or hallucinating graph traversal, use robot flags for deterministic, dependency-aware outputs with precomputed metrics (PageRank, betweenness, critical path, cycles, HITS, eigenvector, k-core).

**Scope boundary:** bv handles *what to work on* (triage, priority, planning). For agent-to-agent coordination (messaging, work claiming, file reservations), use [MCP Agent Mail](https://github.com/Dicklesworthstone/mcp_agent_mail).

**⚠️ CRITICAL: Use ONLY ` + "`--robot-*`" + ` flags. Bare ` + "`bv`" + ` launches an interactive TUI that blocks your session.**

#### The Workflow: Start With Triage

**` + "`bv --robot-triage`" + ` is your single entry point.** It returns everything you need in one call:
- ` + "`quick_ref`" + `: at-a-glance counts + top 3 picks
- ` + "`recommendations`" + `: ranked actionable items with scores, reasons, unblock info
- ` + "`quick_wins`" + `: low-effort high-impact items
- ` + "`blockers_to_clear`" + `: items that unblock the most downstream work
- ` + "`project_health`" + `: status/type/priority distributions, graph metrics
- ` + "`commands`" + `: copy-paste shell commands for next steps

**Count semantics (strict since #165):**
- ` + "`quick_ref.open_count`" + ` / ` + "`project_health.counts.open`" + ` — issues with status exactly ` + "`open`" + `; always equals ` + "`counts.by_status.open`" + `
- ` + "`quick_ref.blocked_count`" + ` / ` + "`counts.blocked`" + ` — issues with status exactly ` + "`blocked`" + `; always equals ` + "`counts.by_status.blocked`" + `
- ` + "`quick_ref.in_progress_count`" + ` — status exactly ` + "`in_progress`" + `
- ` + "`counts.closed`" + ` — closed-like issues (` + "`closed`" + ` + ` + "`tombstone`" + `)
- ` + "`quick_ref.not_closed_count`" + ` / ` + "`counts.not_closed`" + ` — every non-closed issue (` + "`open`" + `+` + "`in_progress`" + `+` + "`blocked`" + `+` + "`deferred`" + `); this is the pre-#165 meaning of ` + "`open_count`" + `
- ` + "`quick_ref.actionable_count`" + ` / ` + "`counts.actionable`" + ` — non-closed issues with no open blocking dependencies (ready to work now)
- ` + "`quick_ref.not_actionable_count`" + ` / ` + "`counts.dependency_blocked`" + ` — non-closed issues blocked by open dependencies regardless of status; this is the pre-#165 meaning of ` + "`blocked_count`" + `
- Partition invariant: ` + "`not_closed == actionable + not_actionable`" + ` (every non-closed issue is exactly one of the two)

**Liveness (#166):** the git-history prologue of ` + "`--robot-triage`" + ` is bounded (default 10s; tune via ` + "`--robot-history-timeout-ms <ms>`" + ` or ` + "`BV_ROBOT_HISTORY_TIMEOUT_MS`" + `, ` + "`0`" + ` = unbounded). On timeout the in-flight git subprocess is killed and triage proceeds without history; ` + "`meta.history_status`" + ` reports ` + "`ok`" + `, ` + "`error`" + `, or ` + "`timeout`" + ` (omitted when history was not attempted).

bv --robot-triage        # THE MEGA-COMMAND: start here
bv --robot-next          # Minimal: just the single top pick + claim command

# Compact triage (#183): only decision-relevant fields — id, title, status,
# assignee, blocked_by, unblocks, score — plus quick_ref/quick_wins/blockers.
# Cuts payload size by ~80% versus the full output.
bv --robot-triage --brief

# Token-optimized output (TOON) for lower LLM context usage:
bv --robot-triage --format toon
export BV_OUTPUT_FORMAT=toon
bv --robot-next

Before claiming, verify the current bead state with ` + "`br show <id> --json`" + ` or
` + "`br ready --json`" + `. ` + "`recommendations`" + ` can include graph-important blocked or
assigned work; only ` + "`quick_ref.top_picks`" + ` and non-empty ` + "`claim_command`" + ` fields
represent claimable work.

#### Other Commands

**Planning:**
| Command | Returns |
|---------|---------|
| ` + "`--robot-plan`" + ` | Parallel execution tracks with ` + "`unblocks`" + ` lists |
| ` + "`--robot-priority`" + ` | Priority misalignment detection with confidence |

**Graph Analysis:**
| Command | Returns |
|---------|---------|
| ` + "`--robot-insights`" + ` | Full metrics: PageRank, betweenness, HITS (hubs/authorities), eigenvector, critical path, cycles, k-core, articulation points, slack |
| ` + "`--robot-label-health`" + ` | Per-label health: ` + "`health_level`" + ` (healthy\|warning\|critical), ` + "`velocity_score`" + `, ` + "`staleness`" + `, ` + "`blocked_count`" + ` |
| ` + "`--robot-label-flow`" + ` | Cross-label dependency: ` + "`flow_matrix`" + `, ` + "`dependencies`" + `, ` + "`bottleneck_labels`" + ` |
| ` + "`--robot-label-attention [--attention-limit=N]`" + ` | Attention-ranked labels by: (pagerank × staleness × block_impact) / velocity |

**History & Change Tracking:**
| Command | Returns |
|---------|---------|
| ` + "`--robot-history`" + ` | Bead-to-commit correlations: ` + "`stats`" + `, ` + "`histories`" + ` (per-bead events/commits/milestones), ` + "`commit_index`" + ` |
| ` + "`--robot-diff --diff-since <ref>`" + ` | Changes since ref: new/closed/modified issues, cycles introduced/resolved |

**Other Commands:**
| Command | Returns |
|---------|---------|
| ` + "`--robot-burndown <sprint>`" + ` | Sprint burndown, scope changes, at-risk items |
| ` + "`--robot-forecast <id\\|all>`" + ` | ETA predictions with dependency-aware scheduling |
| ` + "`--robot-alerts`" + ` | Stale issues, blocking cascades, priority mismatches |
| ` + "`--robot-suggest`" + ` | Hygiene: duplicates, missing deps, label suggestions, cycle breaks |
| ` + "`--robot-graph [--graph-format=json\\|dot\\|mermaid]`" + ` | Dependency graph export |
| ` + "`--export-graph <file.html>`" + ` | Self-contained interactive HTML visualization |

#### Scoping & Filtering

bv --robot-plan --label backend              # Scope to label's subgraph
bv --robot-insights --as-of HEAD~30          # Historical point-in-time
bv --recipe actionable --robot-plan          # Pre-filter: ready to work (no blockers)
bv --recipe high-impact --robot-triage       # Pre-filter: top PageRank scores
bv --robot-triage --robot-triage-by-track    # Group by parallel work streams
bv --robot-triage --robot-triage-by-label    # Group by domain

#### Understanding Robot Output

**All robot JSON includes:**
- ` + "`data_hash`" + ` — Fingerprint of the source JSONL issue file (verify consistency across calls)
- ` + "`status`" + ` — Per-metric state: ` + "`computed|approx|timeout|skipped`" + ` + elapsed ms
- ` + "`as_of`" + ` / ` + "`as_of_commit`" + ` — Present when using ` + "`--as-of`" + `; contains ref and resolved SHA

**Two-phase analysis:**
- **Phase 1 (instant):** degree, topo sort, density — always available immediately
- **Phase 2 (async, 500ms timeout):** PageRank, betweenness, HITS, eigenvector, cycles — check ` + "`status`" + ` flags

**For large graphs (>500 nodes):** Some metrics may be approximated or skipped. Always check ` + "`status`" + `.

#### jq Quick Reference

bv --robot-triage | jq '.quick_ref'                        # At-a-glance summary
bv --robot-triage | jq '.recommendations[0]'               # Top recommendation
bv --robot-plan | jq '.plan.summary.highest_impact'        # Best unblock target
bv --robot-insights | jq '.status'                         # Check metric readiness
bv --robot-insights | jq '.Cycles'                         # Circular deps (must fix!)
bv --robot-label-health | jq '.results.labels[] | select(.health_level == "critical")'

**Performance:** Phase 1 instant, Phase 2 async (500ms timeout). Prefer ` + "`--robot-plan`" + ` over ` + "`--robot-insights`" + ` when speed matters. Results cached by data hash.

Use bv instead of parsing Beads JSONL directly—it computes PageRank, critical paths, cycles, and parallel tracks deterministically.
`

// READMEAgentInstructions returns READMEAgentBlurb with the active Beads CLI substituted.
func READMEAgentInstructions() string {
	tool := beadscli.Tool()
	if tool == "br" {
		return READMEAgentBlurb
	}
	return strings.NewReplacer(
		"`br show", "`"+tool+" show",
		"`br ready", "`"+tool+" ready",
	).Replace(READMEAgentBlurb)
}
