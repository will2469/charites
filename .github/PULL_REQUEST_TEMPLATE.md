## Description

Please provide a clear and concise summary of the changes in this pull request and the rationale behind them.

Closes #(issue)

---

## Type of Change

- [ ] `feat`: New analyzer rule, parser capability, or CLI feature
- [ ] `fix`: Bug fix, false-positive reduction, or line offset correction
- [ ] `docs`: Documentation, README, or 8-Pillars Wiki update
- [ ] `perf`: Performance or memory allocation optimization
- [ ] `refactor`: Internal refactoring with no user-facing behavior changes
- [ ] `test`: Additional test fixtures or Tri-Corpus harness improvements
- [ ] `ci`: Build, release, or GitHub Actions workflow updates

---

## Quality & Architecture Checklist

- [ ] **Semgrep Canonical ID:** Rule identifier follows `<category>.<slug>` (e.g. `theme.hardcode-color`, `a11y.alt-text`).
- [ ] **1-SSOT Tri-Corpus Tests:** Positive (`want`), Negative (`charites:ignore`), and Adversarial fixtures are provided under `tests/correctness/<category>.<slug>/`.
- [ ] **Directive Support:** Rules properly respect `charites:ignore <category>.<slug>` comments.
- [ ] **8-Pillars Documentation:** If adding a new rule, `wiki/<category>.<slug>.md` has been created adhering to the 8-pillars documentation standard.
- [ ] **All Tests Pass:** `make test-race` (`go test -race ./...`) passes with zero failures.
- [ ] **Code Hygiene Clean:** `make lint` (`gofmt` and `go vet`) passes cleanly.
- [ ] **Clean Conventional Commit:** Commit messages adhere to Conventional Commits specification.
