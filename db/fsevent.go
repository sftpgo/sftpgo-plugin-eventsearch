package db

// FsEvent defines a filesystem event
type FsEvent struct {
	ID                string `json:"id" gorm:"primaryKey"`
	Timestamp         int64  `json:"timestamp"`
	Action            string `json:"action"`
	Username          string `json:"username"`
	FsPath            string `json:"fs_path"`
	FsTargetPath      string `json:"fs_target_path,omitempty"`
	VirtualPath       string `json:"virtual_path"`
	VirtualTargetPath string `json:"virtual_target_path,omitempty"`
	SSHCmd            string `json:"ssh_cmd,omitempty"`
	FileSize          int64  `json:"file_size,omitempty"`
	Status            int    `json:"status"`
	Protocol          string `json:"protocol"`
	IP                string `json:"ip,omitempty"`
	SessionID         string `json:"session_id"`
	InstanceID        string `json:"instance_id,omitempty"`
}

// TableName defines the database table name
func (ev *FsEvent) TableName() string {
	return "eventstore_fs_events"
}
