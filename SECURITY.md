# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please follow these steps:

### How to Report

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please email security reports to: **security@mirage.dev**

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: 1-7 days
  - High: 7-14 days
  - Medium: 14-30 days
  - Low: 30-90 days

### Disclosure Policy

- We will work with you to understand and resolve the issue
- We ask that you do not publicly disclose the vulnerability until we've had a chance to address it
- We will credit security researchers who responsibly disclose vulnerabilities
- We will publish security advisories for significant issues

### Bug Bounty

Currently, we do not have a bug bounty program, but we deeply appreciate security researchers who help keep Mirage secure.

## Security Best Practices

When using Mirage:

1. **Production Use**: Mirage is designed for development/testing - do not use in production environments
2. **Sensitive Data**: Never proxy requests containing real credentials or sensitive data
3. **Network Security**: Run Mirage only on trusted networks
4. **Access Control**: Use firewall rules to restrict access to Mirage ports

## Security Updates

Security updates will be released as patch versions and announced via:
- GitHub Security Advisories
- Release notes
- CHANGELOG.md

Subscribe to releases to stay informed about security updates.
