# Improvement Tasks Checklist

Note: Each item is prefixed with [ ] for tracking. Complete in order where possible. Some tasks can be parallelized but dependencies are noted.

1. Architecture and Boundaries
   1. [ ] Define and document high-level architecture (backend services, plugin system, data layer, web UI embedding). Include a simple diagram in docs/. 
   2. [ ] Review and refine module/package boundaries; ensure domain packages don’t depend on delivery (HTTP/UI) packages. 
   3. [ ] Tighten arch-go.yml rules to reflect intended dependencies (e.g., domain can depend on data interfaces, not concrete providers). 
   4. [ ] Introduce an internal/config package for cohesive configuration management (env vars, defaults, validation). 
   5. [ ] Replace hard-coded AWS region in pkg/data/dyndb with config-driven value and validation. 
   6. [ ] Establish a clear SPI boundary for data stores (interfaces in domain, implementations in separate subpackages). 
   7. [ ] Evaluate plugin initialization order and lifecycle; specify lifecycle hooks (init/start/stop) and ordering guarantees.

2. HTTP Server and Middleware
   8. [ ] Add health (GET /healthz) and readiness (GET /readyz) endpoints with dependency checks (e.g., DynamoDB if configured). 
   9. [ ] Implement structured request logging middleware (trace ID, method, path, latency, status) and propagate IDs via context. 
   10. [ ] Add recovery middleware to catch panics and return 500 with safe logging. 
   11. [ ] Add CORS middleware with configurable origins for web UI/local dev. 
   12. [ ] Add security headers middleware (Content-Security-Policy, X-Content-Type-Options, X-Frame-Options, Referrer-Policy). 
   13. [ ] Formalize API routing: define explicit route registration types (method, pattern) and centralize versioning strategy.

3. Authentication and Authorization
   14. [ ] Finalize authn cookie handling: set HttpOnly, SameSite, Secure (configurable for dev), proper Max-Age/Expires. 
   15. [ ] Replace placeholder cookie value (user ID only) with signed/encoded token (e.g., JWT or HMAC) and rotation strategy. 
   16. [ ] Implement logout endpoint to clear auth cookie. 
   17. [ ] Complete basic authorization layer (roles/permissions) and unskip tests blocked by “autz not implemented yet”. 
   18. [ ] Centralize authn/authz errors; standardize 401/403 responses with JSON body.

4. Context, Errors, and Observability
   19. [ ] Replace all context.TODO() usages with passed-through request or init contexts as appropriate. 
   20. [ ] Standardize error wrapping with fmt.Errorf("...: %w", err) and propagate context. 
   21. [ ] Define and use a common error response shape for APIs. 
   22. [ ] Introduce metrics (Prometheus or expvar) for requests, DB ops, and errors; add /metrics endpoint (config-gated). 
   23. [ ] Add structured logs across critical paths (auth, data ops, startup/shutdown) with stable keys.

5. Generic SPI and Resource Lifecycle
   24. [ ] Implement Close() handling in pkg/util/generic_spi.go for previously active instances during Reset() and on multiple-success cases. 
   25. [ ] Define a lightweight closer interface and unit-test SPI lifecycle (init success/failures, concurrent calls, reset). 
   26. [ ] Expose configuration for SPI init timeout and context.

6. Data Layer (DynamoDB)
   27. [ ] Make DynamoDB table prefix and region configurable with explicit validation and helpful startup logs. 
   28. [ ] Improve MigrateTable to wait for ACTIVE state on creation and to provide more actionable mismatch messages. 
   29. [ ] Replace context.TODO() with function-scoped contexts passed down from callers. 
   30. [ ] Add retry/backoff and specific error handling for throughput exceptions. 
   31. [ ] Add repository-level interfaces in domain and ensure tests can run with in-memory or mock store. 
   32. [ ] Create localdev option (e.g., DynamoDB Local) and document usage; add CI job gate to skip cloud-dependent tests unless creds provided.

7. API Design and Versioning
   33. [ ] Write an API contract (OpenAPI or markdown) for existing endpoints (/api/v1/version, auth endpoints). 
   34. [ ] Introduce typed request/response DTOs and validation. 
   35. [ ] Ensure consistent status codes and error envelopes across endpoints. 
   36. [ ] Add deprecation policy and routing for future /api/v2.

8. Web UI Integration
   37. [ ] Document the web UI build pipeline and embed strategy (go:embed, build tags). 
   38. [ ] Add development CORS config and local proxy for UI to backend. 
   39. [ ] Provide a CLI command to export UI assets (zip/tgz) with flags (output path, format); ensure errors bubble properly. 
   40. [ ] Add E2E smoke test (Cypress/Playwright) that runs against a local server in CI.

9. CLI and Commands
   41. [ ] Add consistent logging and error handling to Cobra commands (no panic; return errors). 
   42. [ ] Add --config, --log-level, --addr flags to server command; print effective config at startup (without secrets). 
   43. [ ] Add graceful shutdown signals (SIGINT/SIGTERM) wiring and longer shutdown timeout with draining.

10. Security and Secrets
   44. [ ] Move secrets management to environment variables or secret stores; never print secret values in logs. 
   45. [ ] Add static analysis checks for hard-coded credentials/regions. 
   46. [ ] Add lint rule and CI check to prevent use of http.DefaultClient/server defaults without timeouts.

11. Testing Strategy
   47. [ ] Increase unit test coverage for server handlers, middleware, and auth flows. 
   48. [ ] Add integration tests for DynamoDB (tagged/conditional). 
   49. [ ] Mock AWS in tests where possible; add golden tests for API JSON outputs. 
   50. [ ] Stabilize flaky tests by removing real network dependencies in unit tests.

12. Tooling, Quality Gates, and CI/CD
   51. [ ] Introduce golangci-lint with a curated ruleset; add to CI. 
   52. [ ] Enforce gofmt/goimports, govulncheck, and staticcheck in CI. 
   53. [ ] Add makefile or taskfile for common tasks (build, test, lint, run). 
   54. [ ] Cache Go build/test in CI using actions/setup-go cache; ensure proper restore keys. 
   55. [ ] Separate deploy steps from build/test; use environments and approvals for AWS Lambda updates. 
   56. [ ] Generate and publish coverage badge/artifacts; make coverage threshold configurable.

13. Documentation
   57. [ ] Expand README with setup, build, run, test instructions (backend and web UI). 
   58. [ ] Add docs/configuration.md listing all env vars, flags, and defaults. 
   59. [ ] Add docs/architecture.md describing packages, plugins, and data layer strategies. 
   60. [ ] Document local development workflows (running server + _webui dev server, debugging tips).

14. Performance and Reliability
   61. [ ] Add timeouts to all outbound calls and HTTP server (Read/Write/Idle timeouts). 
   62. [ ] Add connection reuse and client configuration where applicable (AWS SDK, HTTP clients). 
   63. [ ] Add basic load test scenario and track latency/throughput trends over time. 
   64. [ ] Ensure server Start/Stop returns deterministic errors and logs (avoid goroutine leaks).

15. Code Hygiene and Cleanup
   65. [ ] Remove dead code and unused exports; minimize package public surface. 
   66. [ ] Unify naming conventions (files, types, functions); avoid abbreviations except common ones (ctx, id, URL). 
   67. [ ] Add package-level doc.go files for major packages. 
   68. [ ] Replace fmt.Println or panic with slog and proper error returns; only main() may exit with non-zero.

16. Release and Distribution
   69. [ ] Add versioning strategy (semver) and changelog; make build embed version and git commit consistently. 
   70. [ ] Provide reproducible builds (CGO flags, -trimpath) and add size analysis as optional CI step. 
   71. [ ] Create Dockerfile and local compose (optional) for running the server with or without DynamoDB Local; document usage.
