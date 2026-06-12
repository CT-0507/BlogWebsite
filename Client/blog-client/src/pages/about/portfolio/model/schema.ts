import { z } from "zod";

export const contactSchema = z.object({
  email: z
    .email()
    .min(4, "Must be at least 4 characters")
    .max(50, "Password cannot exceed 50 characters")
    .trim()
    .nonempty("This is a required field"),
  content: z
    .string()
    .min(8, "Must be at least 8 characters")
    .max(1000, "Password cannot exceed 1000 characters")
    .nonempty("This is a required field"),
});

export type ContactFormValues = z.infer<typeof contactSchema>;
