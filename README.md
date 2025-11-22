# 📌 `imgpin` — Deterministic Image Digest Pinning for Containers & Dev Environments

[![CI](https://github.com/jloper3/imgpin/actions/workflows/ci.yaml/badge.svg)](https://github.com/jloper3/imgpin/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jloper3/imgpin)](https://goreportcard.com/report/github.com/jloper3/imgpin)

`imgpin` is a fast, lightweight Go CLI that resolves container image tags to immutable `@sha256` digests and rewrites your configuration files to ensure **fully reproducible environments**.

It supports:

* **Dockerfiles**
* **DevContainer definitions (`devcontainer.json`)**
* **Kubernetes manifests**
* **Lockfiles** to freeze image states (`imgpin.lock`)
* **Fast registry caching**
* **Unified diff output**

This prevents silent supply-chain drift, increases build determinism, and enforces EnvSecOps-style “**check the bag**” environment integrity.

---

## 🚀 Features

* 🔐 **Resolve image tags → digests**
* 🛠 **Rewrite files** in place (or print to stdout)
* 📦 **Lockfile support** (`imgpin lock`)
* 🔍 **Drift detection** (`imgpin check`)
* ⚡ **Persistent cache** (`~/.cache/imgpin/digests.json`)
* 📄 **Unified diffs** when changes occur
* 🧩 **Handlers** for Docker, DevContainers, and Kubernetes YAML
* 💨 **Zero dependencies beyond core Go & official containerregistry**

---

## 📦 Installation

### Build from source

```bash
make build
```

The binary will be available at:

```
bin/imgpin
```

### Install into GOPATH

```bash
make install
```

---

## 🧭 Usage

### **Resolve an image to a digest**

```bash
imgpin resolve python:3.11
```

Output:

```
python@sha256:<digest>
```

---

### **Rewrite a Dockerfile**

Dry-run:

```bash
imgpin file Dockerfile
```

In-place modification:

```bash
imgpin file Dockerfile --in-place
```

---

### **Rewrite a DevContainer file**

```bash
imgpin file .devcontainer/devcontainer.json --in-place
```

---

### **Rewrite a Kubernetes manifest**

```bash
imgpin file deploy.yaml --in-place
```

---

### **Generate / update a lockfile**

```bash
imgpin lock
```

Produces or updates:

```
imgpin.lock
```

This captures **all images referenced across your repo** along with their resolved digests.

---

### **Check your repository for drift**

```bash
imgpin check
```

If upstream tags changed or files added new images, `imgpin check` reports:

* Which images drifted
* Exact diffs
* Missing lockfile entries
* Unexpected images

This is ideal for pre-deploy checks or CI pipelines.

---

## 📂 Project Structure

```
.
├── cmd/imgpin/              # CLI entrypoint
├── internal/
│   ├── cli/                 # Commands: resolve, file, lock, check
│   ├── cache/               # JSON/TTL cache
│   ├── resolve/             # Registry resolution + caching
│   ├── diff/                # Unified diff utilities
│   ├── lockfile/            # imgpin.lock read/write
│   └── handlers/
│       ├── dockerfile/
│       ├── devcontainer/
│       └── kubernetes/
└── Makefile
```

---

## 🛠 Development

### Format, vet, tidy

```bash
make fmt
make vet
make tidy
```

### Run tests

```bash
make test
```

HTML coverage report:

```bash
make cover-html
```

---

## 🔒 Philosophy

Container tags drift silently. Dev environments diverge undetected.
With `imgpin`, you can **pin every environment artifact deterministically**, verify them across your repository, and detect drift before it becomes a supply-chain issue.

This is one of the simplest—but highest-leverage—steps toward **attestation-ready environments** and a practical EnvSecOps workflow.

---

## 📜 License

MIT License
