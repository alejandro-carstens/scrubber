package cmd

import (
	"errors"
	"strings"

	"scrubber/app/models"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

var tableRegister map[string]models.Modelable = map[string]models.Modelable{
	"users":           &models.User{},
	"access_controls": &models.AccessControl{},
}

type runMigrations struct {
	action string
	tables []string
}

func (rm *runMigrations) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "run-migrations",
		Short: "runs the MySQL table migrations",
		Args:  rm.Validate,
		Run:   rm.Handle,
	}

	command.Flags().String("action", "migrate", "whether to migrate or drop tables, columns, and indices")
	command.Flags().StringSlice("tables", []string{""}, "tables to be migrated, if empty it will default to 'all'")

	return command
}

// Handle implementation of the Commandable interface
func (rm *runMigrations) Handle(cmd *cobra.Command, args []string) {
	rp.BootResource("mysql")

	db := rp.MySQL().LogMode(true)

	models := map[string]models.Modelable{}

	for _, table := range rm.tables {
		models[table] = tableRegister[table]
	}

	if len(models) == 0 {
		models = tableRegister
	}

	for _, model := range models {
		if rm.action == "drop" {
			if len(model.Indices()) == 0 {
				for index, _ := range model.Indices() {
					db.Model(model).RemoveIndex(index)
				}
			}

			db.DropTableIfExists(model.Table())

			continue
		}

		if db.HasTable(model.Table()) {
			db.AutoMigrate(model)

			continue
		}

		db.CreateTable(model)

		if len(model.Indices()) == 0 {
			continue
		}

		for index, properties := range model.Indices() {
			if strings.Contains(index, "unique_") {
				db.Model(model).AddUniqueIndex(index, properties...)

				continue
			}

			db.Model(model).AddIndex(index, properties...)
		}
	}
}

// Validate implementation of the Commandable interface
func (rm *runMigrations) Validate(cmd *cobra.Command, args []string) error {
	rm.action = stringFromFlags(cmd.Flags(), "action")
	rm.tables = stringSliceFromFlags(cmd.Flags(), "tables")

	if !inStringSlice(rm.action, []string{"migrate", "drop"}) {
		return errors.New("invalid action, please choose either drop or migrate")
	}

	for _, table := range rm.tables {
		if len(table) == 0 {
			break
		}

		if _, valid := tableRegister[table]; !valid {
			return errors.New("invalid table specified")
		}
	}

	return nil
}
