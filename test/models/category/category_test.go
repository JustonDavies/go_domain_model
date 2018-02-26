//-- Package Declaration -----------------------------------------------------------------------------------------------
package category

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/JustonDavies/go_domain_model/models"

	`github.com/JustonDavies/go_domain_model/test`
	"github.com/jinzhu/gorm/dialects/postgres"
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

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestCreate(test *testing.T) {
	//-- Setup model helper ----------
	var helper = utils.InitHelper()
	defer helper.Disconnect()

	//-- Shared Variables ----------
	var name = `Test Create`
	var metadata = `{"tags":["tag2","tag3","tag4"]}`

	//-- Pre-conditions ----------
	{
		var count int
		helper.
			Database().
			Model(&models.Category{}).
			Where("name = ?", name).
			Count(&count)

		if  count != 0 {
			test.Error("Pre-conditions: A record with these parameters already exists")
		}
	}

	//-- Action ----------
	{
		var category = models.Category {
			Name: name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		var result = helper.Database().Create(&category)
		if len(result.GetErrors()) > 0 {
			test.Error("Unable to create: ", result.GetErrors())
		}

	}

	//-- Post-conditions ----------
	{
		var count int
		helper.
			Database().
			Model(&models.Category{}).
			Where("name = ?", name).
			Count(&count)

		if count != 1 {
			test.Error("Post-conditions: Unable to find persisted record")
		}
	}
}

func TestCreateFromJSON(test *testing.T) {
	//-- Setup model helper ----------
	var helper = utils.InitHelper()
	defer helper.Disconnect()

	//-- Shared Variables ----------
	var name = `Test Create From JSON`
	var metadata = `{"tags":["tag2","tag3","tag4"]}`

	//-- Pre-conditions ----------
	{
		var count int
		helper.
			Database().
			Model(&models.Category{}).
				Where("name = ?", name).
					Count(&count)

		if  count != 0 {
			test.Error("Pre-conditions: A record with these parameters already exists")
		}
	}
	//-- Action ----------
	{
		var category models.Category

		var err = json.Unmarshal([]byte(fmt.Sprintf(`{"Name":"%s","Metadata":%s}`, name, metadata)), &category)
		if err != nil {
			test.Error("Unable to unmarshal data: ", err)
		}

		var result = helper.Database().Create(&category)
		if len(result.GetErrors()) > 0 {
			test.Error("Unable to create: ", result.GetErrors())
		}
		
	}
	//-- Post-conditions ----------
	{
		var count int
		helper.
			Database().
			Model(&models.Category{}).
			Where("name = ?", name).
			Count(&count)

		if count != 1 {
			test.Error("Post-conditions: Unable to find persisted record")
		}
	}
}

//TODO: Validation tests
