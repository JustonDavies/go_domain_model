package main

import (
	"encoding/json"
	"log"

	"github.com/JustonDavies/go_domain_model/helpers"
	"github.com/JustonDavies/go_domain_model/models"
	"github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	var helper = helpers.New()
	helper.Connect()
	helper.MigrateAll()

	defer helper.Disconnect()

	//Create
	{
		var category = models.Category{
			Name:     "New Category",
			Metadata: postgres.Jsonb{json.RawMessage(`{"tags": ["tag1", "tag2"]}`)},
		}

		var result = helper.Database().Create(&category)
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

		var result = helper.Database().Create(&category)
		var test = category.Tags()
		log.Print(test)
		if len(result.GetErrors()) > 0 {
			log.Println("Unable to create Category: ", result.GetErrors())
		}
	}

	// Read - JSON Output
	{
		var category models.Category
		var result = helper.Database().Where("id = ?", 1).First(&category)
		//helper.Database().Where(&models.Category{id: 1}).First(&category)

		if len(result.GetErrors()) > 0 {
			log.Println("Unable to find Category: ", result.GetErrors())
		}
		
		var output, err = json.Marshal(category)
		log.Println(err, string(output))
	}

	// Update
	{
		var category models.Category
		helper.Database().Where("id = ?", 1).First(&category)

		category.Name = "Edited Category"
		category.Metadata = postgres.Jsonb{json.RawMessage(`{"tags": ["tag1", "tag2", "tag3"]}`)}
		helper.Database().Save(&category)
	}

	//Update - Single Value
	{
		var category models.Category
		helper.Database().Where("id = ?", 1).First(&category)

		helper.Database().Model(&category).Update("name", "Edited (again) Category")
	}

	// Delete - delete
	{
		var category models.Category
		helper.Database().Where("id = ?", 1).First(&category)

		helper.Database().Delete(&category)
	}
}
