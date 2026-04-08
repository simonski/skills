---
id: docker
version: 1.0.0
description: Docker and container best practices for AI agents
---

# Docker & Container Best Practices

When working with Docker and containers:

- Use minimal base images (distroless, alpine) to reduce attack surface and image size.
- Pin image versions with a specific digest or tag; never rely on `latest` in production.
- Run containers as a non-root user; add `USER` instruction in the Dockerfile.
- Use multi-stage builds to keep production images lean and free of build tools.
- Never store secrets in image layers; use runtime secrets or environment variables.
- Set resource limits (`--memory`, `--cpus`) for production containers.
- Use `.dockerignore` to exclude unnecessary files from the build context.
- Keep each container to a single process/responsibility.
- Use health checks (`HEALTHCHECK`) so orchestrators can detect and replace unhealthy containers.
- Store persistent data in named volumes, not in container layers.
