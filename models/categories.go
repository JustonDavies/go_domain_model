//-- Package Declaration -----------------------------------------------------------------------------------------------
package models

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	`encoding/json`
	"fmt"
	`log`
	"regexp"

	`github.com/jinzhu/gorm`
	`github.com/jinzhu/gorm/dialects/postgres`
	`github.com/qor/validations`
)

//-- Structs Declaration -----------------------------------------------------------------------------------------------
type Category struct {
	//-- User Variables ----------
	Name     string `gorm:"not null;unique"`
	Metadata postgres.Jsonb

	//-- System Variables ----------

	//-- Relations ----------
	Services []Service

	//-- Automated fields (ID, Timestamps) ----------
	gorm.Model
}

type catalogMetadata struct {
	Tags []string
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func (category *Category) Tags() ([]string) {
	var prop catalogMetadata

	if err := json.Unmarshal(category.Metadata.RawMessage, &prop); err != nil {
		log.Println(`Unable to unmarshal Metadata, returning empty struct`)
	}

	return prop.Tags
}

//-- ORM Functions -----------------------------------------------------------------------------------------------------
func (category *Category) BeforeSave() (err error) {
	err = category.sanitizeMetaData();
	if err != nil {
		log.Println(`Callbacks: Unable to sanitize Metadata, returning error`)
		return
	}
	return
}
func (category *Category) Validate(database *gorm.DB) {
	category.validateMetadata(database)
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
func (category *Category) sanitizeMetaData() (err error) {
	var sanitized []byte
	var metadata catalogMetadata

	// 1.) Read existing data into a known format
	//NOTE: Incredibly, Postgres and GORM allow any string through
	{
		if json.Unmarshal([]byte(category.Metadata.RawMessage), &metadata) != nil {
			log.Println(`Sanitizing: Unable to unmarshal Metadata, an empty struct will be used`)
		}
	}

	// 2.) Populate defaults if nothing is present
	{
		if metadata.Tags == nil {
			metadata.Tags = []string{}
		}
	}

	// 3.) Write sanitized data back into the structure
	{
		if sanitized, err = json.Marshal(metadata); err != nil {
			log.Println(`Sanitizing: Unable to marshal clean Metadata, returning error`)
			return
		}
		category.Metadata = postgres.Jsonb{sanitized}
	}

	// 4.) Return
	return
}

func (category *Category) validateMetadata(database *gorm.DB) {
	var metadata catalogMetadata

	// 1.) Read existing data into a known format
	//NOTE: Should be redundant, but we have to unmarshal it anyway
	{
		if json.Unmarshal([]byte(category.Metadata.RawMessage), &metadata) != nil  {
			database.AddError(
				validations.NewError(
					category,
					`Metadata`,
					`Unable to unmarshal Metadata`,
				),
			)
		}
	}

	// 2.) Validate Tags follow convention
	{
		var validTag = regexp.MustCompile(`\A[a-z_]+\z`)
		for _, tag := range category.Tags() {
			if !validTag.MatchString(tag) {
				database.AddError(
					validations.NewError(
						category,
						"Tags",
						fmt.Sprintf(`"Tag '%s' need to follow format of 'tag_name'"`, tag),
					),
				)
			}
		}
	}
	
}