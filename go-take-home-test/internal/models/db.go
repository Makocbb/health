package models

// MigrationOptions holds configuration for the migration service
type MigrationOptions struct {
	VersionFilePath string
	MigrationsPath  string
}

func NewMigrationOptions(inputOptions ...MigrationOption) *MigrationOptions {
	opts := &MigrationOptions{}
	for _, opt := range inputOptions {
		opt(opts)
	}
	return opts
}

// MigrationOption is a functional option for configuring MigrationOptions
type MigrationOption func(*MigrationOptions)

// MigrationWithVersionFilePath sets the path to the version file
func MigrationWithVersionFilePath(path string) MigrationOption {
	return func(opts *MigrationOptions) {
		opts.VersionFilePath = path
	}
}

// MigrationWithMigrationsPath sets the path to the migrations directory
func MigrationWithMigrationsPath(path string) MigrationOption {
	return func(opts *MigrationOptions) {
		opts.MigrationsPath = path
	}
}
