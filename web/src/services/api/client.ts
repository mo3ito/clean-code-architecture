import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

import { API_CONFIG } from "./config";
import { ApiError } from "./api-error";

import type { ApiResponse } from "./api-response";

export const apiClient = axios.create({
  baseURL: API_CONFIG.baseURL,
  timeout: API_CONFIG.timeout,
  withCredentials: true,
  headers: {
    Accept: "application/json",
  },
});

// Request Interceptor
apiClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  return config;
});

// Response Interceptor
apiClient.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse;

    // حتی اگر HTTP Status موفق باشد،
    // success: false را نیز به‌عنوان خطا در نظر می‌گیریم.
    if (body && body.success === false) {
      throw new ApiError(
        body.message || "Request failed",
        response.status,
        body.error,
        body.data,
      );
    }

    return response;
  },

  (error: unknown) => {
    if (axios.isAxiosError<ApiResponse>(error)) {
      const response = error.response;
      const body = response?.data;

      return Promise.reject(
        new ApiError(
          body?.message || error.message || "Network error",
          response?.status ?? 0,
          body?.error,
          body?.data,
        ),
      );
    }

    return Promise.reject(error);
  },
);
