package data

func TablesPostgres() []interface{} {
	return []interface{}{
		(*Device)(nil),
		(*DeviceHandler)(nil),
		(*DeviceWebhook)(nil),
		(*WebhookMessage)(nil),
	}
}
