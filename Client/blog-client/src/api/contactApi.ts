import type { Contact } from "@/types/types";
import { api, axiosAuth } from "./axiosConfig";
import { API_VERSION } from "./blogApi";

interface CreateContactRequest {
  email: string;
  content: string;
}

export async function createContact(
  formData: CreateContactRequest,
): Promise<Contact> {
  const { data } = await api.post(`${API_VERSION}/contact/new`, formData);

  return data;
}

export async function deleteContact(contactID: number): Promise<string> {
  const { data } = await axiosAuth.delete(
    `${API_VERSION}/contact/${contactID}`,
  );

  return data;
}
