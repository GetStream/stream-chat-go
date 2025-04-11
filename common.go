package stream_chat

type PagerRequest struct {
	Limit *int    `json:"limit" validate:"omitempty,gte=0,lte=100"`
	Next  *string `json:"next"`
	Prev  *string `json:"prev"`
}
type SortParamRequestList []*SortParamRequest

type SortParamRequest struct {
	// Name of field to sort by
	Field string `json:"field"`

	// Direction is the sorting direction, 1 for Ascending, -1 for Descending, default is 1
	Direction int `json:"direction"`
}

type PagerResponse struct {
	Next *string `json:"next,omitempty"`
	Prev *string `json:"prev,omitempty"`
}
