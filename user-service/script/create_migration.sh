#!/bin/bash

# This script creates a new migration file with the provided name as an argument.
# It creates a new migration file in the directory "database/migrations"
# and adds a basic template for MongoDB migration functions with dynamic index fields.
#
# To use this script, you need to give execute permissions:
# chmod +x script/create_migration.sh
#
# Run the script to create a new migration file with the provided migration name as an argument:
# ./script/create_migration.sh create_yourCollection_collection

# Check if a migration name is provided as an argument
if [ -z "$1" ]; then
  echo "Usage: $0 migration_name"
  exit 1
fi

# Get the current timestamp in the format YYYYMMDDHHMMSS
timestamp=$(date +"%Y%m%d%H%M%S")

# The migration name from the first argument
migration_name=$1

# Create the migration filename with the format timestamp_migration_name.go
filename="${timestamp}_${migration_name}.go"

# Convert the migration name to a camelCase function name
function_name="$(echo ${migration_name} | awk -F_ '{for (i=1; i<=NF; i++) {if (i == 1) {printf tolower($i)} else {printf toupper(substr($i,1,1)) tolower(substr($i,2))}}}')Migration"

# Migration ID with the format timestamp_migration_name
migration_id="${timestamp}_${migration_name}"

# Extract the collection name from the migration name (take the second part after 'create_')
collection_name=$(echo ${migration_name} | awk -F_ '{print $2}')

# Directory where the migration file will be created
migration_dir="database/migrations"

# Check if the migrations directory exists
if [ ! -d "$migration_dir" ]; then
  echo "Migration directory $migration_dir does not exist. Creating..."
  mkdir -p "$migration_dir"
fi

# Create the migration file in the specified directory
touch "${migration_dir}/${filename}"

# Add the basic template for MongoDB migration to the new file
cat <<EOL > "${migration_dir}/${filename}"
package migrations

import (
	"context"
	"time"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Migration is a struct to define migration
type Migration struct {
	ID       string
	Migrate  func() error
	Rollback func() error
}

// Migration function for $migration_name
func $function_name(database *mongo.Database, indexField string) *Migration {
	return &Migration{
		ID: "$migration_id",
		Migrate: func() error {
			collection := database.Collection("$collection_name")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			indexOptions := options.Index().SetUnique(true)
			indexModel := mongo.IndexModel{
				Keys:    bson.M{indexField: 1},
				Options: indexOptions,
			}

			_, err := collection.Indexes().CreateOne(ctx, indexModel)
			if err != nil {
				return err
			}

			logrus.Printf("Migration: %s completed. Index created on field: %s", "$migration_name", indexField)
			return nil
		},
		Rollback: func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := database.Collection("$collection_name").Drop(ctx)
			if err != nil {
				return err
			}

			logrus.Printf("Rollback: %s completed", "$migration_name")
			return nil
		},
	}
}
EOL

echo "Created migration file: ${migration_dir}/${filename}"
