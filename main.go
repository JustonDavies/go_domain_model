package main

import (
	`encoding/json`
	`log`

	`github.com/JustonDavies/go_domain_model/helpers/database_helper`
	`github.com/JustonDavies/go_domain_model/models`
	`github.com/jinzhu/gorm/dialects/postgres`
)

func main() {
	var database = databaseHelper.New()
	database.MigrateAll()

	defer database.Disconnect()

	//Create
	{
		var category = models.Category{
			Name:     "",
			Metadata: postgres.Jsonb{json.RawMessage(`{"tags": ["tag_one", "tag_two"]}`)},
		}

		var result = database.Connection().Create(&category)
		if len(result.GetErrors()) > 0 {
			log.Println("Unable to create Provider: ", result.GetErrors())
		}
	}

	//Create from JSON
	{
		var input = `{"Name":"Category From JSON","Metadata":{"tags": null}}`
		var category models.Category

		var err = json.Unmarshal([]byte(input), &category)
		if err != nil {
			log.Println("Unable to unmarshal Category data: ", err)
		}

		var result = database.Connection().Create(&category)
		if len(result.GetErrors()) > 0 {
			log.Println("Unable to create Category: ", result.GetErrors())
		}
	}

	// Read - JSON Output
	{
		var category models.Category
		var result = database.Connection().Where("id = ?", 1).First(&category)
		//database.Connection().Where(&models.Category{id: 1}).First(&category)

		if len(result.GetErrors()) > 0 {
			log.Println("Unable to find Category: ", result.GetErrors())
		}

		var output, err = json.Marshal(category)
		log.Println(err, string(output))
	}

	// Update
	{
		var category models.Category
		database.Connection().Where("id = ?", 1).First(&category)

		category.Name = "Edited Category"
		category.Metadata = postgres.Jsonb{json.RawMessage(`{"tags": ["tag1", "tag2", "tag3"]}`)}
		database.Connection().Save(&category)
	}

	//Update - Single Value
	{
		var category models.Category
		database.Connection().Where("id = ?", 1).First(&category)

		database.Connection().Model(&category).Update("name", "Edited (again) Category")
	}

	// Delete - delete
	{
		var category models.Category
		database.Connection().Where("id = ?", 1).First(&category)

		database.Connection().Delete(&category)
	}
}
