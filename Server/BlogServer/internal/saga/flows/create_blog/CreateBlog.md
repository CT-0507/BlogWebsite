# Create Blog

## General flow

1. CreateBlog
2. IncreaseAuthorBlogCount
3. CreateNotifications and pushNotifications

## Context per step

### Step 1: CreateBlog

Kickstart outboxPayload

Get author id and user id from local cache table

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "user-1"
  },
  "payload": {
    "authorID": "author-1",
    "userID": "user-1",
    "status": "active",
    "title": "My first blog",
    "content": "This is my new blog",
    "url_slug": "my-first-blog"
  }
}
```

create saga

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "user-1"
  }
}
```

saga step 1

```json
{
  "input": {
    "authorID": "author-1",
    "userID": "user-1",
    "status": "active",
    "title": "My first blog",
    "content": "This is my new blog",
    "url_slug": "my-first-blog"
  },
  "output": null
}
```

this step's outbox, same as the kickstart outbox

```json
{
  "context": {},
  "payload": {
    "authorID": "author-1",
    "userID": "user-1",
    "status": "active",
    "title": "My first blog",
    "content": "This is my new blog",
    "url_slug": "my-first-blog"
  }
}
```

### Service handles event

#### Success outbox

outbox

```json
{
  "context": {
    "blogID": "blog-1"
  },
  "payload": {
    "userID": "user-1",
    "blogID": "blog-1",
    "authorID": "author-1",
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  }
}
```

### Step 2: IncreaseAuthorBlogCount

Receive payload from outbox event to create blog then emit event to saga to update saga state

saga

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "user-1",
    "blogID": "blog-1"
  }
}
```

saga step 1

```json
{
  "input": {
    "authorID": "author-1",
    "userID": "user-1",
    "status": "active",
    "title": "My first blog",
    "content": "This is my new blog",
    "url_slug": "my-first-blog"
  },
  "output": {
    "blogID": "blog-1"
  }
}
```

Create step 2

```json
{
  "input": {
    "userID": "user-1",
    "blogID": "blog-1",
    "authorID": "author-1",
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  },
  "output": null
}
```

Emit next step with outbox payload

```json
{
  "context": {},
  "payload": {
    "userID": "user-1",
    "blogID": "blog-1",
    "authorID": "author-1",
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  }
}
```

### Service handles event

### Success

outbox

```json
{
  "context": {
    "authorID": "author-1"
  },
  "payload": {
    "blogID": "blog-1",
    "authorID": "author-1",
    "userID": "user-1",
    "followerIds": ["user-2", "user-3"],
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  }
}
```

### Step 3: Create notification

saga

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "user-1",
    "blogID": "blog-1"
  }
}
```

update saga step 2
saga step 2

```json
{
  "input": {
    "blogID": "blog-1",
    "authorID": "author-1",
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  },
  "output": {
    "authorID": "author-1"
  }
}
```

Create step 3

```json
{
  "input": {
    "authorID": "author-1",
    "userID": "user-1",
    "followerIds": ["user-2", "user-3"],
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  },
  "output": null
}
```

outbox

```json
{
  "context": {},
  "payload": {
    "authorID": "author-1",
    "userID": "user-1",
    "followerIds": ["user-2", "user-3"],
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  }
}
```

### Service handles event

Create notifications and put notifications with userID to SSE broker

### Success

outbox

```json
{
  "context": {
    "notificationIds": ["not-1", "not-2"]
  },
  "payload": {
    "notificationIds": ["not-1", "not-2"]
  }
}
```

### Step 4: Complete

saga

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "user-1",
    "blogID": "blog-1",
    "notificationIds": ["not-1", "not-2"]
  }
}
```

saga step 3

```json
{
  "input": {
    "authorID": "author-1",
    "userID": "user-1",
    "followerIds": ["user-2", "user-3"],
    "truncatedTitle": "My first blog",
    "truncatedContent": "This is.."
  },
  "output": {
    "notificationIds": ["not-1", "not-2"]
  }
}
```
