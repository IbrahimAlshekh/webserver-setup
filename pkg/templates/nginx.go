package templates

import "fmt"

// GetNginxConfig returns the Nginx configuration for a Laravel application
// This configures Nginx with security headers, gzip compression, and rate limiting
func GetNginxConfig(domain, webRoot string) string {
	return fmt.Sprintf(`server {
    listen 80;
    server_name %s www.%s;
    root %s/public;
    index index.php index.html index.htm;

    # Security headers to protect against common web vulnerabilities
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # Hide Nginx version for security
    server_tokens off;

    # Gzip compression for better performance
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;

    # Main location block
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    # PHP processing
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php8.4-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
        fastcgi_hide_header X-Powered-By;
    }

    # Static files with long cache times for better performance
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Rate limiting for login to prevent brute force attacks
    location /login {
        limit_req zone=login burst=5 nodelay;
        try_files $uri $uri/ /index.php?$query_string;
    }

    # Rate limiting for API to prevent abuse
    location /api {
        limit_req zone=api burst=20 nodelay;
        try_files $uri $uri/ /index.php?$query_string;
    }

    # Security: deny access to hidden files
    location ~ /\. {
        deny all;
    }
}`, domain, domain, webRoot)
}
