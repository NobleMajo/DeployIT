### BASE
FROM golang AS base

WORKDIR /app

COPY go.mod* go.sum* ./
RUN go mod download

### LOCAL
FROM base AS local

RUN apt-get update \
	&& apt-get install -y --no-install-recommends make git \
	&& apt-get clean \
    && apt-get autoclean \
    && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["bash"]

### BASE DEPLOY
FROM base AS base-deploy
COPY . .
RUN make build

### DEPLOY
FROM ubuntu:24.04 AS deploy

RUN useradd -m appuser --uid 10000
USER 10000

COPY --from=base-deploy --chown=10000 /app/bin /usr/local/bin/appbin

CMD ["appbin"]
