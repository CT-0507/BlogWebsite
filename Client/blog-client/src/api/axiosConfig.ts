import axios from "axios";
import { tokenStore } from "./store/tokenStore";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
  timeout: 100,
});

let isRefreshing = false;
let failedQueue: {
  reject: (value: unknown) => void;
  resolve: (reason?: unknown) => void;
}[] = [];
const processQueue = (error: Error | null, response = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(response);
    }
  });

  failedQueue = [];
};
// async function refreshToken() {
//   const { data } = await axiosAuth.post("/refresh", null, {
//     _retry: true,
//   });
//   return data; // { accessToken }
// }
export const axiosAuth = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  withCredentials: true,
});

axiosAuth.interceptors.request.use(
  (config) => {
    const token = tokenStore.get();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

axiosAuth.interceptors.response.use(
  (res) => res,
  async (error) => {
    const originalRequest = error.config;

    if (
      (originalRequest.url.includes("/refresh") ||
        originalRequest.url.includes("/logout")) &&
      error.response.status === 401
    ) {
      //edge case where the refresh token is invalid or expired
      console.error("❌ Refresh token has expired or is invalid.");
      return Promise.reject(error); // fail directly, no retry
    }

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => axiosAuth(originalRequest))
          .catch((err) => Promise.reject(err));
      }

      isRefreshing = true;

      return new Promise((resolve, reject) => {
        axiosAuth
          .post("/refresh")
          .then((data) => {
            processQueue(null);
            tokenStore.set(data.data.accessToken);
            axiosAuth(originalRequest).then(resolve).catch(reject);
          })
          .catch((refreshError) => {
            processQueue(refreshError, null);
            tokenStore.clear(); // Clear auth state
            reject(refreshError); // fail the original promise chain
          })
          .finally(() => {
            isRefreshing = false;
          });
      });
    }
    //   try {
    //     const { accessToken } = await refreshToken();
    //     tokenStore.set(accessToken);

    //     queue.forEach((cb) => cb(accessToken));
    //     queue = [];

    //     originalRequest.headers.Authorization = `Bearer ${accessToken}`;
    //     return api(originalRequest);
    //   } catch (err) {
    //     console.log(err);
    //     tokenStore.clear();
    //     queue = [];
    //     window.location.href = "/account";
    //     return Promise.reject(error);
    //   } finally {
    //     isRefreshing = false;
    //   }
    // }

    return Promise.reject(error);
  }
);
