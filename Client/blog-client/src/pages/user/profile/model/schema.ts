import { z } from "zod";

export const changeBasicInfoSchema = z.object({
  firstName: z
    .string()
    .min(1, "Must be at least 1 characters")
    .max(15, "First name cannot exceed 15 characters")
    .trim()
    .nonempty("This is a required field"),
  lastName: z
    .string()
    .min(1, "Must be at least 1 characters")
    .max(15, "Last name cannot exceed 15 characters")
    .trim(),
});

export type ChangeBasicInfoFormValues = z.infer<typeof changeBasicInfoSchema>;

export const changePasswordSchema = z
  .object({
    currentPassword: z
      .string()
      .min(8, "Must be at least 8 characters")
      .max(30, "Current password cannot exceed 30 characters")
      .regex(/[a-zA-Z]/, { message: "Contain at least one letter." })
      .regex(/[0-9]/, { message: "Contain at least one number." })
      .regex(/[^a-zA-Z0-9]/, {
        message: "Contain at least one special character.",
      })
      .trim(),
    newPassword: z
      .string()
      .min(8, "Must be at least 8 characters")
      .max(30, "Password cannot exceed 30 characters")
      .regex(/[a-zA-Z]/, { message: "Contain at least one letter." })
      .regex(/[0-9]/, { message: "Contain at least one number." })
      .regex(/[^a-zA-Z0-9]/, {
        message: "Contain at least one special character.",
      })
      .trim(),
    confirmNewPassword: z
      .string()
      .min(8, "Confirm Password must be at least 8 characters"),
  })
  .refine((data) => data.newPassword === data.confirmNewPassword, {
    message: "Passwords do not match",
    path: ["confirmNewPassword"], // attach error to confirmPassword field
  })
  .refine((data) => data.newPassword !== data.currentPassword, {
    message: "New password cannot be the same as current password",
    path: ["newPassword"], // attach error to confirmPassword field
  });

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;

export const changeEmailSchema = z.object({
  email: z
    .email("Email is not valid")
    .max(40, "First name cannot exceed 15 characters")
    .trim()
    .nonempty("This is a required field"),
  confirmCode: z
    .string()
    .regex(/^\d+$/, {
      message: "Must contain only numbers",
    })
    .length(6, "Length is not valid")
    .trim()
    .nonempty("This is a required field"),
});

export type ChangeEmailFormValues = z.infer<typeof changeEmailSchema>;
