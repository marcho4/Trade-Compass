# Stage 1: Builder
FROM python:3.12-slim AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .

RUN pip install --no-cache-dir --prefix=/install -r requirements.txt


# Stage 2: Runtime
FROM python:3.12-slim AS runtime

# Install runtime dependencies (Chrome for Selenium)
RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq5 \
    chromium \
    chromium-driver \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean

# Set Chrome environment variables
ENV CHROME_BIN=/usr/bin/chromium
ENV CHROMEDRIVER_PATH=/usr/bin/chromedriver

WORKDIR /app

# Copy installed packages from builder
COPY --from=builder /install /usr/local

# Create non-root user
RUN groupadd --gid 1000 appgroup \
    && useradd --uid 1000 --gid appgroup --shell /bin/bash --create-home appuser

# Copy application code
COPY --chown=appuser:appgroup . .

# Switch to non-root user
USER appuser

EXPOSE 8000

CMD ["uvicorn", "main:api", "--host", "0.0.0.0", "--port", "8081"]
