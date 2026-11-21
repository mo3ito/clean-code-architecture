
import type { InternalAxiosRequestConfig } from "axios";
import { getAccessToken } from "./auth";
import type { AxiosResponse } from "axios";

export interface CustomAxiosRequestConfig
  extends InternalAxiosRequestConfig {
  requiresAuth?: boolean;
}

export function setupRequestInterceptor(
  config: CustomAxiosRequestConfig
) {
  // اگر درخواست صراحتاً Authentication نخواهد
  if (config.requiresAuth === false) {
    return config;
  }

  // دریافت Access Token
  const token = getAccessToken();

  // اضافه کردن Token به Header
  if (token) {
    config.headers.set(
      "Authorization",
      `Bearer ${token}`
    );
  }

  return config;
}

export function setupResponseInterceptor(
  response: AxiosResponse
) {
  return response;
}

export function setupResponseErrorInterceptor(
  error: unknown
) {
  return Promise.reject(error);
}
