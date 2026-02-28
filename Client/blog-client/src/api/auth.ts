import type {
  SignupFormValues,
  LoginFormValues,
} from "@/pages/auth/model/schema";
import { api, axiosAuth } from "./axiosConfig";

export async function loginRequest(formData: LoginFormValues) {
  const { data } = await api.post("/login", formData);

  return data;
}

export async function signupRequest(formData: SignupFormValues) {
  const { data } = await api.post("/register", formData);

  return data;
}

export async function logoutRequest() {
  await axiosAuth.post("/logout");
}

export async function fetchMe() {
  const { data } = await axiosAuth.get("/me");
  return data || null; // user
}
