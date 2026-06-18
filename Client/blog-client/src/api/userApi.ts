import { API_VERSION_V1, axiosAuth } from "./axiosConfig";

export async function getNotifications() {
  const { data } = await axiosAuth.get(`${API_VERSION_V1}/user/notifications`);

  return data;
}
