//-- Package Declaration -----------------------------------------------------------------------------------------------
package models

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	`encoding/json`
	`errors`
	`fmt`
	`log`
	`regexp`

	`github.com/jinzhu/gorm`
	`github.com/jinzhu/gorm/dialects/postgres`
	`github.com/qor/validations`
)

//-- Structs Declaration -----------------------------------------------------------------------------------------------
type Category struct {
	//-- User Variables ----------
	Name     string `gorm:"not null;unique"valid:"required"`
	Metadata postgres.Jsonb //`valid:"required"`

	//-- System Variables ----------

	//-- Relations ----------
	//Domains []Domain

	//-- Automated fields (ID, Timestamps) ----------
	gorm.Model
}

type catalogMetadata struct {
	Tags []string
}

//-- Exported Functions ------------------------------------------------------------------------------------------------
func (category *Category) ParsedMetadata() (catalogMetadata) {
	var metadata catalogMetadata

	if err := json.Unmarshal(category.Metadata.RawMessage, &metadata); err != nil {
		log.Println(`Tags: Unable to unmarshal Metadata, returning empty struct`)
	}

	return metadata
}

func (category *Category) Tags() ([]string) {
	return category.ParsedMetadata().Tags
}

//-- ORM Functions -----------------------------------------------------------------------------------------------------
func (category *Category) BeforeSave() (error) {
	return nil
}

func (category *Category) Sanitize(database *gorm.DB) {
	category.sanitizeMetaData(database)
}

func (category *Category) Validate(database *gorm.DB) {
	category.validateMetadata(database)
}

//-- Internal Functions ------------------------------------------------------------------------------------------------
func (category *Category) sanitizeMetaData(database *gorm.DB) {
	var sanitized []byte
	var metadata catalogMetadata

	// 1.) Read existing data into a known format
	//NOTE: Incredibly, Postgres and GORM allow any string through
	{
		json.Unmarshal([]byte(category.Metadata.RawMessage), &metadata)
		//if json.Unmarshal([]byte(category.Metadata.RawMessage), &metadata) != nil {
		//	log.Println(`Sanitizing: Unable to unmarshal Metadata, an empty struct will be used`)
		//}
	}

	// 2.) Populate defaults if nothing is present
	{
		if metadata.Tags == nil {
			metadata.Tags = []string{}
		}
	}

	// 3.) Write sanitized data back into the structure
	{
		var err error
		if sanitized, err = json.Marshal(metadata); err != nil {
			database.AddError(
				errors.New(
					fmt.Sprintf(`Sanitizing:  Unable to marshal clean Metadata, something terrible has happened: %s`, err),
				),
			)
		}
		category.Metadata = postgres.Jsonb{sanitized}
	}
}

func (category *Category) validateMetadata(database *gorm.DB) {
	var metadata catalogMetadata

	// 1.) Read existing data into a known format
	//NOTE: Should be redundant, but we have to unmarshal it anyway
	{
		if json.Unmarshal([]byte(category.Metadata.RawMessage), &metadata) != nil {
			database.AddError(
				validations.NewError(
					category,
					`Metadata`,
					`Validation: Unable to unmarshal Metadata`,
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
						`Tag`,
						fmt.Sprintf(`Validation: Tag '%s' need to follow format of 'tag_name'`, tag),
					),
				)
			}
		}
	}

}
