import "axios";

export interface ApiResponse<T> {
  success: boolean;
  message?: string;
  data: T;
}

export interface ApiErrorResponse {
  success?: boolean;
  message: string;
  statusCode?: number;
  errors?: Record<string, string[]>;
}

export interface RefreshResponse {
  accessToken: string;
}

declare module "axios" {
  export interface AxiosRequestConfig {
    requiresAuth?: boolean;

    _retry?: boolean;

    skipAuthRefresh?: boolean;
  }

  export interface InternalAxiosRequestConfig {
    requiresAuth?: boolean;

    _retry?: boolean;

    skipAuthRefresh?: boolean;
  }
}