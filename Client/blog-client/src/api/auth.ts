import type {
  SignupFormValues,
  LoginFormValues,
} from "@/pages/auth/model/schema";
import { api, axiosAuth } from "./axiosConfig";
import { isLocalMode } from ".";
import { loginResponse, me } from "./mockApi";

export async function loginRequest(formData: LoginFormValues) {
  if (isLocalMode) return loginResponse;
  const { data } = await api.post("/login", formData);

  return data;
}

export async function signupRequest(formData: SignupFormValues) {
  if (isLocalMode) return {};
  const { data } = await api.post("/register", formData);

  return data;
}

export async function logoutRequest() {
  if (isLocalMode) {
    document.cookie =
      "refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    return { message: "Successfully logout" };
  }
  await axiosAuth.post("/logout");
}

export async function fetchMe() {
  if (isLocalMode) return me;
  const { data } = await axiosAuth.get("/me");
  return data || null; // user
}
