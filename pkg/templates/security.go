package templates

import "fmt"

// GetFail2banConfig returns the Fail2ban configuration
// This configures Fail2ban to protect against brute force attacks
func GetFail2banConfig(sshPort string) string {
	return fmt.Sprintf(`[sshd]
enabled = true
port = %s
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600
findtime = 600

[nginx-http-auth]
enabled = true
filter = nginx-http-auth
logpath = /var/log/nginx/error.log
maxretry = 3
bantime = 3600

[nginx-limit-req]
enabled = true
filter = nginx-limit-req
logpath = /var/log/nginx/error.log
maxretry = 10
bantime = 600`, sshPort)
}

// GetSSHConfig returns the SSH security configuration
// This hardens SSH to prevent unauthorized access
func GetSSHConfig(sshPort string) string {
	return fmt.Sprintf(`# Security configurations
Port %s
PermitRootLogin no
PasswordAuthentication yes
PubkeyAuthentication yes
AuthorizedKeysFile .ssh/authorized_keys
PermitEmptyPasswords no
ChallengeResponseAuthentication no
UsePAM yes
X11Forwarding no
PrintMotd no
ClientAliveInterval 300
ClientAliveCountMax 2
MaxAuthTries 3
MaxSessions 2
Protocol 2`, sshPort)
}

// GetServerInfoContent returns the server information content
// This provides a summary of the server configuration for reference
func GetServerInfoContent(domain, webRoot, dbName, dbUser, sshPort, username string) string {
	return fmt.Sprintf(`===========================================
Laravel Production Server Setup Complete
===========================================

Domain: %s
Web Directory: %s

Database Information:
- Database Name: %s
- Database User: %s
- Database Password: See mysql_credentials.txt file

Important Security Notes:
- SSH Port changed to: %s
- Firewall (UFW) is enabled
- Fail2ban is configured

Service Status Commands:
- sudo systemctl status nginx
- sudo systemctl status php8.3-fpm
- sudo systemctl status mysql
- sudo systemctl status redis-server
- sudo systemctl status supervisor

Log Locations:
- Nginx: /var/log/nginx/
- PHP-FPM: /var/log/php8.3-fpm.log
- MySQL: /var/log/mysql/
- Laravel: %s/storage/logs/

Security Tools:
- UFW Firewall: sudo ufw status
- Fail2ban: sudo fail2ban-client status

SSH Connection (remember the new port):
ssh -p %s %s@your-server-ip
`, domain, webRoot, dbName, dbUser, sshPort, webRoot, sshPort, username)
}