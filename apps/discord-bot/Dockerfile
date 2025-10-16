FROM python:3.12-slim

# Install prerequisites for uv
RUN apt-get update && apt-get install -y curl git

# Install uv by Astral
RUN curl -LsSf https://astral.sh/uv/install.sh | sh

# Set working directory
WORKDIR /app

# Copy bot files
COPY . .

# Environment variables are handled by docker-compose
ENV PATH="/root/.local/bin:$PATH"

# Default command to run bot
CMD ["uv", "run", "main.py"]