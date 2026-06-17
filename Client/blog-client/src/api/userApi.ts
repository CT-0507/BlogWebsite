import { axiosAuth } from "./axiosConfig";

export async function getNotifications() {
  const { data } = await axiosAuth.get("/user/notifications");

  return data;
}
