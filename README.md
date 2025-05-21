<p align="center">
    <img src="images/pro_logo_light.svg" width="250" alt="Kraken Logo">
</p>
<h1 align="center">Kraken Go API Connector</h1>

A Go module for interacting with the Kraken Spot and Derivatives APIs, offering both REST and WebSocket access to spot and futures exchanges. This library includes utilities for fixed-point arithmetic, callback handlers, and API helpers.

Key features include:

* Trading operations: add, amend, and cancel orders

* Real-time market data streams

* Book builders for L2 and L3 order books with checksum validation

* Access to private account data such as balances and execution reports

* Retrieval of instruments and assets

* Utility functions for asset name normalization and price derivation.

## 🚀 Installation

To install the module in your Go project:
```bash
$ go get github.com/krakenfx/api-go
```

If you're interested in running examples or contributing to development, clone the repo and install dependencies:

```bash
git clone https://github.com/krakenfx/api-go.git
cd api-go
go mod download
```

## 📂 Examples & Scripts

Example packages using the module are located in the [examples](examples/) directory

### 🌐 Environment Configuration

Before running the scripts, you'll need to set the following environment variables to configure API endpoints and authentication keys:

```bash
# REST and WebSocket URLs for Kraken Spot API
export KRAKEN_API_SPOT_REST_URL="https://api.kraken.com"
export KRAKEN_API_SPOT_WS_URL="wss://ws.kraken.com/v2"
export KRAKEN_API_SPOT_WS_AUTH_URL="wss://ws-auth.kraken.com/v2"
export KRAKEN_API_SPOT_PUBLIC=""    # Insert your Spot public API key
export KRAKEN_API_SPOT_SECRET=""    # Insert your Spot secret API key

# REST and WebSocket URLs for Kraken Futures API
export KRAKEN_API_FUTURES_REST_URL="https://futures.kraken.com"
export KRAKEN_API_FUTURES_WS_URL="wss://futures.kraken.com/ws/v1"
export KRAKEN_API_FUTURES_PUBLIC=""  # Insert your Futures public API key
export KRAKEN_API_FUTURES_SECRET=""  # Insert your Futures secret API key
```

### 🔐 Running Authenticated Examples

To simplify execution of authenticated example scripts, use the provided helper script:

```bash
./scripts/run prod ./examples/spotrest/recenttrades
```

This command sets the correct environment and runs the selected script with authentication enabled.

## 🧱 Dependencies

This module uses the following third-party packages:

* [gorilla/websocket](https://github.com/gorilla/websocket): WebSocket protocol implementation in Go.

* [google/uuid](https://github.com/google/uuid): RFC 4122-compliant UUID generation.

* [golang.org/x/sync](https://pkg.go.dev/golang.org/x/sync): Additional Go concurrency primitives.

## 📫 Contributing

Feel free to open issues, suggest features, or submit PRs.