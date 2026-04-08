---
id: testing
version: 1.0.0
description: Software testing best practices for AI agents
---

# Testing Best Practices

When writing tests, follow these guidelines:

- Write tests before or alongside the production code (TDD or BDD encouraged).
- Name tests clearly to describe the behaviour being verified, not the implementation.
- Keep tests fast, isolated, and repeatable; avoid shared mutable state.
- Use the Arrange-Act-Assert (AAA) pattern to structure each test case.
- Mock external dependencies (databases, network calls, filesystems) in unit tests.
- Write integration tests that exercise real infrastructure separately from unit tests.
- Aim for high coverage of critical paths and edge cases, not just line coverage.
- Never skip or comment out failing tests; fix them or delete them.
- Use test fixtures and factories to reduce repetitive test setup code.
- Assert on outcomes and behaviour, not on internal implementation details.
