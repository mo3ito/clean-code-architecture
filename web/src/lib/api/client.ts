import axios from "axios";

import { API_CONFIG } from "./config";

export const api = axios.create({
  baseURL: API_CONFIG.baseURL,

  timeout: API_CONFIG.timeout,

  headers: API_CONFIG.headers,

  withCredentials: true,
});