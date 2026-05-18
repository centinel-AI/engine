## Summary

<!-- What does this PR do and why? Link the related issue if applicable (Closes #NNN). -->

## Type of change

<!-- Check all that apply. -->

- [ ] New resource type (`providers/<cloud>/`)
- [ ] Fix — template or data correction
- [ ] Codegen / `scripts/` change
- [ ] Dockerfile / Compose / Taskfile
- [ ] Documentation
- [ ] Other: <!-- describe -->

## Clouds affected

- [ ] AWS
- [ ] Azure (AzureRM + AzureAD)
- [ ] GCP
- [ ] OCI

## Checklist

- [ ] Branch follows naming: `feat/<resource>`, `fix/<description>`, or `docs/<description>`
- [ ] `terraform fmt -check` passes on changed files
- [ ] `terraform validate` passes (`task workspace:<cloud>` → `terraform validate`)
- [ ] No Terraform modules introduced — only native constructs (`for_each`, `try()`, expressions)
- [ ] `try()` used for optional JSON fields (no `lookup()` or `optional()`)
- [ ] Resource label is `grauss`; all identifiers use `snake_case`
- [ ] Each `.tf` file contains exactly one resource block; no hardcoded values
- [ ] JSON data follows `data/<cloud>/<project>/<resource_type>/<instance-name>.json` layout
- [ ] `.cursor/rules/` updated if a convention was added or changed
- [ ] `providers/` changes are generated (not hand-edited); `.gitkeep` is the only tracked file per directory
