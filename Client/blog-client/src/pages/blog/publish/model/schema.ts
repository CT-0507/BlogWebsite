import z from "zod";
import slugify from "slugify";

export const publishBlogSchema = z.object({
  title: z
    .string()
    .max(200, "Blog title cannot exceed 15 characters")
    .trim()
    .nonempty("This is a required field"),
  urlSlug: z
    .string()
    .min(1, "This is a required field")
    .max(400, "Url slug cannot exceed 400 characters")
    .transform((val) =>
      slugify(val, {
        lower: true,
        strict: true,
        trim: true,
      })
    ),
  content: z
    .string()
    .max(10000, "Blog title cannot exceed 15 characters")
    .trim()
    .nonempty("This is a required field"),
});

export type PublishBlogFormValues = z.infer<typeof publishBlogSchema>;
