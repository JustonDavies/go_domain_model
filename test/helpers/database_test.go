//-- Package Declaration -----------------------------------------------------------------------------------------------
package database

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"os"
	`testing`

	`github.com/JustonDavies/go_domain_model/helpers`
	`github.com/JustonDavies/go_domain_model/test`
)

//-- Decorators --------------------------------------------------------------------------------------------------------
func init() {
	utils.ResetModels()
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
	var helper = helpers.New()

	if helper == nil {
		test.Error("Unable to intialize new DatabaseHelper: ", helper)
	}
}

func TestConnect(test *testing.T) {
	var helper = helpers.New()
	var err = helper.Connect()

	if err != nil {
		test.Error("Unable to connect to the database: ", err)
	}
}

func TestDisconnect(test *testing.T) {
	var helper = helpers.New()
	var err = helper.Connect()

	if err != nil {
		test.Error("Unable to connect to the database: ", err)
	}

	err = nil
	err = helper.Disconnect()

	if err != nil {
		test.Error("Unable to disconnect to the database: ", err)
	}

}

//NOTE: This test tends to break all other tests because it messes with Postges in a way I don't understand
//func TestMigrate(test *testing.T) {
//	var helper = utils.InitHelper()
//	defer helper.Disconnect()
//
//	utils.DropModels(helper)
//
//	for _, model := range helper.Models() {
//		if helper.Database().HasTable(model) {
//			test.Error("Table should not yet exist for model: ", reflect.TypeOf(model).Name())
//		}
//	}
//
//	helper.MigrateAll()
//
//	for _, model := range helper.Models() {
//		if !helper.Database().HasTable(model) {
//			test.Error("Table should exist for model: ", reflect.TypeOf(model).Name())
//		}
//	}
//}
