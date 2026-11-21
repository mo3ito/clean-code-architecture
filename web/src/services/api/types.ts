export type ApiRequestOptions = Omit<
  RequestInit,
  "credentials"
>;

export interface ApiErrorResponse {
  message?: string;
  code?: string;
  errors?: Record<string, string[]>;
}