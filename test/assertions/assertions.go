//-- Package Declaration -----------------------------------------------------------------------------------------------
package assertions

import (
	"fmt"
	"testing"
	
	"github.com/jinzhu/gorm"
	"github.com/qor/validations"
)

//-- Structs -----------------------------------------------------------------------------------------------------------

//-- Singleton Functions -----------------------------------------------------------------------------------------------

func TransactionError(result *gorm.DB, test *testing.T) {
	if len(result.GetErrors()) < 1 {
		test.Error(fmt.Sprintf(`Transaction: Error was not found when one was expected: %s`, result.GetErrors()))
	}
}

func NotTransactionError(result *gorm.DB, test *testing.T) {
	if len(result.GetErrors()) > 0 {
		test.Error(fmt.Sprintf(`Transaction: Error was found when one was NOT expected: %s`, result.GetErrors()))
	}
}

func ValidationError(column string, result *gorm.DB, test *testing.T) {
	var present = validationErrorExists(column, result)

	if present == false {
		test.Error(fmt.Sprintf(`Validation: Error was not found for '%s' when one was expected`, column))
	}
}

func NotValidationError(column string, result *gorm.DB, test *testing.T) {
	var present = validationErrorExists(column, result)

	if present == true {
		test.Error(fmt.Sprintf(`Validation: Error was found for '%s' when one NOT was expected`, column))
	}
}

func RecordExists(id uint, model interface{}, database *gorm.DB, test *testing.T) {
	var count = modelCount(id, model, database)

	if count != 1 {
		test.Error(`Existance: Unable to find persisted record when one was expected`)
	}
}

func NotRecordExists(id uint, model interface{}, database *gorm.DB, test *testing.T) {
	var count = modelCount(id, model, database)

	if count != 0 {
		test.Error(`Existance: Able to find persisted record when one was NOT expected`)
	}
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
func validationErrorExists(column string, result *gorm.DB) bool{
	var present = false
	for _, err := range result.GetErrors() {
		switch err.(type) {
		case *validations.Error:
			if column == err.(*validations.Error).Column {
				present = true
			}
		}
		if present == true {
			break
		}
	}

	return present
}

func modelCount(id uint, model interface{}, database *gorm.DB,) int{
	var count int
	database.Where(`id = ?`, id).Find(model).Count(&count)

	return count
}