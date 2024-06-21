FROM golang:1.22 as builder

WORKDIR /workspace


# Copy the Go Modules manifests
COPY go.mod go.mod
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .
# -w: This flag is used to disable DWARF generation. DWARF (Debugging With Attributed Record Formats) is a debugging format that contains information
# for debugging the binary.
# -s: This flag is used to strip the symbol table. The symbol table contains information about symbols (functions, variables, etc.) in the binary.
# Stripping the symbol table makes it harder to debug or reverse-engineer the binary but also reduces the binary's size.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o manager cmd/autoscaler.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static-debian11
WORKDIR /
COPY --from=builder /workspace/manager .
ENTRYPOINT ["/manager"]