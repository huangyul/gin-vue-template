import axios from "axios";
import { ElMessage, ElLoading } from "element-plus";
import router from "@/router";
import { getToken, removeToken } from "./auth";

/*
 * 创建实例
 * 与后端服务通信
 */
const http = axios.create({
  // baseURL: import.meta.env.VITE_BASE_URL,
  baseURL: "/api",
  timeout: 30000
});

let loadingInstance: any;

// 请求次数
let requestNum = 0;

function handleLoading() {
  if (requestNum > 0) {
    requestNum--;
  }
  if (requestNum == 0) {
    loadingInstance.close();
  }
}

/**
 * 请求拦截器
 * 功能：配置请求头
 */
http.interceptors.request.use(
  config => {
    if (requestNum == 0) {
      loadingInstance = ElLoading.service({
        lock: true,
        text: "处理中...",
        background: "rgba(255,255,255,0.3)"
      });
    }
    requestNum++;
    config.headers.Authorization = `Bearer ${getToken().accessToken || ""}`;
    return config;
  },
  error => {
    console.error("网络错误，请稍后重试");
    handleLoading();
    return Promise.reject(error);
  }
);

/**
 * 响应拦截器
 * 功能：处理异常
 */
http.interceptors.response.use(
  response => {
    console.log("response: ", response);

    handleLoading();
    if (response.status != 200) {
      if (response.data.exmsg) {
        ElMessage.error(response.data.exmsg);
      } else {
        ElMessage.error("服务器出错");
      }
      return Promise.reject();
    } else if (response.config.headers.download == "download") {
      // 下载接口
      return response.data;
    } else if (response.data.code != 0) {
      ElMessage.error(response.data.exmsg);
      return Promise.reject();
    } else {
      return response.data;
    }
  },
  error => {
    console.log("error: ", error);

    handleLoading();
    if (error.response?.status == 401) {
      removeToken();
      ElMessage.error("登录已过期，请重新登录");
      router.replace("/login");
      return;
    }
    if (error.response?.status == 404) {
      ElMessage.error("接口不存在");
      return Promise.reject(error);
    }
    if (error.response?.data) {
      ElMessage.error(error.response.data.exmsg);
    } else {
      ElMessage.error("发生未知错误");
    }
    return Promise.reject(error);
  }
);

export async function apiGet<T>(url: string, params?: any): Promise<T> {
  const response = await http.get<T>(url, { params });
  return response.data;
}

export async function apiPost<T>(
  url: string,
  data?: any,
  config = {}
): Promise<T> {
  const response = await http.post<T>(url, data, { ...config });
  return response.data;
}

export async function apiPut<T>(url: string, data?: any): Promise<T> {
  const response = await http.put<T>(url, data);
  return response.data;
}

export async function apiDel<T>(url: string, params?: any): Promise<T> {
  const response = await http.delete<T>(url, { params });
  return response.data;
}

export async function downloadFile<T>(url: string, data: any) {
  const response = await http.post<T>(url, data, {
    responseType: "blob",
    headers: { download: "download" }
  });
  return response as T;
}

export default http;
