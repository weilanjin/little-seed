package dto

type LogSearchReq struct {
	FileName string `json:"file_name"`
	Include  string `json:"include"`
	Exclude  string `json:"exclude"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
}
