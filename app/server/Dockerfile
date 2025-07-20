FROM golang:1.24.2 AS builder

ARG GIT_USERNAME
ARG GIT_EMAIL

RUN --mount=type=secret,id=git_token \
    git_token=$(cat /run/secrets/git_token) && \
    git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    git config --global url."https://${git_token}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hydraide .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /hydraide/

COPY --from=builder /app/hydraide .

COPY certificate/server.crt /hydraide/certificate/server.crt
COPY certificate/server.key /hydraide/certificate/server.key

RUN apk --no-cache add ca-certificates

CMD ["./hydraide"]

HEALTHCHECK --interval=10s --timeout=3s --start-period=3s --retries=3 \
  CMD curl --fail http://localhost:4445/health || exit 1