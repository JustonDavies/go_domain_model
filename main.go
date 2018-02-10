package main

//import (
//	"github.com/JustonDavies/go_domain_model/helpers"
//	"github.com/JustonDavies/go_domain_model/models"
//)

//func main() {
//	var client = database.Initialize()
//	client.Connect()
//	client.MigrateAll()
//
//	defer client.Disconnect()
//
//	// Read
//	var provider models.Provider
//	client.Database().First(&provider, 1)                    // find with id 1
//	client.Database().First(&provider, "name = ?", "Google") // find with name 'Google'
//
//	// Update - update domain to service.google.com
//	client.Database().Model(&provider).Update("Domain", "service.google.com")
//
//	//// Delete - delete
//	client.Database().Delete(&provider)
//}
