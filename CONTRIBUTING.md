# Contributing

## Merge strategy

All PRs are squash-merged. The PR title becomes the commit message on `main`, so it must follow the conventional commits format described below.

## Commit message format

```
<type>(<scope>): <description>
```

| Type | Version bump | Example |
|------|-------------|---------|
| `feat` | minor | `feat: add SFS metrics` |
| `fix` | patch | `fix: correct ECS label` |
| `feat!` | major | `feat!: rename all metrics` |
| `chore`, `docs`, `refactor`, `test` | none | `chore: update dependencies` |

**Scopes** are optional but encouraged for chart-only changes, e.g. `fix(chart): update dashboard`.

### Version bump rules

- `fix:` → patch bump (0.0.x)
- `feat:` → minor bump (0.x.0)
- `feat!:` or a commit body containing `BREAKING CHANGE:` → major bump (x.0.0)
- `chore:`, `docs:`, `refactor:`, `test:` → no bump; these accumulate without opening a Release PR

To force a patch release after a series of non-bumping commits, include at least one `fix:` PR in the batch.

### Chart-only changes

Use `fix(chart):` for changes that only affect the Helm chart (e.g. dashboard updates, default value tweaks). This triggers a full release including a Docker image rebuild. The resulting image is functionally identical to the previous one.

## Release flow

1. Merge one or more `feat:` / `fix:` PRs to `main`.
2. release-please opens (or updates) a Release PR that bumps `version.go`, `Chart.yaml` (`version` and `appVersion`), and `CHANGELOG.md`.
3. Review and merge the Release PR.
4. release-please creates a `vX.Y.Z` tag.
5. The Release workflow triggers on that tag:
   - **Guard** — fails if `Chart.yaml` `version` or `appVersion` does not match the tag. If this fires, a release-please PR was not merged first.
   - **Quality gates** — `lint-go`, `test-go`, `helm-lint`, `helm-test` run in parallel.
   - **docker** and **goreleaser** run after all quality gates pass.
   - **helm-release** runs after `docker` completes (and `goreleaser`).

## Local development

```sh
make check       # run Go lint + tests + helm lint + helm unit tests
make helm-docs   # regenerate charts/otc-prometheus-exporter/README.md
```

`make helm-docs` uses Docker; Commit the updated `README.md` alongside any `values.yaml` or `Chart.yaml` changes, or the CI `helm-docs-check` job will fail.
