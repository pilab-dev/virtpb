# virtpb 🚀

[![Go Reference](https://pkg.go.dev/badge/go.pilab.hu/cloud/virtpb.svg)](https://pkg.go.dev/go.pilab.hu/cloud/virtpb)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`virtpb` is a Go client library containing Protocol Buffer definitions and generated code for virtualized assets management. This library provides the necessary types and interfaces for interacting with virtualized infrastructure services.

## ✨ Features

- 📦 Protocol Buffer definitions for virtualized assets management
- 🔄 Generated Go code for easy integration
- 🔌 Support for gRPC and Connect protocols
- 📊 OpenTelemetry integration for observability
- 🏗️ Comprehensive type definitions for virtualized infrastructure

## 📥 Installation

```bash
go get go.pilab.hu/cloud/virtpb
```

## 💻 Usage

Import the package in your Go code:

```go
import "go.pilab.hu/cloud/virtpb"
```

The generated code provides strongly-typed interfaces and structures for working with virtualized assets. Refer to the [GoDoc documentation](https://pkg.go.dev/go.pilab.hu/cloud/virtpb) for detailed API reference.

## 📦 Dependencies

- 🐹 Go 1.23.0 or later
- 📜 Protocol Buffers
- 🔌 gRPC
- 🔗 Connect
- 📊 OpenTelemetry

## 🔧 Development

This repository contains the generated Protocol Buffer code. The source `.proto` files are maintained in a separate repository. To regenerate the Go code:

1. ⚙️ Install the Protocol Buffer compiler and Go plugins
2. 🔄 Run the protobuf generation command
3. 📝 Update the generated code in this repository

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 💬 Support

For support, please open an issue in the repository or contact the maintainers. 