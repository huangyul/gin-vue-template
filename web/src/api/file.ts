import type { ListResponse } from "@/types/api";
import type { SelectOptions } from "@/types/common";
import type { File, FileListQueryParam } from "@/types/file";
import { apiGet, apiPost } from "@/utils/request";

export const getOptions = () => {
  return apiGet<{ user: SelectOptions }>("/file/get-option");
};

export const getList = (param: FileListQueryParam) => {
  return apiPost<ListResponse<File[]>>("/file/list", param);
};

export const uploadFile = (param: FormData) => {
  return apiPost<File>("/file/upload", param, {
    headers: {
      "Content-Type": "multipart/form-data"
    }
  });
};

export const deleteFileById = (id: number) => {
  return apiGet<File>(`/file/delete/${id}`);
};
