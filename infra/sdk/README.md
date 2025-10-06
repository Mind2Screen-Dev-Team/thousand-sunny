# 🔌 SDK Modules

⚡ **Power up your integrations** with SDK infrastructure – where external APIs become seamless crew members!

The SDK modules provide a **standardized, scalable foundation** for integrating third-party services and APIs. Built with **Uber FX dependency injection** and following **clean architecture principles**, each SDK is a self-contained module that can be easily plugged into your application.

---

## 🏗️ Architecture Overview

```
infra/sdk/
├── sdk.fx.modules.go           # Main SDK module orchestrator
├── <sdk-name>/                 # Individual SDK packages
│   ├── sdk.<sdk-name>.<action>.go    # SDK implementation
│   └── sdk.<sdk-name>.fx.modules.go  # SDK-specific FX modules
└── README.md                   # This documentation
```

---

## 📋 Features

* 🔧 **Uber FX Integration** – Seamless dependency injection and lifecycle management
* 🏗️ **Modular Design** – Each SDK is self-contained and independently manageable  
* 🔄 **Standardized Interface** – Consistent patterns across all SDK implementations
* 🚀 **Hot-Pluggable** – Add/remove SDKs without affecting core application
* 📊 **Built-in Observability** – Integrated logging, tracing, and metrics
* ⚙️ **Configuration-Driven** – Environment-based SDK configuration
* 🛡️ **Error Handling** – Robust error handling and retry mechanisms

---

## 🚀 Quick Start

### Adding a New SDK

1. **Create SDK directory structure:**
```bash
mkdir -p infra/sdk/example-api
```

2. **Implement the SDK:**
```go
// infra/sdk/example-api/sdk.example-api.client.go
package exampleapi

import (
    "context"
    "net/http"
)

type Client struct {
    httpClient *http.Client
    baseURL    string
    apiKey     string
}

func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        httpClient: &http.Client{},
        baseURL:    baseURL,
        apiKey:     apiKey,
    }
}

func (c *Client) GetData(ctx context.Context, id string) (*Data, error) {
    // Implementation here
    return nil, nil
}
```

3. **Create FX module:**
```go
// infra/sdk/example-api/sdk.example-api.fx.modules.go
package exampleapi

import (
    "go.uber.org/fx"
)

var Modules = fx.Options(
    fx.Module("sdk:example-api",
        fx.Provide(NewClient),
    ),
)
```

4. **Register in main SDK modules:**
```go
// infra/sdk/sdk.fx.modules.go
var Modules = fx.Options(
    fx.Module("sdk:modules",
        exampleapi.Modules,
        // Add other SDK modules here
    ),
)
```

---

## 🎯 Usage Patterns

### Basic SDK Implementation

```go
type SDKAPI interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    HealthCheck(ctx context.Context) error
}

type BaseSDK struct {
    config Config
    logger Logger
    tracer Tracer
}

func (s *BaseSDK) Connect(ctx context.Context) error {
    // some logic here
    return nil
}
```

### Configuration Pattern

```yaml
# config.yaml
provider:
  example.one:
    base.url: "https://api.example.com/api/v1"
    options:
      client.id: "example.one.id"
      client.secret: "example.one.secret"
```

### Dependency Injection

```go
// In your service
type MyServiceAPI interface {}

type MyService struct {
    exampleAPI *exampleapi.Client
}

func NewMyService(exampleAPI *exampleapi.Client) MyServiceAPI {
    return &MyService{
        exampleAPI: exampleAPI,
    }
}

// FX Module
var ServiceModules = fx.Options(
    fx.Provide(NewMyService),
)
```

---

## 🛠️ Best Practices

### 1. **Consistent Naming Convention**
- Package: `<sdkname>` (lowercase, no hyphens)
- Files: `sdk.<sdk-name>.<action>.go`
- Modules: `sdk.<sdk-name>.fx.modules.go`

### 2. **Context Propagation**
Always accept and propagate `context.Context` for:
- Request cancellation
- Timeout handling  
- Tracing correlation
- Request-scoped values
