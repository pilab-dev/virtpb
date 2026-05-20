# virtpb

[![Go Reference](https://pkg.go.dev/badge/go.pilab.hu/cloud/virtpb.svg)](https://pkg.go.dev/go.pilab.hu/cloud/virtpb)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`virtpb` is the canonical Protocol Buffer schema repository and generated Go client library for **PiVirt Cloud**, a full-stack virtualization orchestration platform. It defines every API contract across the platform — from the low-level hypervisor daemon up to the CLI tool and web frontend — in a single, versioned, language-neutral schema.

This repository contains the `.proto` source files (under `proto/`) and the corresponding generated Go code (gRPC + Connect), so consumers can import it directly without needing the protobuf toolchain.

---

## Table of Contents

- [What is PiVirt Cloud?](#what-is-pivirt-cloud)
- [Architecture](#architecture)
- [Packages](#packages)
- [What Problems Does It Solve?](#what-problems-does-it-solve)
- [Getting Started](#getting-started)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

---

## What is PiVirt Cloud?

PiVirt Cloud is an open-source virtualization management platform that orchestrates the complete lifecycle of virtual infrastructure:

- **Virtual machines** — create, start, stop, migrate, snapshot, clone, import/export OVA
- **Virtual networking** — VLAN/VXLAN, IPAM, distributed virtual switches, public IPs, QoS, security groups
- **Storage** — datastore management, iSCSI targets, content library (ISOs, templates)
- **Orchestration** — distributed task queue for agent coordination, job tracking
- **User management** — RBAC with roles, role groups, and permission inheritance
- **Real-time monitoring** — WebSocket event streaming for VM events, host metrics, job progress
- **CLI & Web UI** — full-featured CLI (`vpsctl`) and web frontend APIs

This repository (`virtpb`) is the **single source of truth** for all service interfaces and message types that make these capabilities possible.

---

## Architecture

The schema is organized into four architectural layers with two cross-cutting packages:

```
                   ┌──────────────────────────────────────────────┐
                   │          Frontend / API Layer                │
                   │  frontend/v1  │  cli/v1  │  ws/v1            │
                   │  (Web UI APIs)│ (CLI API)│ (WebSocket events)│
                   └──────────────────────┬───────────────────────┘
                                          │
                   ┌──────────────────────▼───────────────────────┐
                   │          Director / Orchestration Layer      │
                   │  director/v1  │  director/v2                 │
                   │  (deprecated) │  (task queue orchestration)  │
                   └──────────────────────┬───────────────────────┘
                                          │
                   ┌──────────────────────▼───────────────────────┐
                   │          Agent Layer                         │
                   │  agent/v1  │  agent/v2                       │
                   │  (VM, storage,   │  (content library, iSCSI, │
                   │   networking,    │   enhanced VM config,      │
                   │   snapshots,     │   host management)         │
                   │   DCUI)          │                            │
                   └──────────────────────┬───────────────────────┘
                                          │
                   ┌──────────────────────▼───────────────────────┐
                   │          Hypervisor Layer                    │
                   │  pivirtd/v1                                  │
                   │  (VM lifecycle, queries, metrics)            │
                   └──────────────────────────────────────────────┘

        ┌────────────────────────────────────────────────────────────────┐
        │  Cross-cutting packages                                        │
        │  common/v1 — Shared types: VM, Host, Task, Stream, Error       │
        │  network/v2 — Virtual networking: VLAN/VXLAN, IPAM, DVS, QoS   │
        └────────────────────────────────────────────────────────────────┘
```

### Layer Descriptions

| Layer | Purpose |
|-------|---------|
| **Hypervisor** (`pivirtd/v1`) | Low-level daemon running on each host. Handles VM lifecycle (create, start, stop, pause, reboot, delete) and exposes VM queries and metrics. Direct replacement for libvirtd. |
| **Agent** (`agent/v1`, `agent/v2`) | Host-level service that mediates between the director and the hypervisor. V1 covers VM operations, storage, networking, snapshots, and DCUI. V2 adds content library (ISO/template management), OVA export/import, iSCSI targets, enhanced VM configuration (CPU topology, NUMA pinning), and host maintenance. |
| **Director** (`director/v1`, `director/v2`) | Central orchestration layer. V1 (deprecated) used bidirectional streaming for agent-director communication. V2 introduces a distributed task queue: agents subscribe to tasks, execute them, and report results — covering VM lifecycle, migrations, snapshots, cloning, and OVA operations. |
| **Frontend** (`frontend/v1`, `cli/v1`, `ws/v1`) | Multi-channel API surface. `frontend/v1` provides management APIs (users, roles, RBAC, host monitoring, task listings). `cli/v1` consolidates 12 service groups into a single API for the `vpsctl` CLI tool. `ws/v1` streams real-time events (job state, VM events, host stats) over WebSocket. |
| **Network** (`network/v2`) | Full virtual networking stack: network/port CRUD, VLAN range management, VXLAN virtual clouds, VNI endpoints, IPAM (IP address management), distributed virtual switches, public IP management, QoS policies, security groups, cluster bridging, traffic flow statistics. |
| **Common** (`common/v1`) | Shared message types reused across all layers: VM definitions, disks, network interfaces, snapshots, host info, task descriptors, stream chunks, error details, and VM state enums. |

---

## Packages

| Go Package | Services / Messages |
|------------|-------------------|
| `pilab/pivirtd/v1` | `PivirtdService` — CreateVM, StartVM, StopVM, PauseVM, ResumeVM, RebootVM, DeleteVM, queries, metrics |
| `pilab/agent/v1` | `AgentService`, `VmService`, `StorageService`, `SnapshotService`, `NetworkService`, `DcuiService`, `AgentDirector` |
| `pilab/agent/v2` | `AgentService` (content library, OVA, datastores, enhanced VM config, host config), `IscsiService` |
| `pilab/director/v1` | `CloudDirector` (deprecated bidirectional streaming), `JobService` (job progress tracking) |
| `pilab/director/v2` | `TaskService` — distributed task queue CRUD and streaming subscription |
| `pilab/frontend/v1` | `ManagementService` (users, roles, RBAC), `HostService` (host monitoring), `TasksService` (task listing) |
| `pilab/cli/v1` | `AuthService`, `VMService`, `DiskService`, `ImageService`, `SnapshotService`, `NetworkService`, `PublicIPService`, `SecurityService`, `TemplateService`, `MonitoringService`, `AdminService`, `SystemService` |
| `pilab/ws/v1` | `WsEventService` — real-time event streaming over WebSocket |
| `pilab/network/v2` | `NetworkService`, `IpamService`, `DvsService`, `PublicIpService`, `QosSecurityService`, `ClusterBridgeService`, `FlowService`, `StatisticsService` |
| `pilab/common/v1` | Shared messages: `VM`, `Disk`, `NetworkInterface`, `Snapshot`, `Host`, `Task`, `StreamMessage`, `Error` |

---

## What Problems Does It Solve?

### 1. Contract Standardisation
Every service boundary in the platform is defined in a single, versioned protobuf schema. Frontend, CLI, agent, and director all speak the same language — no drift, no duplication, no ad-hoc API clients.

### 2. Multi-Protocol Support
The same `.proto` definitions generate both **gRPC** stubs and **Connect-Go** clients and servers. Consumers can use HTTP/2 gRPC, HTTP/1.1+Connect, or both — the generated code supports all combinations.

### 3. Versioning & Compatibility
Breaking changes are gated by the package version suffix (`v1`, `v2`). The deprecated `director/v1` and modern `director/v2` coexist, allowing gradual migration. Buf's breaking change detection prevents accidental regressions.

### 4. Single Import, Strong Typing
Consumers import one Go module (`go.pilab.hu/cloud/virtpb`) and get fully typed request/response messages and service interfaces. No manual marshalling, no runtime reflection, no guessing field names.

### 5. Schema-Driven Development
The `.proto` files are the source of truth. Code generation produces Go types, gRPC registrations, Connect handlers, and OpenTelemetry metadata — all from a single schema definition.

### 6. Ecosystem Compatibility
Because the schema is pure protobuf, it can generate clients in other languages (TypeScript, Python, etc.) for the same APIs. The `buf.gen.yaml` already includes a commented-out TypeScript target for future use.

---

## Getting Started

### Prerequisites

- Go 1.25.0 or later

### Installation

```bash
go get go.pilab.hu/cloud/virtpb
```

### Usage

Import the root package or individual sub-packages:

```go
import (
    // Import specific packages by version
    pivirtdv1 "go.pilab.hu/cloud/virtpb/pilab/pivirtd/v1"
    networkv2 "go.pilab.hu/cloud/virtpb/pilab/network/v2"

    // Connect protocol clients
    "go.pilab.hu/cloud/virtpb/pilab/pivirtd/v1/pivirtdv1connect"
)
```

The generated types are fully compatible with both gRPC and Connect runtimes. See the [GoDoc](https://pkg.go.dev/go.pilab.hu/cloud/virtpb) for the full API reference.

---

## Development

### Repository Structure

```
.
├── proto/pilab/          # Source .proto files
│   ├── agent/v1/         # Agent v1 services
│   ├── agent/v2/         # Agent v2 services
│   ├── cli/v1/           # CLI API
│   ├── common/v1/        # Shared messages
│   ├── director/v1/      # Director v1 (deprecated)
│   ├── director/v2/      # Director v2 (task queue)
│   ├── frontend/v1/      # Web frontend API
│   ├── network/v2/       # Virtual networking
│   ├── pivirtd/v1/       # Hypervisor daemon
│   └── ws/v1/            # WebSocket events
├── pilab/                # Generated Go code
├── buf.gen.yaml          # Buf code generation config
├── buf.yaml              # Buf workspace config
├── go.mod                # Go module definition
└── doc.go                # Package documentation
```

### Regenerating Go Code

The project uses [buf](https://buf.build) for code generation. To regenerate from the `.proto` sources:

```bash
# Install buf (if not already installed)
# See https://buf.build/docs/installation

# Regenerate all Go code
buf generate
```

### Versioning

This module follows Go module versioning conventions. The current version is `v1.0.4` (see `go.mod`). Breaking schema changes are introduced in new API versions (e.g., `v1` → `v2`) while backward compatibility is maintained within a version.

---

## Contributing

Contributions are welcome and appreciated! Here's how you can help:

- **Report bugs** — Open an issue describing the problem, expected behavior, and steps to reproduce.
- **Suggest features** — Open an issue with the "enhancement" label.
- **Submit pull requests** — Fork the repository, make your changes, and open a PR.

### Guidelines

1. Ensure `.proto` changes are backward-compatible within the existing version package, or use a new version suffix for breaking changes.
2. Run `buf lint` and `buf breaking` before committing schema changes.
3. Regenerate Go code with `buf generate` and commit the results.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details. 