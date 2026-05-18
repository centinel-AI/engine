# Security Policy

## Supported versions

Only the latest version of Grauss on the `main` branch receives security fixes.

## Reporting a vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

If you discover a security issue — including problems in generated Terraform configurations, insecure defaults, or credential exposure risks — please report it privately:

1. Email the maintainers at the address listed in the repository profile, or
2. Use [GitHub's private vulnerability reporting](https://docs.github.com/en/code-security/security-advisories/guidance-on-reporting-and-writing/privately-reporting-a-security-vulnerability) if enabled for this repository.

Include as much detail as possible:

- A description of the vulnerability and its potential impact
- Steps to reproduce or a minimal proof of concept
- The Terraform and provider versions in use
- Any relevant configuration or JSON snippets (redact real credentials)

## Response timeline

| Stage | Target |
|-------|--------|
| Acknowledgement | Within 3 business days |
| Initial assessment | Within 7 business days |
| Fix or mitigation | Depends on severity — critical issues are prioritised |

## Scope

This policy covers:

- Terraform configurations under `providers/`
- Data file handling (JSON parsing, file path traversal)
- Documentation that could lead users to adopt insecure practices

Out of scope:

- Vulnerabilities in Terraform or OpenTofu themselves or in provider plugins (incl. `azurerm`, `azuread`, `aws`, `google`, `oci`) — report those to the respective **engine** or **provider** maintainers
- Issues in cloud provider APIs or services (Microsoft Azure, AWS, Google Cloud, Oracle Cloud Infrastructure, etc.) — report those to the respective vendor
- Findings from automated scanners without a demonstrated impact

## Disclosure

Once a fix is available, we will publish a security advisory describing the vulnerability, its impact, and the remediation steps. Credit will be given to the reporter unless they prefer to remain anonymous.
