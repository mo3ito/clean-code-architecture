import { api } from "./client";

import { setupInterceptors } from "./interceptors";

import { http } from "./http";

let isInitialized = false;

function initializeApi() {
  if (isInitialized) {
    return;
  }

  setupInterceptors();

  isInitialized = true;
}

initializeApi();

export { api, http };

export {
  getAccessToken,
  setAccessToken,
  clearAccessToken,
} from "./auth";

export {
  normalizeApiError,
} from "./errors";

export type {
  ApiResponse,
  ApiErrorResponse,
} from "./types";

export type {
  ApiError,
  ApiErrorType,
} from "./errors";