<p align="center">
  <img src="./assets/banner.svg" alt="Mirage Banner" />
</p>

<p align="center">
  <strong>API Mocking Gateway & Traffic Recorder</strong>
</p>

<p align="center">
  <a href="https://github.com/comethrusws/mirage/actions/workflows/ci.yml"><img src="https://github.com/comethrusws/mirage/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/comethrusws/mirage/releases"><img src="https://img.shields.io/github/v/release/comethrusws/mirage" alt="Release"></a>
  <a href="https://github.com/comethrusws/mirage/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License"></a>
  <a href="https://goreportcard.com/report/github.com/comethrusws/mirage"><img src="https://goreportcard.com/badge/github.com/comethrusws/mirage" alt="Go Report"></a>
</p>

---

## Overview

Mirage is a powerful HTTP/HTTPS proxy server designed for local development and testing. It intercepts API requests, allowing you to mock responses, simulate network conditions, and record traffic for replay.

### Key Features

- 🎭 **Request Mocking** - Define scenarios to return custom responses
- 📝 **Traffic Recording** - Capture and replay real API interactions  
- 🎯 **Pattern Matching** - Match by path (glob), method, and headers
- ⏱️ **Latency Simulation** - Add delays to simulate slow networks
- 🔄 **Scenario Switching** - Toggle mocks on/off in real-time
- 📊 **Web Dashboard** - Monitor requests with a clean UI
- 🚀 **Zero Dependencies** - Single binary, cross-platform
- ⚡ **High Performance** - Built with Go for speed and concurrency

## Installation

### Homebrew (macOS/Linux)

```bash
brew install comethrusws/mirage/mirage
```

### Installation Script

```bash
curl -sL https://raw.githubusercontent.com/comethrusws/mirage/main/install.sh | bash
```

### Go Install

```bash
go install github.com/comethrusws/mirage/cmd/mirage@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/comethrusws/mirage/releases).

### Build from Source

```bash
git clone https://github.com/comethrusws/mirage.git
cd mirage
make build
sudo make install
```

## Quick Start

### Start the Proxy

```bash
mirage start
```

Dashboard available at http://localhost:8080/__mirage/

### With Mocking Scenarios

```bash
mirage start --config examples/config.yaml
```

### Record Traffic

```bash
mirage record --output traffic.json
```

### Replay Traffic

```bash
mirage replay traffic.json
```

## Configuration

Create a YAML file to define mock scenarios:

```yaml
scenarios:
  - name: success-response
    match:
      path: /api/users
      method: GET
    response:
      status: 200
      headers:
        Content-Type: application/json
      body: |
        {
          "users": [
            {"id": 1, "name": "Alice"},
            {"id": 2, "name": "Bob"}
          ]
        }
      delay: 100ms

  - name: server-error
    match:
      path: /api/*
      headers:
        X-Test-Error: "true"
    response:
      status: 500
      body: "Internal Server Error"
      delay: 2s

  - name: not-found
    match:
      path: /api/missing
      method: POST
    response:
      status: 404
      body: '{"error": "Not found"}'
```

### Pattern Matching

- **Path**: Supports glob patterns (`/api/*`, `/users/*/profile`)
- **Method**: GET, POST, PUT, DELETE, PATCH, etc.
- **Headers**: Match specific header values

## Usage Examples

### Development Workflow

```bash
mirage start --config dev-api.yaml --port 8080
```

Configure your app to proxy through `http://localhost:8080`

### Testing Different Scenarios

```bash
mirage scenarios list config.yaml
```

Toggle scenarios via the dashboard or CLI.

### Capturing Real Traffic

```bash
mirage record --output prod-traffic.json
```

Later replay for testing:

```bash
mirage replay prod-traffic.json
```

## Dashboard

Access the web dashboard at `http://localhost:8080/__mirage/`

Features:
- Real-time request log
- Scenario management (enable/disable)
- Request/response details
- Performance metrics

## CLI Reference

```
mirage start [flags]              Start proxy server
mirage record [flags]             Record traffic mode
mirage replay <file>              Replay recorded traffic
mirage scenarios list <config>    List scenarios in config
```

### Flags

```
-p, --port int       Port to run on (default 8080)
-c, --config string  Path to config file
-o, --output string  Output file for recordings
```

## Architecture

```
┌─────────────┐         ┌─────────────┐         ┌─────────────┐
│             │         │             │         │             │
│  Your App   ├────────▶│   Mirage    ├────────▶│  Real API   │
│             │         │             │         │             │
└─────────────┘         └─────────────┘         └─────────────┘
                              │
                              │
                        ┌─────▼─────┐
                        │ Dashboard │
                        │ Recorder  │
                        │ Scenarios │
                        └───────────┘
```

## Use Cases

- **API Development** - Mock external APIs during development
- **Testing** - Simulate error conditions and edge cases
- **Integration Testing** - Record production traffic for test suites
- **Network Simulation** - Test app behavior with slow/unreliable connections
- **Debugging** - Inspect and modify API requests/responses in real-time

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

For security concerns, see [SECURITY.md](SECURITY.md).

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

Built with ❤️ using:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP routing
- [YAML v3](https://github.com/go-yaml/yaml) - Configuration parsing

---

<p align="center">
  <sub>⭐ Star the project if you find it useful!</sub>
</p>
