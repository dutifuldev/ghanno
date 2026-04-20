# AGENTS.md

## Product Boundary

`prtags` depends on `ghreplica`.

`ghreplica` is the GitHub-shaped mirror.

`prtags` is the curation layer on top of mirrored GitHub objects.

That boundary should stay explicit in both implementation and explanations.

- mirrored GitHub resources should come from `ghreplica`
- groups, annotations, field definitions, and target projections belong to `prtags`
- if `prtags` needs product-specific behavior, implement it in `prtags`, not in `ghreplica`

## Read Behavior

`group get` is refs-only by default.

Metadata is opt-in:

- CLI: `--include-metadata`
- HTTP: `?include=metadata`

Metadata reads should stay cache-first through `target_projections`.

## Write Behavior

Group membership writes should succeed from stable refs first.

They should not block on live `ghreplica` fetches before the write succeeds.

Projection refresh should happen after the write through background jobs.

## Documentation Convention

Follow SimpleDoc for repository documentation.

- General, non-dated documents should use capitalized filenames with underscores.
- Dated documents should live under `docs/` and use ISO date prefixes with lowercase kebab-case filenames.
