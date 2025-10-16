FROM --platform=linux/amd64 debian:stable-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libsqlite3-0 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy binary and templates
ADD my-journal /usr/bin/my-journal
COPY template ./template

# Create data directory for database
RUN mkdir -p /app/data

EXPOSE 8080

CMD ["my-journal"]