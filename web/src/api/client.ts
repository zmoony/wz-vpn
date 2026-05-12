import axios from "axios";

const client = axios.create({
  baseURL: "/api",
  withCredentials: true,
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error?.response?.data?.message ?? "请求失败";
    return Promise.reject(new Error(message));
  },
);

export default client;
