import { z } from "zod";

export const loginSchema = z.object({
  username: z
    .string()
    .min(4, "Must be at least 4 characters")
    .trim()
    .nonempty("Bắt buộc"),
  password: z
    .string()
    .min(8, "Must be at least 8 characters")
    .regex(/[a-zA-Z]/, { message: "Contain at least one letter." })
    .regex(/[0-9]/, { message: "Contain at least one number." })
    .regex(/[^a-zA-Z0-9]/, {
      message: "Contain at least one special character.",
    })
    .trim(),
});

export type LoginFormValues = z.infer<typeof loginSchema>;

export const signupSchema = z
  .object({
    username: z
      .string()
      .min(4, "Must be at least 4 characters")
      .trim()
      .nonempty("This is a required field"),
    password: z
      .string()
      .min(8, "Must be at least 8 characters")
      .regex(/[a-zA-Z]/, { message: "Contain at least one letter." })
      .regex(/[0-9]/, { message: "Contain at least one number." })
      .regex(/[^a-zA-Z0-9]/, {
        message: "Contain at least one special character.",
      })
      .trim(),
    confirmPassword: z
      .string()
      .min(8, "Confirm Password must be at least 8 characters"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match",
    path: ["confirmPassword"], // attach error to confirmPassword field
  });

export type SignupFormValues = z.infer<typeof signupSchema>;
