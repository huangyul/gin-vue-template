export interface ListResponse<T> {
  data: T;
  total: number;
}

export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}
