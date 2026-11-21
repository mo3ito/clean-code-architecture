import type {
  AxiosError,
  InternalAxiosRequestConfig,
} from "axios";

import { api } from "./client";

import {
  getAccessToken,
  setAccessToken,
  clearAccessToken,
  refreshAccessToken,
} from "./auth";

function handleRequest(
  config: InternalAxiosRequestConfig
) {
  if (config.requiresAuth === false) {
    return config;
  }

  const token = getAccessToken();

  if (config.requiresAuth === true && !token) {
    return Promise.reject(
      new Error("Authentication required")
    );
  }

  if (token) {
    config.headers.set(
      "Authorization",
      `Bearer ${token}`
    );
  }

  return config;
}

async function handleResponseError(
  error: AxiosError
) {
  const originalRequest = error.config;

  if (
    error.response?.status !== 401 ||
    !originalRequest ||
    originalRequest.requiresAuth === false ||
    originalRequest.skipAuthRefresh ||
    originalRequest._retry
  ) {
    return Promise.reject(error);
  }

  originalRequest._retry = true;

  try {
    const newToken = await refreshAccessToken();

    originalRequest.headers.set(
      "Authorization",
      `Bearer ${newToken}`
    );

    return api(originalRequest);
  } catch (refreshError) {
    clearAccessToken();

    return Promise.reject(refreshError);
  }
}

export function setupInterceptors() {
  api.interceptors.request.use(
    handleRequest,
    (error) => Promise.reject(error)
  );

  api.interceptors.response.use(
    (response) => response,
    handleResponseError
  );
}