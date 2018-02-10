//-- Package Declaration -----------------------------------------------------------------------------------------------
package utils

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/JustonDavies/go_domain_model/helpers"
)

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Singleton Functions -----------------------------------------------------------------------------------------------
func InitHelper()(*helpers.DatabaseHelper){
	var helper = helpers.New()
	var err = helper.Connect()
	helper.MigrateAll()

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}

	return helper
}

func ResetModels() {
	var helper = InitHelper()
	defer helper.Disconnect()

	//SnapshotModels(helper)
	DropModels(helper)
	helper.MigrateAll()
}

func DropModels(helper *helpers.DatabaseHelper) {
	for _, model := range helper.Models() {
		log.Printf("Dropping `%s`...", reflect.TypeOf(model).Name())

		//-- Nice Aliases ----------
		var result *gorm.DB

		result = helper.Database().DropTableIfExists(model)
		if result.Error != nil {
			log.Printf("\t Error: %s", result.Error)
		}
	}
}

func SnapshotModels(helper *helpers.DatabaseHelper)  {
	var revision = time.Now().UnixNano()

	for _, model := range helper.Models() {
		log.Printf("Snapshoting `%s` (version %d)...", reflect.TypeOf(model).Name(), revision)

		//-- Nice Aliases ----------
		var result *gorm.DB
		var tableName = helper.Database().NewScope(model).TableName() //TODO: There has to be a less dumb way to do this
		var newTableName = fmt.Sprintf("snapshot_%s_%d", tableName, revision)

		//-- Existence check ----------
		if helper.Database().HasTable(model) == false {
			log.Printf("\t Error origin table %s does not exist!", tableName)
		}

		//-- Create and Copy to destination table ----------
		var sql = fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM %s;", newTableName, tableName)
		result = helper.Database().Exec(sql)
		if result.Error != nil {
			log.Printf("\t Error copying into snapshot table: %s", result.Error)
		}
	}
}