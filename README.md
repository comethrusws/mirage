# Mirage

Mirage is an API mocking gateway and traffic recorder written in Go. It sits between your application and external APIs, allowing you to intercept, record, and mock HTTP requests during development.

## Features

- **HTTP Proxy**: Forward requests to real APIs transparently.
- **Mocking**: Define scenarios to return custom responses based on path, method, and headers.
- **Recording**: Capture real API traffic and save it to JSON for inspection or replay.
- **Web Dashboard**: View intercepted requests in real-time and toggle scenarios on/off.
- **CLI**: Simple command-line interface.

## Installation

```bash
go build -o mirage ./cmd/mirage
```

## Usage

### Start Proxy

Start the proxy server on port 8080:

```bash
./mirage start
```

### Start with Scenarios

Load a configuration file with mock scenarios:

```bash
./mirage start --config examples/config.yaml
```

The Web Dashboard is available at [http://localhost:8080/__mirage/](http://localhost:8080/__mirage/).

### Record Traffic

Start the proxy in recording mode. All traffic will be saved to `traffic.json`.

```bash
./mirage record --output my_traffic.json
```

### Replay Traffic

Replay recorded requests (fire-and-forget):

```bash
./mirage replay my_traffic.json
```

## Configuration

Defines scenarios in a YAML file:

```yaml
scenarios:
  - name: user-success
    match:
      path: /api/users
      method: GET
    response:
      status: 200
      body: '{"users": [{"id": 1, "name": "Alice"}]}'
      headers:
        Content-Type: application/json
```

## Dashboard

Visit `http://localhost:8080/__mirage/` to see the request log and enable/disable scenarios dynamically.
