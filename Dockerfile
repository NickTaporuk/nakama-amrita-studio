# Build stage
FROM golang:1.14.3-buster as builder

# Force the go compiler to use modules
ENV GO111MODULE=on

# Set working directory to current directory
WORKDIR /nakama-amrita-studio

# Copy all the files in the project
COPY . .
RUN ls -la

# Download all the dependencies that are specified in the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=plugin  -trimpath /nakama-amrita-studio/user.go
RUN ls -la /nakama-amrita-studio

# Final Stage
FROM heroiclabs/nakama:2.12.0
#
COPY --from=builder /nakama-amrita-studio/user.so /nakama/data/modules/user.so
RUN ls -la /nakama/data/modules
EXPOSE 7349 7350 7351

ENTRYPOINT ["tini", "--", "/nakama/nakama"]

HEALTHCHECK --interval=5m --timeout=10s \
  CMD curl -f http://localhost:7350/ || exit 1