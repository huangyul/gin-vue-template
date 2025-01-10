export interface User {
  id: number;
  username: string;
  nickname: string;
  created_at: string;
}

export interface UserListQueryParam {
  page: number;
  page_size: number;
  username: string;
  nickname: string;
}
