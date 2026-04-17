# **Delete Author Flow**

Saga context takes outbox context to update itself.
outbox payload is for step service to make side effect

## General flow

1. Delete Author
   Event name: DeleteAuthor
   Soft delete

- Step Emit:
  - DeleteAuthorUserCache
  - DeleteAuthorFailure

2. Delete Author-User cache on Blog module
   Event name: DeleteAuthorUserCache

- Step Emit:
  - Completed
  - DeleteAuthorUserCacheFailure

3. Hard delete (optinal)
   Event name: hard delete

## Context per step

Get userID and avatar url by authorID/slug before passing to delete saga

### Kickstart outbox

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "userID",
    "status": "deleted",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  },
  "payload": {
    "authorID": "author-1",
    "userID": "userID",
    "status": "active",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```

### 1. Delete Author

Saga initialize

```json
{
  "context": {
    "authorID": "author-1",
    "userID": "userID",
    "status": "deleted",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```

Saga step

```json
{
  "input": {
    "authorID": "author-1",
    "userID": "userID",
    "status": "deleted",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  },
  "output": null
}
```

Create outbox event

```json
{
  "context": null,
  "payload": {
    "authorID": "author-1",
    "userID": "userID",
    "status": "active",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  }
}
```

#### Failure

Do nothing, Emit failure event to mark saga context

#### Success

Emit Outbox event

```json
{
  "context": {
    "authorID": "author-1",
    "previousStatus": "active",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png"
  },
  "payload": {
    "authorID": "author-1"
  }
}
```

### 2. Delete Author User cache table

Saga step

```json
{
  "input": {
    "authorID": "author-1",
    "userID": "userID",
    "status": "deleted",
    "avatar": "/uploads/2026/04/08/author-90214.png"
  },
  "output": {
    "authorID": "author-1",
    "previousStatus": "active",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png"
  }
}
```

Update Saga context

```json
{
  "context": {
    "status": "deleted",
    "previousStatus": "active",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png",
    "authorID": "author-1",
    "userID": "userID"
  }
}
```

Create step 2

```json
{
  "input": {
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
    "authorID": "author-1"
  }
}
```

#### Failure

Emite event failure to go to previous step compensation event

#### Success:

```json
{
  "context": {
    "authorID": "author-1"
  },
  "payload": null
}
```

### 3. Complete

Update saga step

```json
{
  "input": {
    "authorID": "author-1"
  },
  "output": {
    "authorID": "author-1"
  }
}
```

Update saga context and mark as completed

```json
{
  "context": {
    "authorID": "author-1",
    "cache-id": "user-1-author-1",
    "userID": "user-1",
    "avatar": "/uploads/temp/2026/04/08/author-90214.png",
    "status": "deleted",
    "previousStatus": "active"
  }
}
```
