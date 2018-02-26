//-- Package Declaration -----------------------------------------------------------------------------------------------
package utils

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/JustonDavies/go_domain_model/helpers"
)

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Singleton Functions -----------------------------------------------------------------------------------------------
func InitHelper() (*helpers.DatabaseHelper) {
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
	helper.MigrateAll()

	TruncateModels(helper)
}

func ResetDatabase() {
	var helper = InitHelper()
	helper.MigrateAll()

	DropModels(helper)
	helper.MigrateAll()
}

func TruncateModels(helper *helpers.DatabaseHelper) {
	var transaction = helper.Database().Begin()

	for _, model := range helper.Models() {
		if helper.Database().HasTable(model) {
			var query = fmt.Sprintf("TRUNCATE %s CASCADE;", helper.Database().NewScope(model).TableName())
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

func DropModels(helper *helpers.DatabaseHelper) {
	var transaction = helper.Database().Begin()

	for _, model := range helper.Models() {
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

func SnapshotModels(helper *helpers.DatabaseHelper) {
	var revision = time.Now().UnixNano()
	var transaction = helper.Database().Begin()

	for _, model := range helper.Models() {
		log.Printf("Snapshoting `%s` (version %d)...", reflect.TypeOf(model).Name(), revision)

		//-- Nice Aliases ----------
		var tableName = helper.Database().NewScope(model).TableName() //TODO: There has to be a less dumb way to do this
		var newTableName = fmt.Sprintf("snapshot_%s_%d", tableName, revision)

		//-- Existence check ----------
		if helper.Database().HasTable(model) == false {
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
