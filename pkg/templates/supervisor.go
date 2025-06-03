package templates

import "fmt"

// GetSupervisorConfig returns the Supervisor configuration for Laravel queue workers
// This ensures Laravel queue jobs are processed reliably and automatically restarted if they fail
func GetSupervisorConfig(webRoot, webUser string) string {
	return fmt.Sprintf(`[program:laravel-worker]
process_name=%%(program_name)s_%%(process_num)02d
command=php %s/artisan queue:work --sleep=3 --tries=3 --max-time=3600
autostart=true
autorestart=true
stopasgroup=true
killasgroup=true
user=%s
numprocs=2
redirect_stderr=true
stdout_logfile=%s/storage/logs/worker.log
stopwaitsecs=3600`, webRoot, webUser, webRoot)
}