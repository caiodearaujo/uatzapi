package data

// TablesPostgres returns a slice of all the models that are associated with PostgreSQL tables.
// This function is used to register these models with the ORM (Bun) for creating, migrating, or interacting with the database tables.
func TablesPostgres() []interface{} {
	return []interface{}{
		(*Device)(nil),         // Device model representing the devices connected to WhatsApp.
		(*DeviceHandler)(nil),  // DeviceHandler model for managing device handler states.
		(*DeviceWebhook)(nil),  // DeviceWebhook model for storing webhook configurations.
		(*WebhookMessage)(nil), // WebhookMessage model for managing messages sent via webhooks.
	}
}
