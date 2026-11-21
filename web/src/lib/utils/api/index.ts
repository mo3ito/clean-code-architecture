import type { AxiosRequestConfig } from "axios";
import { api } from "./client";

export const http = {
  get: <T>(url: string, config?: AxiosRequestConfig) => api.get<T>(url, config),

  post: <T, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig) =>
    api.post<T>(url, data, config),

  put: <T, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig) =>
    api.put<T>(url, data, config),

  patch: <T, D = unknown>(url: string, data?: D, config?: AxiosRequestConfig) =>
    api.patch<T>(url, data, config),

  delete: <T>(url: string, config?: AxiosRequestConfig) =>
    api.delete<T>(url, config),
};

export { api };
