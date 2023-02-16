package packager

import (
	"github.com/AubreeH/goApiDb/database"
)

type Package struct {
	Database *database.Database
}

func Setup(config database.DatabaseConfig) (*Package, error) {
	db, err := database.SetupDatabase(config)
	if err != nil {
		return nil, err
	}

	output := &Package{Database: db}

	return output, nil
}
