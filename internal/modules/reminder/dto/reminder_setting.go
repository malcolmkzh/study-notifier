package dto

type UpdateReminderSettingRequest struct {
	Enabled bool `json:"enabled"`
}

type ReminderSettingResponse struct {
	UserID   string `json:"user_id"`
	Enabled  bool   `json:"enabled"`
	Timezone string `json:"timezone"`
}
