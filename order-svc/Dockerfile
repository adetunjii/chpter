FROM golang:1.23 AS build-stage

WORKDIR /order-svc

COPY go.mod go.sum ./

# Mount cache for go dependencies and artifacts from builds.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build 

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags "-s -w" -o /order-svc/ordersvc

FROM gcr.io/distroless/static-debian12:nonroot

ENV LOG_LEVEL=debug

WORKDIR /app

COPY --from=build-stage /order-svc/ordersvc ./

# Doesn't really do anything. Only used for documentation for end user.
# Container defaults to 5051 
EXPOSE 5051

CMD ["/app/ordersvc"]