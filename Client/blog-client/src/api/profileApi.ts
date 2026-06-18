import type {
  ChangeEmailFormValues,
  ChangeBasicInfoFormValues,
  ChangePasswordFormValues,
} from "@/pages/user/profile/model/schema";
import { axiosAuth } from "./axiosConfig";
import { API_VERSION_V1 } from "./axiosConfig";

export async function changeEmailRequest(formData: ChangeEmailFormValues) {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/user/change-email`,
    formData,
  );

  return data;
}

export async function changeBasicInfoRequest(
  formData: ChangeBasicInfoFormValues,
) {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/user/change-basic-info`,
    formData,
  );

  return data;
}

export async function changePasswordRequest(
  formData: ChangePasswordFormValues,
) {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/user/change-password`,
    formData,
  );

  return data;
}

export async function sendEmailCode(email: string) {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/user/change-email-code`,
    email,
  );

  return data;
}
