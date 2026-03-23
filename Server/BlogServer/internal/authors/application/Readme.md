Author Module — Use Cases
Module Responsibility

- The Author Module is responsible for:

- Author identity & lifecycle

- Author profile management

- Author relationships (followers)

- Author-owned presentation (e.g. featured posts — IDs only)

**1. Author Identity & Lifecycle**

✅ Create Author Profile

```
createAuthor(userId, profileData)
```

Description:
Create an author profile linked to a user.

Rules:

One user ↔ one author (or configurable)

Generate unique slug

Initialize default profile

✅ Get Author by ID

```
getAuthorById(authorId)
```

✅ Get Author by Slug

```
getAuthorBySlug(slug)
```

✅ List Authors (with filters)

```
listAuthors(filterParams)
```

Possible filters:

- keyword (name)

- category/tag

- popularity (followers count)

- created date

✅ Update Author Profile

```
updateAuthorProfile(authorId, profileData)
```

✅ Delete Author (Soft Delete Recommended)

```
deleteAuthor(authorId)
```

Behavior:

Mark as deleted / inactive

Emit domain event: AuthorDeleted
✅ Suspend Author (Admin)

```
suspendAuthor(authorId)
```

**2. Author Slug & Identity Management**

✅ Update Author Slug

```
updateAuthorSlug(authorId, newSlug)
```

Rules:

- Must be unique

- Should handle redirect (optional)

**3. Author Social (Follow System)**

✅ Follow Author

```
followAuthor(userId, authorId)
```

✅ Unfollow Author

```
unfollowAuthor(userId, authorId)
```

✅ Get Author Followers

```
getAuthorFollowers(authorId, pagination)
```

✅ Get Followed Authors (by user)

```
getFollowedAuthors(userId, pagination)
```

**4. Author Social (Follow System)**

✅ Set Featured Posts (Curated by Author)

```
setAuthorFeaturedPosts(authorId, postIds[])
```

Description:
Author selects posts to highlight on their profile.

Rules:

- Max limit (e.g. 5–10 posts)

- Posts must belong to the author (validated via Post module)

- Order is preserved

- Store only post IDs

✅ Get Featured Posts (IDs only)

```
getAuthorFeaturedPostIds(authorId)
```

Post module will resolve full data

**5. Author Metadata & Insights**

✅ Get Author Stats

```
getAuthorStats(authorId)
```

Returns:

- follower count

- post count (optional via Post module)

- other aggregates

**6. Author Discovery**

✅ Search Authors

```
searchAuthors(keyword, pagination)
```

✅ Get Trending Authors (Optional)

```
getTrendingAuthors()
```

👉 Implementation may depend on:

- followers growth

- activity

**7. Internal / Domain Events**

📢 AuthorCreated

Triggered when:

```
createAuthor()
```

📢 AuthorUpdated

Triggered when:

```
updateAuthorProfile()
```

📢 AuthorDeleted

Triggered when:

```
deleteAuthor()
```

📢 AuthorFollowed

Triggered when:

```
followAuthor()
```
