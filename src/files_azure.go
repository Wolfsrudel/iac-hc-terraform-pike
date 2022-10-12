package pike

import (
	_ "embed" // required for embed
)

//go:embed mapping/azure/resource/resourcegroups/azurerm_resource_group.json
var azurermResourceGroup []byte

//go:embed mapping/azure/resource/serverfarms/azurerm_service_plan.json
var azurermServicePlan []byte

//go:embed mapping/azure/resource/keyvault/azurerm_key_vault.json
var azurermKeyVault []byte

//go:embed mapping/azure/resource/documentdb/azurerm_cosmosdb_account.json
var azureCosmosdbAccount []byte

//go:embed mapping/azure/resource/documentdb/azurerm_cosmosdb_table.json
var azureCosmosdbTable []byte