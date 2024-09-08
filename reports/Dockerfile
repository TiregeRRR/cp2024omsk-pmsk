FROM ubuntu:22.04 as build

ENV PYTHON_VERSION=3.11

RUN export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y update \
    && apt-get -y install --no-install-recommends \
    python${PYTHON_VERSION} \
    python3-pip python3.11-venv \
    git \
    abiword \
    && rm -rf /var/lib/apt/lists/*

RUN ln -s -f /usr/bin/python${PYTHON_VERSION} /usr/bin/python3 && \
    ln -s -f /usr/bin/python${PYTHON_VERSION} /usr/bin/python && \
    ln -s -f /usr/bin/pip3 /usr/bin/pip && \
    python3 -m venv /venv && /venv/bin/pip install --timeout=1200 --no-cache-dir --upgrade pip setuptools wheel

FROM build AS build-venv

COPY requirements.txt /requirements.txt
RUN /venv/bin/pip install --timeout=12000 --no-cache-dir --disable-pip-version-check -r /requirements.txt

FROM build

COPY --from=build-venv /venv /venv

WORKDIR /app

COPY . .

EXPOSE 8000
ENTRYPOINT ["/venv/bin/python3", "-m", "granian", "--interface", "asgi", "app.main:app", "--port", "8000", "--host", "0.0.0.0"]
# For local test
### granian --interface asgi app.main:app --port 8000 --host 0.0.0.0
