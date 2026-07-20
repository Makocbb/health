package models

// MigrationOptions holds configuration for the migration service
type MigrationOptions struct {
	VersionFilePath string
	MigrationsPath  string
}

// MigrationOption is a functional option for configuring MigrationOptions
type MigrationOption func(*MigrationOptions)

// WithVersionFilePath sets the path to the version file
func WithVersionFilePath(path string) MigrationOption {
	return func(opts *MigrationOptions) {
		opts.VersionFilePath = path
	}
}

// WithMigrationsPath sets the path to the migrations directory
func WithMigrationsPath(path string) MigrationOption {
	return func(opts *MigrationOptions) {
		opts.MigrationsPath = path
	}
}
