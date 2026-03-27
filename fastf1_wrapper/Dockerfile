FROM python:3.14-slim

# Install curl for healthcheck
RUN apt-get update && apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root user with UID 1000 for volume compatibility
RUN groupadd -g 1000 appgroup && \
    useradd -u 1000 -g appgroup -m -s /bin/sh appuser

WORKDIR /app

# Copy configuration files
COPY pyproject.toml ./

# Install dependencies without cache and not editable
RUN pip install . --no-cache-dir

# Copy source and set ownership
COPY src/ ./src/
RUN mkdir -p /app/.cache && chown -R appuser:appgroup /app && chmod 777 /app/.cache

EXPOSE 8081

# Healthcheck to probe liveness
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl --fail http://localhost:8081/health || exit 1

USER appuser

CMD ["python", "-m", "src.core.http.server"]
