FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.work /go.work.sum ./ 
COPY /shared/ ./shared/
COPY /user-svc ./user-svc
COPY /order-svc ./order-svc

WORKDIR /app/user-svc

# Mount cache for go dependencies and artifacts from builds.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build 

RUN go mod download 

RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags "-s -w" -o /bin/usersvc

FROM gcr.io/distroless/static-debian12:nonroot

ENV LOG_LEVEL=debug

WORKDIR /app

COPY --from=build-stage bin/usersvc ./

# Doesn't really do anything. Only used for documentation for end user.
# Container defaults to 5050
EXPOSE 5050

CMD ["/app/usersvc"]