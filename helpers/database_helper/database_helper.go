//-- Package Declaration -----------------------------------------------------------------------------------------------
package databaseHelper

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"fmt"
	`log`
	`os`
	`reflect`

	`github.com/jinzhu/gorm`

	`github.com/qor/validations`

	`github.com/JustonDavies/go_domain_model/models`
	`github.com/JustonDavies/go_domain_model/utilities/sanitizer`
)

//-- Structs Declaration -----------------------------------------------------------------------------------------------
type Helper struct {
	connection *gorm.DB
	models     []interface{}
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func New() (*Helper) {
	// 1.) Create container
	var helper = new(Helper)

	// 2.) Import models
	helper.models = []interface{}{
		models.Domain{},
		models.Category{},
	}

	// 3.) Connect
	{
		var databaseDriver = os.Getenv("MODEL_DATABASE_DRIVER")
		var databaseParams = os.Getenv("MODEL_DATABASE_PARAMETERS")

		var database, err = gorm.Open(databaseDriver, databaseParams)
		if err != nil {
			panic(fmt.Sprintf(`Unable to connect to database: %s`, err))
		}

		helper.connection = database
	}

	// 4.) Register enhances callbacks
	sanitizer.RegisterCallbacks(helper.connection)
	validations.RegisterCallbacks(helper.connection)

	//5.)
	return helper
}

func (helper *Helper) Connection() (*gorm.DB) {
	return helper.connection
}

func (helper *Helper) Models() ([]interface{}) {
	return helper.models
}

func (helper *Helper) Disconnect() (error) {
	var err = helper.connection.Close()
	if err != nil {
		return err
	}
	return nil
}

func (helper *Helper) MigrateAll() {
	var transaction = helper.Connection().Begin()

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

//-- Internal Functions ------------------------------------------------------------------------------------------------
