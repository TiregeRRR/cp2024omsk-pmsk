FROM library/node:16-bookworm as builder



ARG DEBIAN_FRONTEND=noninteractive
RUN apt update \
    && apt install -y --no-install-recommends sudo \
    && apt autoremove -y \
    && rm -rf /var/lib/apt/lists/* \
    && echo "node ALL=(ALL) NOPASSWD: ALL" >/etc/sudoers.d/node \
    && chmod 0440 /etc/sudoers.d/node

WORKDIR /app
COPY . .
RUN npm i

ARG BFF_HOST
ENV VITE_BFF_HOST=$BFF_HOST

RUN npm run build

#CMD ["bash"]
FROM python:3.12.5-alpine3.20

WORKDIR /app
COPY --from=builder /app/dist .

CMD ["python3.12", "-m", "http.server"]