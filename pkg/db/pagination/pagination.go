package pagination

import (
	"encoding/base64"
	"encoding/json"
)

type Pagination struct {
	Cursor string `form:"cursor"`
	Limit  int    `form:"limit,default=10" validate:"gte=1,lte=250"` // Min 1, Max 250
}

type Cursor struct {
	CreatedAt string `json:"created_at,omitempty"`
	ID        string `json:"id,omitempty"`
}

type PageInfo struct {
	NextCursor     string `json:"next_cursor"`
	PreviousCursor string `json:"previous_cursor"`
	HasMore        bool   `json:"has_more"`
}

func EncodeCursor(data Cursor) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeCursor(data string) (*Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var cursor Cursor
	if err := json.Unmarshal(b, &cursor); err != nil {
		return nil, err
	}

	return &cursor, nil
}

func BuildCursorPageInfo[T any](data []*T, limit int, extractCursor func(*T) string) (pageInfo *PageInfo) {
	if len(data) > limit {
		return &PageInfo{
			HasMore:    true,
			NextCursor: extractCursor(data[limit-1]),
		}
	}
	return
}
