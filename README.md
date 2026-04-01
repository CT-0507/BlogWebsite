# 🧩 Modular Monolith Blog Platform

A scalable, production-ready blog platform built with a modular monolith architecture designed to be microservice-ready. This project demonstrates modern backend patterns, event-driven design, and cloud-native deployment practices.

## 🚀 Overview

This project is a full-stack blogging platform where users can create and manage author profiles, publish blogs, interact through comments, and follow authors. The system is designed with scalability and maintainability in mind, leveraging a modular monolith approach that can evolve into microservices.

## 🏗️ Architecture

Modular Monolith (Microservice-Ready)

The system is built as a modular monolith, where each domain (user, author, blog, comment, etc.) is isolated into independent modules with clear boundaries.

### Why this approach?

1.  Easier development and deployment compared to microservices
2.  Avoids distributed system complexity early on
3.  Clear separation of concerns
4.  Ready to extract into microservices when scaling demands it

## Tech Stack

### Frontend

- React
- TanStack Query + Axios
- Material UI

#### Why?

- React: Component-based architecture for scalable UI development
- TanStack Query: Efficient server-state management (caching, retries, background - sync)
- Axios: Simple and flexible HTTP client
- Material UI: Fast, consistent, and accessible UI components

### Backend

- Golang + Gin
- PostgreSQL
- sqlc
- In-memory Event Bus
- Outbox Pattern
- Saga Pattern

#### Why?

- Golang: High performance, concurrency support, and simplicity
- Gin: Lightweight and fast HTTP framework
- PostgreSQL: Reliable, ACID-compliant relational database
- sqlc:
- Generates type-safe SQL code
- Avoids ORM overhead and hidden behavior
- Ensures full control over queries and performance
- Event-driven patterns:
- In-memory event bus: Decouples modules without network overhead
- Outbox pattern: Guarantees reliable event publishing
- Saga pattern: Handles distributed transactions and consistency across modules

### Infrastructure & DevOps

- Docker
- Nginx (Load Balancer, Rate Limiter, API Gateway) (not yet implement)
- GitHub Actions (CI/CD) (not yet implement)
- Prometheus (Monitoring) (not yet implement)
- Redis (not yet implement)
- AWS EC2 (E3 instance) (not yet implement)
- AWS S3 (Storage) (not yet implement)
- AWS PostgreSQL (RDS) (not yet implement)

#### Why?

- Docker: Consistent environments and easy deployment
- Nginx:
  - Load balancing across backend instances
  - Rate limiting for API protection
  - Acts as API gateway entry point
- GitHub Actions:
  - Automated build, test, and deployment pipelines
  - Ensures code quality and faster iteration
- Prometheus:
  - Real-time monitoring and metrics collection
  - Observability for system performance
- Redis:
  - Caching layer for performance optimization
  - Supports future real-time features (e.g., chat)
- AWS:
  - EC2: Flexible compute hosting
  - S3: Scalable object storage for assets
  - RDS PostgreSQL: Managed, reliable database

## 📊 Load Testing

The system is tested under load to validate:

- Scalability
- Throughput
- System reliability under stress

## ✨ Features

### 👤 User & Author System

Users can create and manage:

- User profile
- Author profile
- Ability to switch identity when interacting (user vs author)

### ✍️ Blogging

Authors can:

- Create, update, and delete blogs
- Manage comments on their blogs

### 💬 Comments

- Users can comment on blogs
- Users with author profiles can choose which identity to comment as

### ⭐ Author Ranking System

Authors are ranked based on:

- Likes / dislikes
- Number of blogs
- Follower count

### 🔔 Follower & Notification System

- Users can follow authors
- Followers receive notifications when authors publish new blogs

### 🛠️ Admin Dashboard

Admin capabilities include:

- Manage Saga states
- Monitor and handle Dead Letter Queue (DLQ)
- Full CRUD operations on:
  - Users
  - Authors
  - Blogs
  - Comments

### 💬 Chat System (Planned)

- Real-time chat between users (future enhancement)
- Likely to leverage Redis and event-driven architecture

### 📦 Key Design Highlights

- Event-Driven Architecture inside a monolith
- Strong consistency with Saga orchestration
- High performance via sqlc (no ORM overhead)
- Cloud-native and containerized
- Observability-first design

### 🧪 Future Improvements

- Extract modules into independent microservices
- Implement real-time chat with WebSockets
- Introduce message broker (e.g., Kafka) for distributed event bus
- Enhance recommendation and ranking algorithms

### 📌 Summary

This project demonstrates how to:

- Build a scalable backend using modular monolith architecture
- Apply advanced distributed system patterns without full microservice complexity
- Design a system that is production-ready and evolution-friendly
