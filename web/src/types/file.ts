export interface FileListQueryParam {
  page: number;
  page_size: number;
  file_name: string;
  user_id: string;
}

export interface File {
  id: number;
  file_name: string;
  link: string;
  upload_user: number;
  upload_time: string;
}
