About table design

- Why do all tables do not have any relation

In your design, each module behaves like a mini service inside the monolith.

That means:

❗ The database should NOT tightly couple modules together

Foreign keys create hard coupling at the database level, which breaks that principle.

##

If you had this:

FOREIGN KEY (user_id) REFERENCES users(id)

You are saying:

Author module depends on User table

You cannot:

- delete user independently

- deploy modules independently

- split into microservices later

👉 This violates modular boundaries

## With eventually consistent

Modules communicate via events, not joins

Data is duplicated, not shared via FK

So architecture is closer to:

```
Author Module DB  ←→  Blogt Module DB  ←→  User Module DB
        (no FK)           (no FK)           (no FK)
             ↘ events / async sync ↙
```

## Benifits

✅ Loose Coupling

- Author module doesn’t “know” User schema
- Post module changes won’t break Author DB

✅ Independent Evolution

You can:

- split into microservices later
- scale modules independently

✅ Better Performance (your goal)

- No cross-table joins
- Faster reads via denormalization

✅ Event-Driven Design Becomes Natural

Instead of FK constraints:

- PostCreated → Author.post_count++
- UserUpdated → update cached follower info

### Downside

❌ No Referential Integrity at DB Level

The DB will NOT stop this:

```sql
author_featured_posts.post_id = "non_existing_post"
```

👉 You must enforce it in application layer

❌ Risk of Inconsistent Data

Example:

- Blog deleted
- But still exists in featured_blogs

👉 You fix this via:

- events
- background jobs

❌ More Responsibility in Code

You now must handle:

- validation
- cleanup
- sync logic

### Solution

✅ Application-Level Validation

Example:

setAuthorFeaturedPosts(authorId, postIds)

Flow:

- Call Post module → validate posts exist
- Ensure posts belong to author
- Save IDs

✅ Domain Events

- BlogDeleted → remove from featured_blogs
- UserDeleted → remove followers
  ✅ Background Reconciliation (Optional but powerful)

periodic job:

- remove orphaned blog_ids
- fix counters
