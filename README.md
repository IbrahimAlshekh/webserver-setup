# Laravel Server Setup

A Go tool to automate the setup of a production-ready Laravel server on Ubuntu.

## Features

- System update and essential packages installation
- PHP 8.4 installation with optimized configuration
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

4. (Optional) Upload the binary and example configuration file to a remote server:
   ```
   make upload
   ```
   This will build the Linux binary and upload both the binary and the example configuration file to the server.

## Usage

Run the tool:

```
./laravel-setup
```

Or if you installed it to /usr/local/bin:

```
laravel-setup
```

### Module Selection

You can choose which modules to run by using command-line flags. This is useful if you've already completed some steps and don't want to repeat them:

```
laravel-setup --skip-mysql --skip-nginx
```

Available skip flags:

- `--skip-system-update`: Skip system update step
- `--skip-essentials`: Skip installing essential packages
- `--skip-php`: Skip PHP installation
- `--skip-mysql`: Skip MySQL installation
- `--skip-nginx`: Skip Nginx installation
- `--skip-security`: Skip security configuration
- `--skip-laravel`: Skip Laravel setup
- `--skip-services`: Skip services configuration

### Configuration File

You can use a TOML configuration file to store your settings and skip flags. The tool will look for a `config.toml` file in your home directory by default, or you can specify a custom path:

```
laravel-setup --config-path=/path/to/config.toml
```

The configuration file can include all settings and skip flags. The tool will only prompt for values that are not defined in the config file. Command-line flags take precedence over configuration file settings.

Example `config.toml`:

```toml
# Laravel Setup Configuration

# Basic settings
Domain = "example.com"
RepoURL = "https://github.com/user/laravel-project.git"
DBName = "production_db"
DBUser = "db_user"
DBPassword = "your-secure-password"  # Leave empty to generate a random password
DBRootPassword = "your-secure-root-password"  # Leave empty to generate a random password
WebUser = "www-data"
SSHPort = "2222"
WebRoot = "/var/www/example.com"  # Leave empty to use /var/www/[Domain]

# Skip flags - set to true to skip the corresponding step
SkipSystemUpdate = false
SkipEssentials = false
SkipPHP = false
SkipMySQL = false
SkipNginx = false
SkipSecurity = false
SkipLaravel = false
SkipServices = false
```

A sample configuration file is available in the `examples` directory.

### Cleanup

To clean up temporary files created during the setup process:

```
laravel-setup --cleanup
```

### Setup Process

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
