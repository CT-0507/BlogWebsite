import type {
  ChangeEmailFormValues,
  ChangeBasicInfoFormValues,
  ChangePasswordFormValues,
} from "@/pages/user/profile/model/schema";
import { axiosAuth } from "./axiosConfig";
import { isLocalMode } from ".";
import { updateUserData } from "./mockApi";

export async function changeEmailRequest(formData: ChangeEmailFormValues) {
  if (isLocalMode) return updateUserData;
  const { data } = await axiosAuth.post("/user/change-email", formData);

  return data;
}

export async function changeBasicInfoRequest(
  formData: ChangeBasicInfoFormValues
) {
  if (isLocalMode) return updateUserData;
  const { data } = await axiosAuth.post("/user/change-basic-info", formData);

  return data;
}

export async function changePasswordRequest(
  formData: ChangePasswordFormValues
) {
  if (isLocalMode) return updateUserData;
  const { data } = await axiosAuth.post("/user/change-password", formData);

  return data;
}

export async function sendEmailCode(email: string) {
  if (isLocalMode) return { code: "123456" };
  const { data } = await axiosAuth.post("/user/change-email-code", email);

  return data;
}
