# Laravel Server Setup

A Go tool to automate the setup of a production-ready Laravel server on Ubuntu.

## Features

- System update and essential packages installation
- PHP 8.3 installation with optimized configuration
- MySQL installation and secure configuration
- Nginx installation with optimized configuration for Laravel
- Security hardening (firewall, fail2ban, SSH)
- Laravel application setup from Git repository
- Service configuration and startup
- SSL certificate setup with Let's Encrypt

## Requirements

- Ubuntu server (tested on Ubuntu 20.04 LTS)
- Go 1.18 or higher (for building from source)
- Non-root user with sudo privileges

## Installation

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/laravel-setup.git
   cd laravel-setup
   ```

2. Build the binary:
   ```
   make build
   ```

   Or build specifically for Linux:
   ```
   make build-linux
   ```

   Or build specifically for x64 architecture (macOS Intel):
   ```
   make build-x64
   ```

   Or build specifically for ARM64 architecture (Apple Silicon):
   ```
   make build-arm64
   ```

   Or build specifically for x86 architecture (32-bit):
   ```
   make build-x86
   ```

3. (Optional) Install the binary to /usr/local/bin:
   ```
   make install
   ```

## Usage

Run the tool:

```
./laravel-setup
```

Or if you installed it to /usr/local/bin:

```
laravel-setup
```

The tool will guide you through the setup process, asking for:

- Domain name for your Laravel project
- Git repository URL for your Laravel project

## Project Structure

The project is organized into modules for better maintainability:

- `cmd/laravel-setup`: Main entry point
- `pkg/config`: Configuration structures and functions
- `pkg/utils`: Utility functions
- `pkg/system`: System update and essential packages installation
- `pkg/php`: PHP installation and configuration
- `pkg/mysql`: MySQL installation and configuration
- `pkg/nginx`: Nginx installation and configuration
- `pkg/security`: Security configurations
- `pkg/laravel`: Laravel application setup
- `pkg/services`: Service configuration and startup
- `pkg/templates`: Configuration templates

## Security Considerations

This tool implements several security measures:

- Firewall configuration with UFW
- Intrusion prevention with fail2ban
- SSH hardening (custom port, key-based authentication)
- MySQL secure installation
- Nginx security headers and rate limiting

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
