import axios, { AxiosInstance } from 'axios';

const AxiosClient: AxiosInstance = axios.create({
  baseURL: '/api/',
  withCredentials: true,
  timeout: 1000,
  headers: {
    "Content-Type": "application/json",
  }
});

export default AxiosClient;
