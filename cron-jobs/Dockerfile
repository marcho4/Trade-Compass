# Placeholder Dockerfile for Cron Jobs
# TODO: Replace with actual cron jobs implementation

FROM alpine:latest

WORKDIR /app

# Install crond
RUN apk add --no-cache dcron

# Create a simple placeholder cron job
RUN echo '#!/bin/sh\necho "Cron job placeholder executed at $(date)"' > /app/job.sh && \
    chmod +x /app/job.sh

# Add cron job (runs every hour)
RUN echo "0 * * * * /app/job.sh >> /var/log/cron.log 2>&1" > /etc/crontabs/root

CMD ["crond", "-f", "-l", "2"]
