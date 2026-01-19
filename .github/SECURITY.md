# Security Policy

As Praetor is used almost exclusively in penetration testing environments, security is our utmost priority. We believe Praetor should only add to engagements, not be the reason of compromise.

<br/>

## Supported Versions

Due to the nature of the tool, we only actively support the latest version. You are free to report vulnerabilities in older versions, however it's a significantly lower priority for our primary maintainers and therefore may take longer to be addressed.

<br/>

## Threat Model

Praetor is not designed to work in an absolutely compromised environment. While it is built to handle compliance standards and potentially compromised networks or listeners, a clear threat model is in place to define exactly what security benefits Praetor provides and what is out of scope or too loosely related to Praetor's functionality.

**In-scope**:
* MITM Attacks when syncing events between servers
* Praetor Server compromise (NOT this repository) or malicious server operator
* Modifications of event log without explicit use of Praetor commands
* Deletion of events within event log without explicit use of Praetor commands
* Read-only access by a malicious user to a penetration tester's filesystem

**Out-of-scope**:
* R/W or command execution ability on a penetration tester's system
* Stolen encryption keys
* Pre-existing malware (e.g. keyloggers) on a penetration tester's system
* Timing analysis

<br/>

## Reporting a Vulnerability

We use GitHub's built-in private vulnerability reporting system for reporting any security flaws, which can be found [here](https://github.com/lachlanharrisdev/praetor/security/advisories/new). Please fill in all the fields to the best of your ability, however this is not a requirement for your vulnerability to be triaged.

While the system should provide a simple-to-follow template for reporting, we ask that you try to provide us with as much of the following information as possible:


- The type of vulnerability (e.g., buffer overflow, privilege escalation, IDOR etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit OR release version)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit the issue

**Please do not report any security vulnerabilities in public spaces, including but not limited to Github Discussions, Issues, Pull Requests, or any social media platforms**

> [Report a vulnerability](https://github.com/lachlanharrisdev/praetor/security/advisories/new)
