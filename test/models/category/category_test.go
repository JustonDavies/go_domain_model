//-- Package Declaration -----------------------------------------------------------------------------------------------
package category

//-- Imports -----------------------------------------------------------------------------------------------------------
import (
	`encoding/json`
	`fmt`
	`os`
	`testing`

	`github.com/JustonDavies/go_domain_model/helpers/database_helper`
	`github.com/JustonDavies/go_domain_model/models`
	`github.com/JustonDavies/go_domain_model/test`
	`github.com/JustonDavies/go_domain_model/test/assertions`

	`github.com/jinzhu/gorm`
	`github.com/jinzhu/gorm/dialects/postgres`
)

//-- Decorators --------------------------------------------------------------------------------------------------------
func init() {
	//Run once, once testing suite starts
	utility.ResetModels()
}

func setup() {
	//Run before every test
}

func teardown() {
	//Run after every test
}

func TestMain(test_function *testing.M) {
	setup()
	var code = test_function.Run()
	teardown()
	os.Exit(code)
}

//-- Local Constants ---------------------------------------------------------------------------------------------------

//-- Tests -------------------------------------------------------------------------------------------------------------
func TestCreate(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test create`
	var metadata = `{"tags":["tag_one","tag_two","tag_three"]}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)
	}
}

func TestCreateFromJSON(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test create from JSON`
	var metadata = `{"tags":["tag_one","tag_two","tag_three"]}`
	var jsonInput = fmt.Sprintf(`{"Name":"%s","Metadata":%s}`, name, metadata)

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		var err = json.Unmarshal([]byte(jsonInput), &category)
		if err != nil {
			test.Error("TEST: Unable to unmarshal data: ", err)
		}

		result = database.Connection().Create(&category)
	}
	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)
	}
}

func TestValidName(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test valid name`
	var metadata = `{"tags":["tag_one","tag_two","tag_three"]}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.NotValidationError(`Name`, result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)
	}
}

func TestNilName(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Shared Variables ----------
	var name string
	var metadata = `{"tags":["tag_one","tag_two","tag_three"]}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.TransactionError(result, test)
		assertions.ValidationError(`Name`, result, test)
		assertions.NotRecordExists(category.ID, &category, database.Connection(), test)
	}
}

func TestValidMetadata(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test valid metadata`
	var metadata = `{"tags":["tag_one","tag_two","tag_three"]}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.NotValidationError(`Metadata`, result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)
	}
}

func TestEmptyMetadata(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test empty metadata`
	var metadata = `null`
	//var jsonInput = fmt.Sprintf(`{"Name":"%s","Metadata":null}`, name)

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{

		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.NotValidationError(`Metadata`, result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)

		var storedMetadata, _ = json.Marshal(category.Metadata)
		if metadata == string(storedMetadata) {
			test.Error(fmt.Sprintf("Invalid metadata was not properly sanitized and persisted: %s -> %s", metadata, storedMetadata))
		}
	}
}

func TestMalformedMetadata(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test malformed metadata`
	var metadata = `{"malformed":"Sanitize me!"}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.NotValidationError(`Metadata`, result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)

		var storedMetadata, _ = json.Marshal(category.Metadata)
		if metadata == string(storedMetadata) {
			test.Error(fmt.Sprintf("Invalid metadata was not properly sanitized and persisted: %s -> %s", metadata, storedMetadata))
		}
	}
}

func TestInvalidMetadata(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test invalid metadata`
	var metadata = `"bananas"`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.NotValidationError(`Metadata`, result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)

		var storedMetadata, _ = json.Marshal(category.Metadata)
		if metadata == string(storedMetadata) {
			test.Error(fmt.Sprintf("Invalid metadata was not properly sanitized and persisted: %s -> %s", metadata, storedMetadata))
		}

	}
}

func TestValidTags(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test valid tags`
	var metadata = `{"tags":["tag_one","tag_two","tag_three"]}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.NotTransactionError(result, test)
		assertions.NotValidationError(`Tag`, result, test)
		assertions.RecordExists(category.ID, &category, database.Connection(), test)
	}
}

func TestInvalidTags(test *testing.T) {
	//-- Setup model database ----------
	var database = databaseHelper.New()
	defer database.Disconnect()

	//-- Shared Variables ----------
	var result *gorm.DB
	var category models.Category

	//-- Test Parameters ----------
	var name = `Test valid tags`
	var metadata = `{"tags":["tag one","tag_two","tag_three"]}`

	//-- Pre-conditions ----------
	{}

	//-- Action ----------
	{
		category = models.Category{
			Name:     name,
			Metadata: postgres.Jsonb{json.RawMessage(metadata)},
		}

		result = database.Connection().Create(&category)
	}

	//-- Post-conditions ----------
	{
		assertions.TransactionError(result, test)
		assertions.ValidationError(`Tag`, result, test)
		assertions.NotRecordExists(category.ID, &category, database.Connection(), test)
	}
}