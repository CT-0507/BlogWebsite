import z from "zod";

const MAX_CONTENT_LENGTH = 100;
const MAX_TITLE_LENGTH = 100;
const MAX_AUTHOR_LENGTH = 100;

export const searchBlogsSchema = z.object({
  title: z
    .string()
    .trim()
    .max(MAX_TITLE_LENGTH, `Cannot exceed ${MAX_TITLE_LENGTH} characters`)
    .optional()
    .nullable(),
  content: z
    .string()
    .trim()
    .max(MAX_CONTENT_LENGTH, `Cannot exceed ${MAX_CONTENT_LENGTH} characters`)
    .optional()
    .nullable(),
  author: z
    .string()
    .trim()
    .max(MAX_AUTHOR_LENGTH, `Cannot exceed ${MAX_AUTHOR_LENGTH} characters`)
    .optional()
    .nullable(),
  sortBy: z.enum(["title", "created_at", "relevance"], {
    error: "sortBy invalid",
  }),
  sortDir: z.enum(["asc", "desc"], {
    error: "sortBy invalid",
  }),
  limit: z.number().min(10, "limit is under 5").max(100, "limit is over 100"),
});

export type SearchBlogsFormValues = z.infer<typeof searchBlogsSchema>;
