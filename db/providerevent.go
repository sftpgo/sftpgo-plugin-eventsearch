package db

// ProviderEvent defines a provider event
type ProviderEvent struct {
	ID         string `json:"id" gorm:"primaryKey"`
	Timestamp  int64  `json:"timestamp"`
	Action     string `json:"action"`
	Username   string `json:"username"`
	IP         string `json:"ip,omitempty"`
	ObjectType string `json:"object_type"`
	ObjectName string `json:"object_name"`
	ObjectData []byte `json:"object_data"`
	InstanceID string `json:"instance_id,omitempty"`
}

// TableName defines the database table name
func (ev *ProviderEvent) TableName() string {
	return "eventstore_provider_events"
}
