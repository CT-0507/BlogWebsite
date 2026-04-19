# **Create Author Flow**

Saga context takes outbox context to update itself.
outbox payload is for step service to make side effect

## General flow

1. Create Author
   Event name: CreateAuthor

- Step Emit:
  - CreateAuthorUserCache
  - CreateAuthorFailure

2. Create Author-User cache on Blog module
   Event name: CreateAuthorUserCache

- Step Emit:
  - Completed
  - CreateAuthorUserCacheFailure

## Context per step

### Kickstart outbox event

Create author Handler after saving the image and put the url into avatar then creates outbox event.

#### Kickstart outbox

```json
{
  "context": {
    "userID": "user-1",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png"
  },
  "payload": {
    "userID": "user-1",
    "displayName": "First Author",
    "bio": "Say some thing here",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png",
    "slug": "first-author",
    "socialLink": null,
    "email": null
  }
}
```

### 1. Create Author

Saga initialize

```json
{
  "context": {
    "userID": "user-1",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png"
  }
}
```

Saga Step 1

```json
{
  "input": {
    "userID": "user-1",
    "displayName": "First Author",
    "bio": "Say some thing here",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png",
    "slug": "first-author",
    "socialLink": null,
    "email": null
  },
  "output": null
}
```

Create outbox event

```json
{
  "context": null,
  "payload": {
    "userID": "user-1",
    "displayName": "First Author",
    "bio": "Say some thing here",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png",
    "slug": "first-author",
    "socialLink": null,
    "email": null
  }
}
```

#### Execute step

#### On Failure

Do nothing

#### On Success

outbox event

```json
{
  "context": {
    "authorID": "author-1",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  },
  "payload": {
    "userID": "user-1",
    "authorID": "author-1"
  }
}
```

### 2. Create Author-User cache on blog module

Update step 1 output

```json
{
  "input": {
    "userID": "user-1",
    "displayName": "First Author",
    "bio": "Say some thing here",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png",
    "slug": "first-author",
    "socialLink": null,
    "email": null
  },
  "output": {
    "authorID": "author-1",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```

Update Saga context

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "user-1",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```

Create step 2

```json
{
  "input": {
    "userID": "user-1",
    "authorID": "author-1"
  },
  "output": null
}
```

Create outbox event

```json
{
  "context": null,
  "payload": {
    "userID": "user-1",
    "authorID": "author-1"
  }
}
```

#### Failure

Only change in outbox event name to failure

```json
{
  "context": null,
  "payload": null
}
```

Takes output from previous step to do compensation event.
Delete author and avatar
Step 1 output

```json
{
  "output": {
    "authorID": "author-1",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```

#### Success:

Create outbox event

```json
{
  "context": {
    "userID": "user-1",
    "authorID": "author-1"
  },
  "payload": null
}
```

### 3. Completed

Update Step 2 output

```json
{
  "input": {
    "userID": "user-1",
    "authorID": "author-1"
  },
  "output": {
    "userID": "user-1",
    "authorID": "author-1"
  }
}
```

Update Saga context and mark as complete

```json
{
  "context": {
    "authorID": "author-1",
    "cache-id": "user-1-author-1",
    "userID": "user-1",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```
