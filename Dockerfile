ARG GO_VERSION=1.21

FROM golang:${GO_VERSION} AS base
LABEL authors="pavel.karpovich"
ENV CGO_ENABLED=0

FROM base AS build
WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd ./app && go build -o /bin/bot .


FROM alpine:3.19 AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=build /bin/bot /bin/


ENTRYPOINT [ "/bin/bot" ]