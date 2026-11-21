import axios from "axios";
import { API_CONFIG } from "./config";
import type { RefreshResponse } from "./types";

let accessToken: string | null = null;
let refreshPromise: Promise<string> | null = null;

export function getAccessToken() {
  return accessToken;
}

export function setAccessToken(token: string | null) {
  accessToken = token;
}

export function clearAccessToken() {
  accessToken = null;
}

// یک Axios Instance جدا برای Refresh
const refreshClient = axios.create({
  baseURL: API_CONFIG.baseURL,
  timeout: API_CONFIG.timeout,

  // ارسال خودکار Cookie شامل Refresh Token
  withCredentials: true,

  headers: {
    Accept: "application/json",
  },
});

// درخواست دریافت Access Token جدید
async function performRefresh(): Promise<string> {
  const response = await refreshClient.post<RefreshResponse>(
    API_CONFIG.refreshEndpoint
  );

  const newAccessToken = response.data.accessToken;

  setAccessToken(newAccessToken);

  return newAccessToken;
}

// جلوگیری از ارسال چند درخواست Refresh هم‌زمان
export function refreshAccessToken(): Promise<string> {
  if (!refreshPromise) {
    refreshPromise = performRefresh().finally(() => {
      refreshPromise = null;
    });
  }

  return refreshPromise;
}