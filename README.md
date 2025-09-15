# Zero-Touch Provisioning MCP Server

A Model Context Protocol (MCP) server that provides seamless integration with Ubuntu MAAS (Metal as a Service) for zero-touch provisioning of machines and virtual hosts.

## 🚀 Features

- **Machine Management**: List and query information about physical machines in your MAAS environment
- **VM Host Operations**: Manage virtual machine hosts including listing, querying, and composing new VMs
- **OAuth 1.0 Authentication**: Secure communication with MAAS API using OAuth 1.0 with PLAINTEXT signature
- **Multiple Transport Modes**: Support for stdio, HTTP, and SSE transport protocols
- **Structured Logging**: Comprehensive logging with Zap logger for monitoring and debugging

## 📋 Prerequisites

- Go 1.23.3 or later
- Ubuntu MAAS instance with API access
- Valid MAAS API credentials

## 🛠️ Installation

1. Clone the repository:
```bash
git clone https://github.com/JarcauCristian/ztp-mcp.git
cd ztp-mcp
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the project:
```bash
go build -o ztp-mcp cmd/main.go
```

## ⚙️ Configuration

Set the following environment variables:

```bash
# Required: MAAS configuration
export MAAS_BASE_URL="https://your-maas-server.com"
export MAAS_API_KEY="consumer_key:token:secret"

# Optional: MCP server configuration
export MCP_TRANSPORT="stdio"  # Options: stdio, http, sse
export MCP_ADDRESS=":8080"    # Required for http/sse modes
```

### MAAS API Key Format

The `MAAS_API_KEY` must be in the format: `consumer_key:token:secret`

You can obtain your MAAS API key from your MAAS web interface under your user preferences.

## 🚀 Usage

### Stdio Mode (Default)
```bash
./ztp-mcp
```

### HTTP Mode
```bash
export MCP_TRANSPORT=http
export MCP_ADDRESS=":8080"
./ztp-mcp
```

### SSE Mode
```bash
export MCP_TRANSPORT=sse
export MCP_ADDRESS=":8080"
./ztp-mcp
```

## 🔧 Available Tools

### Machine Operations

#### `list_machines`
List all available machines in the MAAS environment.

**Usage:**
- No parameters required
- Returns: JSON array of machine objects with their current status, configuration, and metadata

### VM Host Operations

#### `list_vm_hosts`
Retrieve all VM hosts available in the MAAS environment.

**Usage:**
- No parameters required
- Returns: JSON array of VM host objects with their capabilities and current utilization

#### `list_vm_host`
Get detailed information about a specific VM host.

**Parameters:**
- `id` (required): The numeric ID of the VM host to query

**Usage:**
```json
{
  "id": "1"
}
```

#### `compose_vm_host`
Create a new virtual machine on a specified VM host.

**Parameters:**
- `id` (required): ID of the VM host to compose the machine on
- `cores` (required): Number of CPU cores for the VM
- `memory` (required): RAM allocation in MiB
- `storage` (required): Storage allocation in GB
- `hostname` (required): Name for the new VM (alphanumeric, dots, and hyphens allowed)

**Usage:**
```json
{
  "id": "1",
  "cores": "4",
  "memory": "8192",
  "storage": "100",
  "hostname": "my-new-vm"
}
```

## 📁 Project Structure

```
ztp-mcp/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   └── server/
│       ├── maas_client/
│       │   └── maas-client.go  # MAAS API client with OAuth 1.0 support
│       ├── parser/
│       │   └── parse.go        # URI parsing utilities
│       └── tools/
│           ├── tool.go         # MCP tool interface definition
│           ├── machines.go     # Machine management tools
│           └── vm-hosts.go     # VM host management tools
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── LICENSE                     # MIT License
└── README.md                   # This file
```

## 🔐 Security

- Uses OAuth 1.0 with PLAINTEXT signature method for MAAS API authentication
- Secure credential management through environment variables
- Request timeouts and proper error handling
- No sensitive data stored in code or logs

## 📝 API Response Format

All tools return responses in the following JSON structure:

```json
{
  "Body": "response_content",
  "StatusCode": 200,
  "Headers": {
    "Content-Type": ["application/json"],
    "..."
  }
}
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Troubleshooting

### Common Issues

1. **Authentication Errors**: Verify your `MAAS_API_KEY` format and credentials
2. **Connection Issues**: Check your `MAAS_BASE_URL` and network connectivity
3. **Permission Errors**: Ensure your MAAS user has appropriate permissions for the operations you're trying to perform

### Logging

The server uses structured logging with different log levels. Set the log level using:

```bash
export LOG_LEVEL=debug  # Options: debug, info, warn, error
```

## 📞 Support

For issues and questions:
- Create an issue in this repository
- Check existing issues for similar problems
- Provide logs and configuration details when reporting bugs