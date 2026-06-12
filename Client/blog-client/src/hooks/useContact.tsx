import { createContact, deleteContact } from "@/api/contactApi";
import { useMutation } from "@tanstack/react-query";

export function useCreateContact() {
  return useMutation({
    mutationFn: createContact,

    onSuccess: (data) => {
      console.log(data);
    },
    onError: (err) => {
      console.log(err);
    },
  });
}

export function useDeleteContact() {
  return useMutation({
    mutationFn: deleteContact,
    onSuccess: (data) => {
      console.log(data);
    },
    onError: (err) => {
      console.log(err);
    },
  });
}
