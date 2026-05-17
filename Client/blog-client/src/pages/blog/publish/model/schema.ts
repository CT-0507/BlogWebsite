import z from "zod";
import slugify from "slugify";

export const editorContentSchema = z.object({
  // json: z.object({
  //   blocks: z.array(z.any()).min(1, "Editor must contain at least one block"),
  // }),

  plainText: z
    .string()
    .min(20, "Content must be at least 20 characters")
    .max(10000, "Blog title cannot exceed 10000 characters"),
});

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
  content: z.any().optional(),
  tags: z.array(z.string().min(1).max(20)).max(5).optional().nullable(),
  thumbnail: z.instanceof(File).nullable().optional(),
});

export type PublishBlogFormValues = z.infer<typeof publishBlogSchema>;
