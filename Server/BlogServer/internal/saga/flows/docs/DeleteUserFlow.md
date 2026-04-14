# **Delete User Flow**

Saga context takes outbox context to update itself.
outbox payload is for step service to make side effect

## General flow

1. Soft delete User
Event name: DeleteUser
- Step Emit
    - DeleteAuthorProfile
    - DeleteUserFailure
2. DeleteAUthorProfile
- Step Emit
    - HardDeleteUser
    - DeleteAuthorProfileFailure
3. DeleteBlogAuthorUserCache
- Step Emit
    - DeleteUserCompleted
    - HardDeleteUserFailure

## Context per step

### Kickstart outbox event

Check user status
if not exists or is deleted then exit.

```json
{
    "context": {
        "userID": "user-1",
        "status": "deleted",
    },
    "payload": {
        "status": "deleted",
        "userID": "user-1"
    }
}
```

### 1. Delete User

Saga initialize

```json
{
    "context": {
        "userID": "user-1",
        "status": "deleted",
    }
}
```

Saga step 1

```json
{
    "input": {
        "userID": "user-1",
        "status": "deleted",
    },
    "output": null
}
```

Create outbox event

```json
{
    "context": null,
    "payload": {
        "userID": "user",
        "status": "deleted",
    }
}
```

### Success

Event type : DeleteAuthorProfile

```json
{
    "context": {
        "userPreviousStatus": "active",
        "userID": "user-1",
    },
    "payload": {
        "userID": "user-1"
    }
}
```

### Failure
Emit
Event type : DeleteUserFailed with no context nor payload

### 2. DeleteAuthorProfile

Update Step 1
```json
{
    "input": {
        "userID": "user-1",
        "status": "deleted",
    },
    "output": {
        "userPreviousStatus": "active",
        "userID": "user-1",
    }
}
```

Update saga context
```json
{
    "context": {
        "userID": "user-1",
        "userPreviousStatus": "active",
        "status": "deleted",
    }
}
```

Create step 2
```json
{
    "input": {
        "userID": "user-1",
    },
    "output": null
}
```

Create outbox event

```json
{
    "context": null,
    "payload": {
        "userID": "user-1"
    }
}
```

#### Success
```json
{
    "context": {
        "author-id": "author-1",
        "authorPreviousState": "active",
    },
    "payload": {
        "userID": "user-1",
        "author-id": "author-1"
    }
}
```

#### Failure 
Emit DeleteAuthorProfileFailed


### 3. Delete Blog Author User cache
Update Step 2
```json
{
    "input": {
        "userID": "user-1",
    },
    "output": {
        "author-id": "author-1",
        "authorPreviousState": "active",
    }
}
```

Update saga context
```json
{
    "context": {
        "userID": "user-1",
        "previousStatus": "active",
        "authorPreviousState": "active",
        "status": "deleted",
        "author-id": "author-1",
    }
}
```
Create Saga step 3
```json
{
    "input": {
        "userID": "user-1",
        "author-id": "author-1",
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
        "author-id": "author-1",
    }
}
```

#### Success

Create outbox event step 3
```json
{
    "payload": null,
    "output": {
        "userID": "user-1",
        "author-id": "author-1",
        "cachePreviousState": "active",
    }
}
```
#### Failure
Emit Delete Blog Author User cache failed, mark as failed and start compensation

### 4. Complete
Emit Delete Blog Author User cache success and mark saga context as success

Update Step 3
```json
{
    "input": {
        "userID": "user-1",
        "author-id": "author-1",
    },
    "output": {
        "userID": "user-1",
        "author-id": "author-1",
        "cachePreviousState": "active",
    }
}
```

Update saga context
```json
{
    "context": {
        "userID": "user-1",
        "authorPreviousStatus": "active",
        "userPreviousState": "active",
        "status": "deleted",
        "author-id": "author-1",
        "cachePreviousState": "active",
    }
}
```