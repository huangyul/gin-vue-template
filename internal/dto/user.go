package dto

type UserListQueryParams struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type UserResp struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	Username  string `json:"username"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

type UserListResp struct {
	Data  []UserResp `json:"data"`
	Total int64      `json:"total"`
}
