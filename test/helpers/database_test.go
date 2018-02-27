//-- Package Declaration -----------------------------------------------------------------------------------------------
package database

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"os"
	`testing`

	`github.com/JustonDavies/go_domain_model/helpers/database_helper`
	`github.com/JustonDavies/go_domain_model/test`
)

//-- Decorators --------------------------------------------------------------------------------------------------------
func init() {
	utility.ResetModels()
}

func setup() {}

func teardown() {}

func TestMain(test_function *testing.M) {
	setup()
	var code = test_function.Run()
	teardown()
	os.Exit(code)
}

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestNew(test *testing.T) {
	var database = databaseHelper.New()

	if database == nil {
		test.Error("Unable to intialize new Connection: ", database)
	}
}

// NOTE: Rolled into New() function
//func TestConnect(test *testing.T) {
//	var database = databaseHelper.New()
//	database.Connect()
//}

func TestDisconnect(test *testing.T) {
	var database = databaseHelper.New()

	var err = database.Disconnect()

	if err != nil {
		test.Error("Unable to disconnect to the database: ", err)
	}

}

//NOTE: This test tends to break all other tests because it messes with Postges in a way I don't understand
//func TestMigrate(test *testing.T) {
//	var database = utils.NewDatabaseHelper()
//	defer database.Disconnect()
//
//	utils.DropModels(database)
//
//	for _, model := range database.Models() {
//		if database.Connection().HasTable(model) {
//			test.Error("Table should not yet exist for model: ", reflect.TypeOf(model).Name())
//		}
//	}
//
//	database.MigrateAll()
//
//	for _, model := range database.Models() {
//		if !database.Connection().HasTable(model) {
//			test.Error("Table should exist for model: ", reflect.TypeOf(model).Name())
//		}
//	}
//}
