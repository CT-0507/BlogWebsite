import z from "zod";

export const publishBlogSchema = z.object({
  title: z
    .string()
    .max(200, "Blog title cannot exceed 15 characters")
    .trim()
    .nonempty("This is a required field"),
  content: z
    .string()
    .max(10000, "Blog title cannot exceed 15 characters")
    .trim()
    .nonempty("This is a required field"),
});

export type PublishBlogFormValues = z.infer<typeof publishBlogSchema>;
