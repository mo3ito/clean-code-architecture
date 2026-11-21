import axios from "axios";
import {
  setupRequestInterceptor,
  setupResponseInterceptor,
  setupResponseErrorInterceptor,
} from "./interceptors";

export const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

api.interceptors.request.use(setupRequestInterceptor);
api.interceptors.response.use(
  setupResponseInterceptor,
  setupResponseErrorInterceptor,
);
