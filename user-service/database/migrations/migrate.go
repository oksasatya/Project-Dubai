package migrations

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"user-service/database/seeder"
)

// Migrate is a function to migrate all tables
func Migrate(db *mongo.Database) error {
	migrations := []*Migration{
		createUsersCollectionMigration(db, "email"),
	}
	autoMigrate := os.Getenv("AUTO_MIGRATE")
	autoDrop := os.Getenv("AUTO_DROP")

	if autoDrop == "true" && autoMigrate == "true" {
		logrus.Println("Running AutoDrop (Rollback all migrations) and AutoMigrate...")

		for i := len(migrations) - 1; i >= 0; i-- {
			if err := migrations[i].Rollback(); err != nil {
				return fmt.Errorf("rollback migration %s failed: %v", migrations[i].ID, err)
			}
		}

		for _, migration := range migrations {
			if err := migration.Migrate(); err != nil {
				return fmt.Errorf("migration failed after drop: %v", err)
			}
		}

		// Seeders
		seeder.SeedAll(db)

		logrus.Println("AutoMigrate and Seeders completed.")
	} else if autoDrop == "true" {
		logrus.Println("Running AutoDrop (Rollback all migrations)...")

		for i := len(migrations) - 1; i >= 0; i-- {
			if err := migrations[i].Rollback(); err != nil {
				return fmt.Errorf("rollback migration %s failed: %v", migrations[i].ID, err)
			}
		}
	} else if autoMigrate == "true" {
		logrus.Println("Running AutoMigrate...")
		for _, migration := range migrations {
			if err := migration.Migrate(); err != nil {
				return fmt.Errorf("migration failed: %v", err)
			}
		}
		seeder.SeedAll(db)
		logrus.Println("AutoMigrate completed.")
	} else {
		logrus.Println("Skipping AutoMigrate and AutoDrop.")
	}
	return nil
}
