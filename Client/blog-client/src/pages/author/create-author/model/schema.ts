import slugify from "slugify";
import { z } from "zod";

export const createAuthorSchema = z.object({
  displayName: z
    .string()
    .min(4, "Must be at least 4 characters")
    .max(100, "Cannot exceed 150 characters")
    .trim()
    .nonempty("This is a required field"),
  bio: z.string().max(1000, "Cannot exceed 1000 characters").trim(),
  avatar: z.file().nullable(),
  slug: z
    .string()
    .min(1, "This is a required field")
    .max(400, "Url slug cannot exceed 400 characters")
    .trim()
    .transform((val) =>
      slugify(val, {
        lower: true,
        strict: true,
        trim: true,
      })
    ),
  socialLink: z.string().max(300, "Cannot exceed 300 characters").trim(),
  email: z
    .union([z.literal(""), z.email()])
    .transform((e) => (e === "" ? undefined : e))
    .optional(),
});

export type CreateAuthorFormValues = z.infer<typeof createAuthorSchema>;
