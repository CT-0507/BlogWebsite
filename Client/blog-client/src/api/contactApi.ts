import type { Contact } from "@/types/types";
import { api, axiosAuth, API_VERSION_V1 } from "./axiosConfig";

interface CreateContactRequest {
  email: string;
  content: string;
}

export async function createContact(
  formData: CreateContactRequest,
): Promise<Contact> {
  const { data } = await api.post(`${API_VERSION_V1}/contact/new`, formData);

  return data;
}

export async function deleteContact(contactID: number): Promise<string> {
  const { data } = await axiosAuth.delete(
    `${API_VERSION_V1}/contact/${contactID}`,
  );

  return data;
}
