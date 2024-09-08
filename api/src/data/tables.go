package data

func TablesPostgres() []interface{} {
	return []interface{}{
		(*DeviceHandler)(nil),
		(*DeviceWebhook)(nil),
		(*WebhookMessage)(nil),
	}
}
