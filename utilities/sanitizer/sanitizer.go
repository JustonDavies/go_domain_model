//-- Package Declaration -----------------------------------------------------------------------------------------------
package sanitizer

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	"github.com/jinzhu/gorm"
)

//-- Structs Declaration -----------------------------------------------------------------------------------------------

//-- Exported Functions ------------------------------------------------------------------------------------------------
func RegisterCallbacks(database *gorm.DB) {
	var callback = database.Callback()

	if callback.Create().Get("sanitizer:sanitize") == nil {
		callback.Create().Before("gorm:before_create").Register("sanitizer:sanitize", sanitize)
	}
	if callback.Update().Get("sanitizer:sanitize") == nil {
		callback.Update().Before("gorm:before_update").Register("sanitizer:sanitize", sanitize)
	}
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
func sanitize(scope *gorm.Scope) {
	scope.CallMethod("Sanitize")
}
