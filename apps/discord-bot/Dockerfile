FROM python:3.12-slim

RUN apt-get update && apt-get install -y curl git

RUN curl -LsSf https://astral.sh/uv/install.sh | sh

WORKDIR /app

COPY . .

ENV PATH="/root/.local/bin:$PATH"

CMD ["uv", "sync"]

CMD ["uv", "run", "main.py"]
CMD ["uv", "run", "api_server.py"]
