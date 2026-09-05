# Security Policy

The Charites team takes the security of our static analysis compiler and AST linters seriously. We appreciate responsible disclosure of security vulnerabilities.

---

## Supported Versions

Only the latest minor release branch receives active security updates and patches.

| Version | Supported          |
| :------ | :----------------- |
| `v1.x`  | :white_check_mark: |
| `< 1.0` | :x:                |

---

## Tool Description Policy (MCP)

Charites Model Context Protocol (MCP) tools follow a **descriptive, not imperative** policy for tool descriptions to prevent "tool steering" or "tool description poisoning" - an adversarial pattern where tool metadata is used to imperatively command AI agents rather than describe functionality.

### Design Principles:
- **Avoid Imperative Language:** Descriptions specify what a tool does, not when/how it MUST be invoked.
- **No Agent Steering:** We never inject hidden prompts to force AI agents to call specific tools.
- **User Agency:** Agents and developers make informed decisions about tool invocation based on user intent.

---

## Safe Traversal & Containment

Charites inspects user code and templates safely:
- **Path Traversal Defense:** All file system inputs are normalized and checked to prevent arbitrary filesystem traversal outside the intended scan root.
- **Bounded AST Traversal:** Node traversal uses iterative range-over-func loops rather than unbounded recursive call stacks, defending against stack overflow attacks from deeply nested markup.
- **Fail-Closed Directives:** Malformed suppression comments do not silently disable rules.

---

## Reporting a Vulnerability

**Please do NOT report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in Charites (such as an arbitrary code execution vector, malicious AST exploitation, parser denial-of-service, or rule bypass), please report it responsibly:

### Option 1: GitHub Private Vulnerability Reporting (Preferred)
Submit a report through GitHub's private vulnerability advisory dashboard:
 **[Report a Security Vulnerability](https://github.com/will2469/charites/security/advisories/new)**

### Option 2: Direct Email
Send an email with the subject line `[SECURITY] Vulnerability in Charites` to:
 **[will.i.is.ega@gmail.com](mailto:will.i.is.ega@gmail.com)**

Please include:
- A detailed description of the vulnerability.
- A minimal reproducible web template (`.astro`, `.tsx`, or `.css`) demonstrating the issue.
- Any potential impact or attack vectors.
- (Optional) Proposed remediation or patch.

---

## Response Timeline & Process

- **Acknowledgment:** We aim to acknowledge receipt of your vulnerability report within **48 hours**.
- **Assessment:** We will validate the issue, assess severity, and provide a remediation timeline within **7 business days**.
- **Coordinated Disclosure:** A security patch will be prepared in a private fork and released alongside an official advisory.
- **Credit:** We will publicly credit your contribution in the release notes and security advisory (unless you prefer to remain anonymous).
