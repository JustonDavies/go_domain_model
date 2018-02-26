//-- Package Declaration -----------------------------------------------------------------------------------------------
package helpers

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	`log`
	`os`
	`reflect`

	`github.com/jinzhu/gorm`

	`github.com/qor/validations`

	`github.com/JustonDavies/go_domain_model/models`
)

//-- Structs -----------------------------------------------------------------------------------------------------------
type DatabaseHelper struct {
	database *gorm.DB
	models   []interface{}
}

//-- Singleton Functions -----------------------------------------------------------------------------------------------
func New() (*DatabaseHelper) {
	var helper = new(DatabaseHelper)
	helper.models = []interface{}{
		models.Provider{},
		models.Category{},
		models.Service{},
	}

	return helper
}

//-- Public Functions --------------------------------------------------------------------------------------------------
func (helper *DatabaseHelper) Database() (*gorm.DB) {
	return helper.database
}

func (helper *DatabaseHelper) Models() ([]interface{}) {
	return helper.models
}

func (helper *DatabaseHelper) Connect() (error) {
	var databaseDriver = os.Getenv("MODEL_DATABASE_DRIVER")
	var databaseParams = os.Getenv("MODEL_DATABASE_PARAMETERS")

	var database, err = gorm.Open(databaseDriver, databaseParams)
	if err != nil {
		return err
	}
	helper.database = database

	validations.RegisterCallbacks(helper.database)
	return nil
}

func (helper *DatabaseHelper) Disconnect() (error) {
	var err = helper.database.Close()
	if err != nil {
		return err
	}
	return nil
}

func (helper *DatabaseHelper) MigrateAll() {
	var transaction = helper.Database().Begin()

	for _, model := range helper.models {
		var result = transaction.AutoMigrate(model)
		if result.Error != nil {
			log.Printf("\t Error migrating %s: %s", reflect.TypeOf(model).Name(), result.Error)
		}
	}

	var result = transaction.Commit()
	if result.Error != nil {
		log.Printf("\t Error commiting transation for all model migration: %s", result.Error)
	}
}

//-- Private Functions -------------------------------------------------------------------------------------------------
