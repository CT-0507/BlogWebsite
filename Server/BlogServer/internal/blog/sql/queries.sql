-- name: GetBlog :one
SELECT 
    b.blog_id, 
    b.title,
    b.url_slug,
    b.author_id,
    b.content,
    b.like_count,
    b.dislike_count,
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE b.blog_id = $1 AND b.deleted_at IS NULL;

-- name: GetBlogWithUserReaction :one
SELECT 
    b.*,
    i.slug,
    i.display_name,
    r.type AS reaction_type
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
LEFT JOIN blogs.blog_reactions r
    ON r.blog_id = b.blog_id
    AND r.user_id = $1
WHERE b.url_slug = $2 AND b.deleted_at IS NULL;

-- name: GetBlogByUrlSlug :one
SELECT 
    b.*,
    i.slug,
    i.display_name,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE b.url_slug = $1 AND b.deleted_at IS NULL;

-- name: ListBlogsByAuthor :many
SELECT
    b.blog_id, 
    b.title,
    b.url_slug,
    b.author_id,
    b.content,
    b.like_count,
    b.dislike_count,
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE i.author_id = $1 AND b.deleted_at IS NULL AND b.status = $2;

-- name: ListBlogsByAuthorSlug :many
SELECT
    b.blog_id, 
    b.title,
    b.url_slug,
    b.author_id,
    b.content,
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    i.slug,
    i.display_name
FROM blogs.blogs b
JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
WHERE i.slug = $1 AND b.deleted_at IS NULL AND b.status = $2;

-- name: ListAllBlogs :many
SELECT blog_id, title, content, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM blogs.blogs;

-- name: GetListBlogsCount :one
WITH params AS (
    SELECT
        websearch_to_tsquery('english', COALESCE(sqlc.narg('title')::TEXT, '')) AS title_q,
        websearch_to_tsquery('english', COALESCE(sqlc.narg('content')::TEXT, '')) AS content_q,
        sqlc.narg('author_display_name')::TEXT AS author_q
)
SELECT
    COUNT(*) as total_result

FROM blogs.blogs b
JOIN blogs.idx_user_author_profile a ON a.author_id = b.author_id
CROSS JOIN params p

WHERE
    b.status = 'active'
    AND (
        (sqlc.narg('title')::TEXT IS NULL OR b.title_vector @@ p.title_q)
        AND (sqlc.narg('content')::TEXT IS NULL OR b.content_vector @@ p.content_q)
        AND (sqlc.narg('author_display_name')::TEXT IS NULL OR a.display_name ILIKE '%' || sqlc.narg('author_display_name')::TEXT || '%')
    );

-- name: ListBlogs :many
-- SELECT 
--     b.blog_id,
--     b.author_id,
--     b.title, 
--     b.url_slug,
--     b.content,
--     b.like_count,
--     b.dislike_count, 
--     b.status,
--     b.created_at, 
--     b.created_by, 
--     b.updated_at, 
--     b.updated_by,
--     i.slug,
--     i.display_name
-- FROM blogs.blogs b
-- JOIN blogs.idx_user_author_profile i ON i.author_id = b.author_id
-- WHERE b.deleted_at IS NULL;
WITH params AS (
    SELECT
        websearch_to_tsquery('english', COALESCE(sqlc.narg('title')::TEXT, '')) AS title_q,
        websearch_to_tsquery('english', COALESCE(sqlc.narg('content')::TEXT, '')) AS content_q,
        sqlc.narg('author_display_name')::TEXT AS author_q
)
SELECT
    b.blog_id,
    b.author_id,
    b.title, 
    b.url_slug,
    b.content,
    b.like_count,
    b.dislike_count, 
    b.status,
    b.created_at, 
    b.created_by, 
    b.updated_at, 
    b.updated_by,
    a.slug,
    a.display_name,

    (
        COALESCE(ts_rank(b.title_vector, p.title_q)::BIGINT, 0) * 2 +
        COALESCE(ts_rank(b.content_vector, p.content_q)::BIGINT, 0)
    ) AS rank

FROM blogs.blogs b
JOIN blogs.idx_user_author_profile a ON a.author_id = b.author_id
CROSS JOIN params p

WHERE
    b.status = 'active'
    AND (
        (sqlc.narg('title')::TEXT IS NULL OR b.title_vector @@ p.title_q)
        AND (sqlc.narg('content')::TEXT IS NULL OR b.content_vector @@ p.content_q)
        AND (sqlc.narg('author_display_name')::TEXT IS NULL OR a.display_name ILIKE '%' || sqlc.narg('author_display_name')::TEXT || '%')
    )

ORDER BY
    CASE WHEN sqlc.narg('sort_by')::TEXT = 'title' AND sqlc.narg('sort_dir')::TEXT = 'asc'  THEN b.title END ASC,
    CASE WHEN sqlc.narg('sort_by')::TEXT = 'title' AND sqlc.narg('sort_dir')::TEXT = 'desc' THEN b.title END DESC,

    CASE WHEN sqlc.narg('sort_by')::TEXT = 'created_at' AND sqlc.narg('sort_dir')::TEXT = 'asc'  THEN b.created_at END ASC,
    CASE WHEN sqlc.narg('sort_by')::TEXT = 'created_at' AND sqlc.narg('sort_dir')::TEXT = 'desc' THEN b.created_at END DESC,

    CASE
        WHEN sqlc.narg('sort_by')::TEXT = 'relevance' AND sqlc.narg('sort_dir')::TEXT = 'asc' THEN
            (COALESCE(ts_rank(b.title_vector, p.title_q)::BIGINT, 0) * 2 +
             COALESCE(ts_rank(b.content_vector, p.content_q)::BIGINT, 0)) END ASC,
    CASE
        WHEN sqlc.narg('sort_by')::TEXT = 'relevance' AND sqlc.narg('sort_dir')::TEXT = 'desc' THEN
            (COALESCE(ts_rank(b.title_vector, p.title_q)::BIGINT, 0) * 2 +
             COALESCE(ts_rank(b.content_vector, p.content_q)::BIGINT, 0)) END DESC,

    -- default sort (created_at DESC)
    b.created_at DESC

LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CreateBlog :one
INSERT INTO blogs.blogs(
    author_id,
    title,
    url_slug,
    content,
    created_by,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateBlog :one
UPDATE blogs.blogs
    SET title = $1,
    content = $2
WHERE blog_id = $3
RETURNING blog_id;

-- name: HardDeleteBlog :one
DELETE FROM blogs.blogs
WHERE blog_id = $1
RETURNING blog_id;

-- name: DeleteBlog :one
UPDATE blogs.blogs
    SET deleted_by = $1,
    deleted_at = NOW(),
    status = 'deleted'
WHERE blog_id = $2
RETURNING blog_id;

-- name: CreateUserAuthorProfileIDCacheRecord :exec
INSERT INTO blogs.idx_user_author_profile (
    user_id,
    author_id,
    slug,
    display_name
) VALUES (
    $1, $2, $3, $4
);

-- name: VerifyAuthorIDByUserID :one
SELECT author_id
FROM blogs.idx_user_author_profile
WHERE user_id = $1;

-- name: UpdateBlogStatusForDeletedAuthor :exec
UPDATE blogs.blogs
SET status = 'author_deleted',
deleted_at = NOW()
WHERE blogs.author_id = $1;

-- name: DeleteAuthorHardDeletedBlogs :exec
DELETE FROM blogs.blogs
WHERE author_id = $1;

-- name: DeleteAuthorCache :exec
DELETE FROM blogs.idx_user_author_profile
WHERE author_id = $1;

-- name: MarkAuthorCacheAsDeleted :exec
UPDATE blogs.idx_user_author_profile
SET status = 'deleted', deleted_at = NOW()
WHERE author_id = $1;

-- name: RestoreBlog :exec
UPDATE blogs.blogs
SET status = $1,
deleted_at = null,
deleted_by = null
WHERE blog_id = $2;


-- comments

-- name: CreateComment :one
WITH vals AS (
    SELECT gen_random_uuid() AS u1,
            gen_random_uuid() AS u2
)
INSERT INTO blogs.comments (
    id, blog_id, content, actor_type, actor_id, actor_display_name, parent_comment_id, root_comment_id, depth
)
SELECT u1, $1, $2, $3, $4, $5, $6,
    CASE 
        WHEN $8 = 0 THEN u1 
        ELSE $7
    END AS root_comment_id, 
    $8
FROM vals
RETURNING *;

-- name: GetBlogRootComment :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments c
    WHERE c.status = 'active'
    GROUP BY parent_comment_id
)
SELECT
    p.*,
    COALESCE(cc.cnt, 0) AS reply_count
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
WHERE
    p.blog_id = $1
    AND p.status = 'active'
    AND p.depth = 0;

-- name: GetBlogRootCommentWithUserReaction :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments c
    WHERE c.status = 'active'
    GROUP BY parent_comment_id
)
SELECT
    p.*,
    COALESCE(cc.cnt, 0) AS reply_count,
    r.type AS reaction_type
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
LEFT JOIN blogs.comment_reactions r
    ON r.comment_id = p.id
    AND r.user_id = $2
WHERE
    p.blog_id = $1
    AND p.depth = 0
    AND (
        p.status = 'active'
        OR (p.status = 'hidden' AND p.actor_id = $2)
    );

-- name: GetBlogRootCommentCount :one
SELECT COUNT(*) AS total
FROM blogs.comments c
WHERE c.blog_id = $1
    AND c.status = 'active';

-- name: GetCommentsByRootComment :many
SELECT *
FROM blogs.comments
WHERE root_comment_id = $1  AND status <> 'hidden';

-- name: GetCommentsByParentComment :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments c
    WHERE c.status = 'active'
    GROUP BY parent_comment_id
)
SELECT 
    p.*,
    COALESCE(cc.cnt, 0) AS child_comment_count
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
WHERE 
    p.parent_comment_id = $1 
    AND p.status = 'active';

-- name: GetCommentsByParentCommentUserWithReaction :many
WITH child_counts AS (
    SELECT parent_comment_id, COUNT(*) AS cnt
    FROM blogs.comments c
    WHERE c.status = 'active'
    GROUP BY parent_comment_id
)
SELECT 
    p.*,
    COALESCE(cc.cnt, 0) AS child_comment_count,
    r.type AS reaction_type
FROM blogs.comments p
LEFT JOIN child_counts cc
    ON cc.parent_comment_id = p.id
LEFT JOIN blogs.comment_reactions r
    ON r.comment_id = p.id
    AND r.user_id = $2
WHERE 
    p.parent_comment_id = $1 
    AND (
        p.status = 'active'
        OR (p.status = 'hidden' AND p.actor_id = $2)
    );

-- name: GetCommentByID :one
SELECT *
FROM blogs.comments
WHERE id = $1;

-- name: UpdateComment :one
UPDATE blogs.comments
SET
    -- wipe data on status deleted
    content = CASE
        WHEN COALESCE(sqlc.narg('status'), status) = 'deleted'
             AND status <> 'deleted'
        THEN 'deleted by user'
        ELSE COALESCE(sqlc.narg('content'), content)
    END,
    -- wipe data on status deleted
    actor_display_name = CASE
        WHEN COALESCE(sqlc.narg('status'), status) = 'deleted'
             AND status <> 'deleted'
        THEN 'deleted by user'
        ELSE actor_display_name
    END,
    actor_id= CASE
        WHEN COALESCE(sqlc.narg('status'), status) = 'deleted'
             AND status <> 'deleted'
        THEN NULL
        ELSE actor_id
    END,
    actor_avatar_url= CASE
        WHEN COALESCE(sqlc.narg('status'), status) = 'deleted'
             AND status <> 'deleted'
        THEN NULL
        ELSE actor_avatar_url
    END,
    updated_at = NOW(),
    deleted_at = CASE
        WHEN COALESCE(sqlc.narg('status'), status) = 'deleted'
             AND status <> 'deleted'
        THEN NOW()
        WHEN COALESCE(sqlc.narg('status'), status) <> 'deleted'
        THEN NULL
        ELSE deleted_at
    END
WHERE id = $1 AND (actor_id = $2 OR sqlc.arg('isAdmin')::BOOLEAN = TRUE)
RETURNING id;
-- WITH RECURSIVE tree AS (
--     SELECT id
--     FROM blogs.comments ch
--     WHERE ch.id = $1

--     UNION ALL

--     SELECT c.id
--     FROM blogs.comments c
--     JOIN tree t ON c.parent_comment_id = t.id
-- ),
-- updated AS (
--     UPDATE blogs.comments cm
--     SET
--         content = CASE 
--             WHEN cm.id = $1 THEN COALESCE(sqlc.narg('content'), cm.content)
--             ELSE cm.content
--         END,
--         status = COALESCE(sqlc.narg('status'), cm.status),
--         updated_at = NOW(),
--         deleted_at = CASE
--             WHEN COALESCE(sqlc.narg('status'), cm.status) = 'deleted'
--                  AND cm.status <> 'deleted'
--             THEN NOW()
--             WHEN COALESCE(sqlc.narg('status'), cm.status) <> 'deleted'
--             THEN NULL
--             ELSE cm.deleted_at
--         END
--     WHERE 
--         cm.id IN (SELECT id FROM tree)
--         AND (
--             cm.id <> $1
--             OR (cm.actor_id = $2 OR sqlc.arg('isAdmin')::BOOLEAN = TRUE)
--         )
--     RETURNING cm.id
-- )
-- SELECT id FROM updated;

-- name: UpdateBlogReactionCount :exec
UPDATE blogs.blogs
    SET like_count = like_count + $1,
    dislike_count = dislike_count + $2
WHERE blog_id = $3;

-- name: UpsertBlogReaction :one
WITH existing AS (
    SELECT type
    FROM blogs.blog_reactions e
    WHERE e.blog_id = $1 AND e.user_id = $2
    FOR UPDATE
),
upsert AS (
    INSERT INTO blogs.blog_reactions (blog_id, user_id, type)
    VALUES ($1, $2, $3)
    ON CONFLICT (blog_id, user_id)
    DO UPDATE SET type = EXCLUDED.type
    RETURNING type AS new_type
)
SELECT
    COALESCE((SELECT type FROM existing), 'none')::VARCHAR(20) AS old_type,
    (SELECT new_type FROM upsert)::VARCHAR(20) AS new_type;

-- name: UpsertCommentReaction :one
WITH existing AS (
    SELECT type
    FROM blogs.comment_reactions e
    WHERE e.comment_id = $1 AND e.user_id = $2
    FOR UPDATE
),
upsert AS (
    INSERT INTO blogs.comment_reactions (comment_id, user_id, type)
    VALUES ($1, $2, $3)
    ON CONFLICT (comment_id, user_id)
    DO UPDATE SET type = EXCLUDED.type
    RETURNING type AS new_type
)
SELECT
    COALESCE((SELECT type FROM existing), 'none')::VARCHAR(20) AS old_type,
    COALESCE((SELECT new_type FROM upsert), 'none')::VARCHAR(20) AS new_type;

-- name: UpdateCommentReactionCount :exec
UPDATE blogs.comments
    SET like_count = like_count + $1,
    dislike_count = dislike_count + $2
WHERE id = $3;

-- name: SyncBlogLikeAndDislike :exec
UPDATE blogs.blogs b
SET
    like_count = COALESCE(x.like_count, 0),
    dislike_count = COALESCE(x.dislike_count, 0)
FROM (
    SELECT
        blog_id,
        COUNT(*) FILTER (WHERE type = 'like' AND status = 'active') AS like_count,
        COUNT(*) FILTER (WHERE type = 'dislike' AND status = 'active') AS dislike_count
    FROM blogs.blog_reactions
    GROUP BY blog_id
) x
WHERE x.blog_id = b.blog_id;

-- name: GetAuthorCacheByUserID :one
SELECT *
FROM blogs.idx_user_author_profile
WHERE user_id = $1;

-- name: ListRankingTable :many
-- params:
-- :sort_by      -> 'daily' | 'weekly' | 'likes' | 'score' | 'rank'
-- :sort_dir     -> 'asc' | 'desc'
-- :limit        -> int
-- :offset       -> int
-- :get_all      -> boolean
WITH filtered AS (
    SELECT *
    FROM blogs.blog_ranking br
    WHERE
        CASE
            WHEN sqlc.arg('type')::TEXT = 'allTime' THEN br.rank_all_time
            WHEN sqlc.arg('type')::TEXT = 'trending' THEN br.rank_trending
        END IS NOT NULL
),
    top20 AS (
SELECT
    ft.*
FROM filtered ft
ORDER BY
    CASE
        WHEN sqlc.arg('type')::TEXT = 'allTime' THEN ft.rank_all_time
        WHEN sqlc.arg('type') = 'trending' THEN ft.rank_trending
    END ASC,
    -- daily access
    CASE WHEN sqlc.arg('sort_by')::TEXT = 'daily'  AND sqlc.arg('sort_dir')::TEXT = 'asc'  THEN ft.daily_access_count END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'daily'  AND sqlc.arg('sort_dir') = 'desc' THEN ft.daily_access_count END DESC,

    -- weekly access
    CASE WHEN sqlc.arg('sort_by') = 'weekly' AND sqlc.arg('sort_dir') = 'asc'  THEN ft.weekly_access_count END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'weekly' AND sqlc.arg('sort_dir') = 'desc' THEN ft.weekly_access_count END DESC,

    -- like count
    CASE WHEN sqlc.arg('sort_by') = 'likes'  AND sqlc.arg('sort_dir') = 'asc'  THEN ft.like_count END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'likes'  AND sqlc.arg('sort_dir') = 'desc' THEN ft.like_count END DESC,

    -- score
    CASE WHEN sqlc.arg('sort_by') = 'score' AND sqlc.arg('type') = 'allTime'  AND sqlc.arg('sort_dir') = 'asc'  THEN ft.score_all_time END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'score' AND sqlc.arg('type') = 'allTime' AND sqlc.arg('sort_dir') = 'desc' THEN ft.score_all_time END DESC,
    CASE WHEN sqlc.arg('sort_by') = 'score' AND sqlc.arg('type') = 'trending' AND sqlc.arg('sort_dir') = 'asc'  THEN ft.score_trending END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'score' AND sqlc.arg('type') = 'trending' AND sqlc.arg('sort_dir') = 'desc' THEN ft.score_trending END DESC,

    -- rank
    CASE WHEN sqlc.arg('sort_by') = 'rank' AND sqlc.arg('type') = 'allTime' AND sqlc.arg('sort_dir') = 'asc'  THEN ft.rank_all_time END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'rank' AND sqlc.arg('type') = 'allTime' AND sqlc.arg('sort_dir') = 'desc' THEN ft.rank_all_time END DESC,
    CASE WHEN sqlc.arg('sort_by') = 'rank' AND sqlc.arg('type') = 'trending' AND sqlc.arg('sort_dir') = 'asc'  THEN ft.rank_trending END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'rank' AND sqlc.arg('type') = 'trending' AND sqlc.arg('sort_dir') = 'desc' THEN ft.rank_trending END DESC
    LIMIT 20
)
SELECT t.*, b.title, b.author_id, b.url_slug, a.avatar, a.display_name, a.slug, 
    COUNT(t.rank_all_time) OVER() as total_all_time, 
    COUNT(t.rank_trending) OVER() as total_trending
FROM top20 t
LEFT JOIN blogs.blogs b ON b.blog_id = t.blog_id AND b.status = 'active'
LEFT JOIN blogs.idx_user_author_profile a ON a.author_id = b.author_id AND a.deleted_at IS NULL 
-- paging control
LIMIT
CASE
    WHEN sqlc.arg('get_all')::BOOLEAN THEN 20
    ELSE LEAST(sqlc.arg('limit')::INT, 20)
END
OFFSET
CASE
    WHEN sqlc.arg('get_all')::BOOLEAN THEN 0
    ELSE sqlc.arg('offset')::INT
END;

-- name: TruncateBlogRankingTable :exec
TRUNCATE blogs.blog_ranking;

-- name: UpdateBlogRankingResult :exec
WITH comment_stats AS (
    SELECT
        c.blog_id,
        COUNT(*) AS comment_count
    FROM blogs.comments c
    GROUP BY c.blog_id
),
base AS (
    SELECT
        b.blog_id,
        b.created_at,
        b.daily_access_count,
        b.weekly_access_count,
        b.like_count,
        b.dislike_count,
        COALESCE(cs.comment_count, 0) AS comment_count,

        -- base score (shared)
        (
            b.daily_access_count * 2 +
            b.weekly_access_count * 1 +
            (b.like_count + 0.25 * b.dislike_count) * 2 +
            (b.like_count - b.dislike_count) * 1 +
            COALESCE(cs.comment_count, 0) * 2
        ) AS base_score,

        EXTRACT(EPOCH FROM (NOW() - b.created_at)) / 3600 AS hours_since_created

    FROM blogs.blogs b
    LEFT JOIN comment_stats cs
        ON cs.blog_id = b.blog_id
    WHERE b.status = 'active'
),
scored AS (
    SELECT
        *,
        base_score AS score_all_time,
        base_score / (1 + hours_since_created) AS score_trending
    FROM base
),
ranked AS (
    SELECT
        *,
        RANK() OVER (ORDER BY score_all_time DESC) AS rank_all_time,
        RANK() OVER (ORDER BY score_trending DESC) AS rank_trending
    FROM scored
)
INSERT INTO blogs.blog_ranking (
    blog_id,
    rank_all_time,
    rank_trending,
    score_all_time,
    score_trending,
    like_count,
    dislike_count,
    comment_count,
    weekly_access_count,
    daily_access_count,
    created_at,
    computed_at
)
SELECT
    blog_id,

    -- only keep rank if in top 20
    CASE WHEN rank_all_time <= 20 THEN rank_all_time ELSE NULL END,
    CASE WHEN rank_trending <= 20 THEN rank_trending ELSE NULL END,

    score_all_time,
    score_trending,
    like_count,
    dislike_count,
    comment_count,
    weekly_access_count,
    daily_access_count,
    created_at,
    NOW()

FROM ranked
WHERE rank_all_time <= 20
   OR rank_trending <= 20;

-- name: InsertTags :exec
WITH input_tags AS (
  SELECT DISTINCT LOWER(TRIM(UNNEST(sqlc.arg('name')::text[]))) AS name
),
inserted AS (
  INSERT INTO blogs.tags (name)
  SELECT name FROM input_tags
  ON CONFLICT (name) DO NOTHING
),
all_tags AS (
  SELECT id, name FROM blogs.tags
  WHERE name IN (SELECT name FROM input_tags)
)
INSERT INTO blogs.blog_tags (blog_id, tag_id)
SELECT $1, id FROM all_tags
ON CONFLICT DO NOTHING;