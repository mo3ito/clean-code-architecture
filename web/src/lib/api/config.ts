export const API_CONFIG = {
  baseURL: process.env.NEXT_PUBLIC_API_URL,

  timeout: 15_000,

  refreshEndpoint: "/auth/refresh",

  headers: {
    Accept: "application/json",
  },
} as const;