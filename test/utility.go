//-- Package Declaration -----------------------------------------------------------------------------------------------
package utility

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"fmt"
	"log"
	"reflect"
	"time"
	
	"github.com/JustonDavies/go_domain_model/helpers/database_helper"
)

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Singleton Functions -----------------------------------------------------------------------------------------------
func timestamp() int64{
	return time.Now().Unix()
}

func ResetModels() {
	var database = databaseHelper.New()
	defer database.Disconnect()
	
	database.MigrateAll()
	truncateModels(database)
}

func ResetDatabase() {
	var database = databaseHelper.New()
	defer database.Disconnect()

	database.MigrateAll()

	dropModels(database)
	database.MigrateAll()
}

func truncateModels(database *databaseHelper.Helper) {
	var transaction = database.Connection().Begin()

	for _, model := range database.Models() {
		if database.Connection().HasTable(model) {
			var query = fmt.Sprintf("TRUNCATE %s CASCADE;", database.Connection().NewScope(model).TableName())
			var result = transaction.Exec(query)
			if result.Error != nil {
				log.Printf("\t Error truncating %s: %s", query, result.Error)
			}
		}
	}

	var result = transaction.Commit()
	if result.Error != nil {
		log.Printf("\t Error commiting transation for all model drop: %s", result.Error)
	}
}

func dropModels(database *databaseHelper.Helper) {
	var transaction = database.Connection().Begin()

	for _, model := range database.Models() {
		var result = transaction.DropTableIfExists(model)
		if result.Error != nil {
			log.Printf("\t Error dropping %s: %s", reflect.TypeOf(model).Name(), result.Error)
		}
	}

	var result = transaction.Commit()
	if result.Error != nil {
		log.Printf("\t Error commiting transation for all model drop: %s", result.Error)
	}
}

func snapshotModels(database *databaseHelper.Helper) {
	var revision = time.Now().UnixNano()
	var transaction = database.Connection().Begin()

	for _, model := range database.Models() {
		log.Printf("Snapshoting `%s` (version %d)...", reflect.TypeOf(model).Name(), revision)

		//-- Nice Aliases ----------
		var tableName = database.Connection().NewScope(model).TableName() //TODO: There has to be a less dumb way to do this
		var newTableName = fmt.Sprintf("snapshot_%s_%d", tableName, revision)

		//-- Existence check ----------
		if database.Connection().HasTable(model) == false {
			log.Printf("\t Error origin table %s does not exist!", tableName)
		}

		//-- Create and Copy to destination table ----------
		var sql = fmt.Sprintf("CREATE TABLE %s AS SELECT * FROM %s;", newTableName, tableName)
		var result = transaction.Exec(sql)
		if result.Error != nil {
			log.Printf("\t Error copying into snapshot table: %s", result.Error)
		}
	}

	var result = transaction.Commit()
	if result.Error != nil {
		log.Printf("\t Error commiting transation for all model snapshot: %s", result.Error)
	}
}
