import { createContact, deleteContact } from "@/api/contactApi";
import { useMutation } from "@tanstack/react-query";

export function useCreateContact() {
  return useMutation({
    mutationFn: createContact,
  });
}

export function useDeleteContact() {
  return useMutation({
    mutationFn: deleteContact,
  });
}
