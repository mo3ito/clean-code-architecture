import axios from "axios";

import type {
  ApiErrorResponse,
} from "./types";

export type ApiErrorType =
  | "HTTP_ERROR"
  | "NETWORK_ERROR"
  | "TIMEOUT_ERROR"
  | "CANCEL_ERROR"
  | "UNKNOWN_ERROR";

export interface ApiError {
  type: ApiErrorType;

  message: string;

  status?: number;

  code?: string;

  data?: ApiErrorResponse;
}

export function normalizeApiError(
  error: unknown
): ApiError {
  if (axios.isCancel(error)) {
    return {
      type: "CANCEL_ERROR",
      message: "Request cancelled",
    };
  }

  if (!axios.isAxiosError(error)) {
    return {
      type: "UNKNOWN_ERROR",
      message: "An unexpected error occurred",
    };
  }

  if (
    error.code === "ECONNABORTED" ||
    error.code === "ETIMEDOUT"
  ) {
    return {
      type: "TIMEOUT_ERROR",
      message: "Request timeout",
      code: error.code,
    };
  }

  if (!error.response) {
    return {
      type: "NETWORK_ERROR",
      message: "Unable to connect to the server",
      code: error.code,
    };
  }

  const responseData = error.response.data as
    | ApiErrorResponse
    | undefined;

  return {
    type: "HTTP_ERROR",

    message:
      responseData?.message ??
      error.message ??
      "Something went wrong",

    status: error.response.status,

    code: error.code,

    data: responseData,
  };
}