import { isLocalMode } from ".";
import { axiosAuth } from "./axiosConfig";
import { notifications } from "./mockApi";

export async function getNotifications() {
  if (isLocalMode) return notifications;
  const { data } = await axiosAuth.get("/user/notifications");

  return data;
}
