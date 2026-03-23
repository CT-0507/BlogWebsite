# **🧾 Author Module — REST API Documentation**

🔧 Base Info

### **Base URL:**

```
/api/v1/authors
```

### **Auth:**

- Required for write operations
- Optional for public reads

#

## **1️⃣ Author Identity & Profile**

### **✅ Create Author Profile**

```
POST /api/v1/authors
```

**Request**

```json
{
  "displayName": "John Doe",
  "bio": "Tech writer",
  "avatarUrl": "https://...",
  "socialLinks": {
    "twitter": "...",
    "github": "..."
  }
}
```

**Response**

```json
{
  "id": "author_123",
  "slug": "john-doe",
  "displayName": "John Doe",
  "bio": "Tech writer",
  "avatarUrl": "...",
  "createdAt": "2026-03-18T10:00:00Z"
}
```

### **✅ Get Author by ID**

```
POST /api/v1/authors
```

**Request**

```json
{
  "displayName": "John Doe",
  "bio": "Tech writer",
  "avatarUrl": "https://...",
  "socialLinks": {
    "twitter": "...",
    "github": "..."
  }
}
```

**Response**

```json
{
  "id": "author_123",
  "slug": "john-doe",
  "displayName": "John Doe",
  "bio": "Tech writer",
  "avatarUrl": "...",
  "createdAt": "2026-03-18T10:00:00Z"
}
```

### **✅ Get Author by ID**

```
GET /api/v1/authors/{authorId}
```

**Response**

```json
{
  "id": "author_123",
  "slug": "john-doe",
  "displayName": "John Doe",
  "bio": "Tech writer",
  "avatarUrl": "...",
  "stats": {
    "followers": 120
  }
}
```

### **✅ Get Author by Slug**

```
GET /api/v1/authors/slug/{slug}
```

### **✅ List Authors (Filter/Search)**

```
GET /api/v1/authors
```

**Query Params**

```
?keyword=john&category=tech&page=1&limit=20
```

**Response**

```json
{
  "data": [
    {
      "id": "author_123",
      "slug": "john-doe",
      "displayName": "John Doe",
      "avatarUrl": "..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100
}
```

### **✅ Update Author Profile**

```
PATCH /api/v1/authors/{authorId}
```

**Request**

```json
{
  "displayName": "John Updated",
  "bio": "New bio"
}
```

### **✅ Delete Author (Soft Delete)**

```
DELETE /api/v1/authors/{authorId}
```

**Response**

```json
{
  "success": true
}
```

### **✅ Update Author Slug**

```
PATCH /api/v1/authors/{authorId}/slug
```

**Response**

```json
{
  "slug": "new-john-doe"
}
```

##

## **2️⃣ Follow System**

### **✅ Follow Author**

```
POST /api/v1/authors/{authorId}/follow
```

**Response**

```json
{
  "message": "Success"
}
```

### **✅ Unfollow Author**

```
DELETE /api/v1/authors/{authorId}/follow
```

**Response**

```json
{
  "message": "Success"
}
```

### **✅ Get Author Followers**

```
GET /api/v1/authors/{authorId}/followers
```

Query

```
?page=1&limit=20
```

### **✅ Get Followed Authors (by user)**

```
GET /api/v1/me/following/authors
```

##

## **3️⃣ Author Profile Presentation (Featured Blogs)**

### **✅ Set Featured Blogs**

```
PUT /api/v1/authors/{authorId}/featured-blogs
```

**Request**

```json
{
  "blogIds": ["blog_1", "blog_2", "blog_3"]
}
```

Rules

Max 10 blogs

Order matters

Blogs validated via Blog module

### **✅ Get Featured Post IDs**

```
GET /api/v1/authors/{authorId}/featured-posts
```

**Response**

```json
{
  "postIds": ["post_1", "post_2", "post_3"]
}
```

##

## **4️⃣ Author Stats & Metadata**

### **✅ Get Author Stats**

```
GET /api/v1/authors/{authorId}/stats
```

**Response**

```javascript
{
  "followers": 120,
  "blogs": 45
}
```

## **5️⃣ Author Discovery**

### **✅ Search Authors**

```
GET /api/v1/authors/search
```

**Query**

```
?q=john&page=1&limit=20

```

### **✅ Get Trending Authors (Optional)**

```
GET /api/v1/authors/trending
```

### Final

```
/authors
/authors/{id}
/authors/slug/{slug}
/authors/{id}/follow
/authors/{id}/followers
/authors/{id}/featured-blogs
/authors/{id}/stats

/me/following/authors
```

##

# Rules

Author API:

✅ returns author data

✅ returns blog IDs only
