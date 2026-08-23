FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN go build -o /out/chronos-mesh ./cmd/chronos-mesh
FROM gcr.io/distroless/base-debian12
COPY --from=build /out/chronos-mesh /chronos-mesh
COPY configs /configs
EXPOSE 8080 8123/udp 8443
ENTRYPOINT ["/chronos-mesh"]
