package config

// Config holds all the settings for the Laravel setup
type Config struct {
	Domain         string
	RepoURL        string
	DBName         string
	DBUser         string
	DBPassword     string
	DBRootPassword string
	WebUser        string
	SSHPort        string
	WebRoot        string
	ScriptDir      string
}

// NewConfig initializes a new configuration with default values
func NewConfig() *Config {
	return &Config{
		DBName:  "production_db",
		DBUser:  "db_user",
		SSHPort: "2222",
		WebUser: "www-data",
	}
}