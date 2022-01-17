import AxiosClient from "../lib/axios";

export const checkLogin =  (email: string, password: string): Promise<any> => {
  return AxiosClient
    .post("/login", { email, password })
    .then((res) => res.data)
    .catch((error) => {
      throw error;
    })
};
