import type {
  SignupFormValues,
  LoginFormValues,
} from "@/pages/auth/model/schema";
import { api, API_VERSION_V1, axiosAuth } from "./axiosConfig";

export async function loginRequest(formData: LoginFormValues) {
  const { data } = await api.post(`${API_VERSION_V1}/login`, formData);

  return data;
}

export async function signupRequest(formData: SignupFormValues) {
  const { data } = await api.post(`${API_VERSION_V1}/register`, formData);

  return data;
}

export async function logoutRequest() {
  await axiosAuth.post(`${API_VERSION_V1}/logout`);
}

export async function fetchMe() {
  const { data } = await axiosAuth.get(`${API_VERSION_V1}/me`);
  return data || null; // user
}
