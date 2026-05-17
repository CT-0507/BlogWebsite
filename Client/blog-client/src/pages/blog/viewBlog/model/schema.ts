import { z } from "zod";

const MAX_COMMENT_LENGTH = 1000;
const MIN_COMMENT_LENGTH = 1;

export const postCommentSchema = z.object({
  content: z
    .string()
    .trim()
    .min(
      MIN_COMMENT_LENGTH,
      `Must be at least ${MIN_COMMENT_LENGTH} characters`
    )
    .max(MAX_COMMENT_LENGTH, `Cannot exceed ${MAX_COMMENT_LENGTH} characters`)
    .nonempty("This is a required field"),
  actorType: z.enum(["user", "creator"], {
    error: "actorType must be user or creator",
  }),
  parentCommentId: z.string().nullable().optional(),
  rootCommentId: z.string().nullable().optional(),
  blogID: z.number(),
  depth: z.number().min(0).max(2),
});

export type PostCommentFormValues = z.infer<typeof postCommentSchema>;

export const updateCommentContentSchema = z.object({
  content: z
    .string()
    .trim()
    .min(
      MIN_COMMENT_LENGTH,
      `Must be at least ${MIN_COMMENT_LENGTH} characters`
    )
    .max(MAX_COMMENT_LENGTH, `Cannot exceed ${MAX_COMMENT_LENGTH} characters`)
    .nonempty("This is a required field"),
  commentId: z.string().nonempty(),
});

export type UpdateCommentContentFormValues = z.infer<
  typeof updateCommentContentSchema
>;

export const reportSchema = z.object({
  blogID: z.number(),
  reason: z.string().trim().min(2).max(500),
});

export type BlogReportFormValues = z.infer<typeof reportSchema>;
