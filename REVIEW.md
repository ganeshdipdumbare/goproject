## Dependencies

When a change adds or upgrades a third-party package (package.json, lockfiles,
requirements.txt, pyproject.toml, go.mod), check each added or changed package
with the Socket MCP tools (`depscore`, and `alerts` where available) and report
any low score or alert as a finding.
