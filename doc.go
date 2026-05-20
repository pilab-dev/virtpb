// Package virtpb provides Protocol Buffer schema definitions and generated Go code
// for the PiVirt Cloud virtualization orchestration platform.
//
// # What is PiVirt Cloud?
//
// PiVirt Cloud is a full-stack virtualization management platform that covers the
// entire lifecycle of virtual machines, virtual networks, storage, and the
// supporting infrastructure (RBAC, monitoring, real-time events, CLI tooling).
//
// This repository (virtpb) serves as the single source of truth for all API
// contracts across the platform. Every service boundary — from the low-level
// hypervisor daemon up to the CLI and web frontend — is defined here as
// Protocol Buffer schemas and made available as idiomatic Go types via code
// generation.
//
// # Architecture
//
// The schema is organized into four architectural layers:
//
//	Layer              Go Package                          Role
//	------            ----------                          ----
//	Hypervisor         pilab/pivirtd/v1                    Low-level VM lifecycle on individual hosts
//	Agent              pilab/agent/v1, pilab/agent/v2      Host-level services (VM ops, storage, networking, iSCSI, DCUI)
//	Director           pilab/director/v1, pilab/director/v2  Distributed task orchestration & agent coordination
//	Frontend/API       pilab/frontend/v1, pilab/cli/v1, pilab/ws/v1  Web UI, CLI tool, real-time WebSocket events
//
// Two cross-cutting packages support every layer:
//   - pilab/common/v1   — Shared message types (VM, Host, Task, Stream, Error)
//   - pilab/network/v2  — Virtual networking (VLAN/VXLAN, IPAM, DVS, QoS, security, flow statistics)
//
// # What It's Good For
//
// Building Go clients or servers that interact with any part of the PiVirt Cloud
// platform:
//   - VM lifecycle management (create, start, stop, migrate, snapshot, clone, OVA import/export)
//   - Virtual networking (network CRUD, port management, IPAM, DVS, public IPs, QoS, security groups)
//   - Storage orchestration (datastores, iSCSI targets, ISO/template content library)
//   - Task queue management (distributed job distribution to agent nodes)
//   - User & role administration (RBAC, role groups, permissions)
//   - Real-time event streaming (VM events, host metrics, job progress via WebSocket)
//   - CLI integration (12 service groups covering auth, VMs, disks, images, snapshots, networks, monitoring, system)
//
// # Problems It Solves
//
//   - Contract standardisation: Every service boundary is defined in a single,
//     versioned, language-neutral schema — no drift between frontend, CLI,
//     and backend implementations.
//   - Multi-protocol support: The same .proto definitions are used to generate
//     both gRPC stubs and Connect-Go clients/servers, giving consumers the
//     choice of HTTP/2 gRPC or HTTP/1.1+Connect.
//   - Versioning & compatibility: Breaking changes are gated by the package
//     version suffix (v1, v2), and buf breaking-change detection runs against
//     the schema to prevent accidental regressions.
//   - Single import, strong typing: Consumers import one Go module
//     (go.pilab.hu/cloud/virtpb) and get fully typed request/response messages
//     and service interfaces — no manual marshalling or ad-hoc API clients.
//   - OpenTelemetry integration: Generated code carries observability metadata
//     for distributed tracing and metrics out of the box.
//
// # Usage
//
//	import "go.pilab.hu/cloud/virtpb"
//
// Sub-packages can be imported individually, for example:
//
//	import pivirtdv1 "go.pilab.hu/cloud/virtpb/pilab/pivirtd/v1"
//	import networkv2 "go.pilab.hu/cloud/virtpb/pilab/network/v2"
//
// # Proto Source
//
// The canonical .proto source files live under the proto/ directory and are
// managed with buf (https://buf.build). Generated Go code is checked into this
// repository so that consumers can depend on it directly without a protoc
// toolchain. To regenerate, run:
//
//	buf generate
//
// # Module
//
//	module: go.pilab.hu/cloud/virtpb
//	version: v1.0.4 (see go.mod)
package virtpb
