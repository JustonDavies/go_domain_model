//-- Package Declaration -----------------------------------------------------------------------------------------------
package provider

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/JustonDavies/go_domain_model/models"

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

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestCreate(t *testing.T) {
	var helper = utils.InitHelper()
	defer helper.Disconnect()

	var provider = models.Provider{
		Name:   fmt.Sprintf("thingme (%d)", time.Now().Unix()),
		Domain: fmt.Sprintf("%d.domain.com", time.Now().Unix()),
	}

	var result = helper.Database().Create(&provider)

	if len(result.GetErrors()) > 0 {
		t.Error("Unable to create Provider: ", result.GetErrors())
	}
}

//TODO: Validation tests
