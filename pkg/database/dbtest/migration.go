package dbtest

import (
	"context"
	"fmt"
	"io/fs"

	"ariga.io/atlas-go-sdk/atlasexec"
)

func Migrate(url string, dir fs.FS) error {
	// Define the execution context, supplying a migration directory
	// and potentially an `atlas.hcl` configuration file using `atlasexec.WithHCL`.
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(dir),
	)
	if err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	// atlasexec works on a temporary directory, so we need to close it
	defer workdir.Close()

	// Initialize the client.
	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Run `atlas migrate apply`
	_, err = client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: url,
	})
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
