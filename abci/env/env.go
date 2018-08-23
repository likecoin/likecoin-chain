package env

import "os"

// IsProduction returns true in production environment
var IsProduction = os.Getenv("ENV") == "production"
