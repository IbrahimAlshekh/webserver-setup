package templates

import "fmt"

// GetMySQLConfig returns the MySQL configuration SQL script
// This configures MySQL for Laravel with appropriate user permissions and security settings
func GetMySQLConfig(dbName, dbUser, dbPassword, dbRootPassword string) string {
	return fmt.Sprintf(`
-- Create database with UTF-8 support for Laravel
CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Create Laravel user with restricted permissions
DELETE FROM mysql.user WHERE User='%s' AND Host='localhost';
CREATE USER IF NOT EXISTS '%s'@'localhost' IDENTIFIED BY '%s';
GRANT ALL PRIVILEGES ON %s.* TO '%s'@'localhost';

-- Create admin user for easier database management
DELETE FROM mysql.user WHERE User='admin' AND Host='localhost';
CREATE USER IF NOT EXISTS 'admin'@'localhost' IDENTIFIED BY '%s';
GRANT ALL PRIVILEGES ON *.* TO 'admin'@'localhost' WITH GRANT OPTION;

-- Secure the installation by removing anonymous users and test database
DELETE FROM mysql.user WHERE User='';
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');
DROP DATABASE IF EXISTS test;
DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%%';

-- Flush privileges to apply changes
FLUSH PRIVILEGES;
`, dbName, dbUser, dbPassword, dbName, dbUser, dbRootPassword)
}

// GetMySQLCredentialsContent returns the MySQL credentials content
// This provides a secure way to store and access MySQL credentials
func GetMySQLCredentialsContent(dbName, dbUser, dbPassword, dbRootPassword string) string {
	return fmt.Sprintf(`MySQL Credentials:
==================
Laravel Database: %s
Laravel User: %s
Laravel Password: %s

Admin User: admin
Admin Password: %s

Connection Examples:
mysql -u %s -p%s %s
mysql -u admin -p%s
`, dbName, dbUser, dbPassword, dbRootPassword, dbUser, dbPassword, dbName, dbRootPassword)
}
