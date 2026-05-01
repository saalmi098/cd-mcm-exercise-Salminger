
## Task 3: Dockerfile Analysis

### Multi-Stage Build Stages

#### Stage 1: `builder` (`golang:1.26-alpine`)

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-server ./cmd/api
```

| Step | Purpose |
|------|---------|
| `FROM golang:1.26-alpine` | Full Go toolchain (~300 MB) - only needed to compile |
| `COPY go.mod go.sum` then `go mod download` | Separate layer for deps - Docker cache skips re-download if only source changes |
| `COPY . .` | Source copied after deps to maximize cache hits |
| `CGO_ENABLED=0 go build` | Compile statically linked binary for Linux |

#### Stage 2: runtime (`alpine:3.19`)

```dockerfile
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /api-server .
EXPOSE 8080
ENTRYPOINT ["./api-server"]
```

| Step | Purpose |
|------|---------|
| `FROM alpine:3.19` | Minimal base - ~7 MB, no Go toolchain |
| `ca-certificates` | Needed for HTTPS outbound calls (TLS root CAs) |
| `COPY --from=builder` | Pulls only compiled binary from Stage 1; entire build environment discarded |
| `EXPOSE 8080` | Documents port (does not publish - `docker run -p` still required) |

### `CGO_ENABLED=0` Explained

CGO = C-Go bridge. When enabled, Go can call C libraries via `cgo`. Disabling it forces pure static compilation - all code linked into one self-contained binary, no shared library dependencies.

**Why it matters here:**

- Alpine uses `musl libc`, not `glibc`. A CGO-enabled binary built against glibc inside the builder image would crash on Alpine at runtime.
- Static binary runs on any Linux distro (or even `FROM scratch`) with zero runtime deps.
- --> Ensures predictable behavior across environments.

### Image Size: Multi-Stage vs Single-Stage

#### How it was measured

Two Dockerfiles were built and compared:

- `Dockerfile` — the existing multi-stage build (builder + runtime stage)
- `Dockerfile.single` — a single-stage equivalent kept in one `golang:1.26-alpine` image

```bash
docker build -t product-catalog:multi .
docker build -f Dockerfile.single -t product-catalog:single .
docker images | grep product-catalog
```

#### Results

```
IMAGE                    DISK USAGE
product-catalog:multi      18.2 MB
product-catalog:single      349 MB
```

| Image | Size | Contains |
|-------|------|---------|
| Single-stage (`golang:1.26-alpine`) | **349 MB** | Go toolchain, compiler, stdlib sources, build cache, binary |
| Multi-stage (`alpine:3.19`) | **18.2 MB** | Alpine base, ca-certificates, compiled binary only |

#### Why the difference is so large

Go is a compiled language — the binary is the only artifact needed at runtime. The Go toolchain (compiler, linker, standard library sources, build cache) that produces it is never needed again once compilation is done.

A single-stage build keeps all of that in the final image:

- `golang:1.26-alpine` base alone is ~330 MB (compiler, runtime, stdlib)
- Build cache and intermediate object files add more on top

Multi-stage discards the entire builder layer. The runtime stage starts from `alpine:3.19` (~7 MB) and copies only the compiled binary (~10 MB). Nothing else survives into the final image.

**Multi-stage is 19× smaller (18.2 MB vs 349 MB).**