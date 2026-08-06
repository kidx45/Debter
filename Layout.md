I actually think you're making the right decision here.

The mistake wasn't that you stopped after the tutorial. The mistake would be trying to jump straight into "distributed AI F1 telemetry on AWS with Rust and Go" before you're comfortable building a complete application from memory.

You've already noticed the real problem:

> "I can read the code, but I can't create it."

Those are two very different skills.

Reading code is passive.
Building is active.

The only thing that fixes that is building over and over.

---

# Here's the rule we're going to follow.

For the next **12 weeks**, you are **not allowed** to build anything "cool."

No AI.
No Kubernetes.
No AWS.
No fancy architecture.
No microservices.

Only complete applications.

Why?

Because companies don't hire people who know Redis.

They hire people who can finish software.

---

# Your Daily Rules

These are non-negotiable.

* **2 hours** building
* **30 minutes** learning (video, docs, book)
* **No YouTube Shorts before the project is finished**
* Every project must be pushed to GitHub every day
* Every project must have a README
* Every feature must be committed separately

Example:

```
git commit -m "create database schema"

git commit -m "implement authentication"

git commit -m "add login endpoint"

git commit -m "write unit tests"

git commit -m "dockerize application"
```

Treat your Git history like a story someone else could read.

---

# This Week's Project

You're restarting.

So the project should touch almost everything you've learned.

Not everything.

Just enough.

---

# Project

## Expense Tracker API

Not because it's exciting.

Because it forces you to use backend fundamentals.

Imagine you're building the backend for a Flutter app.

---

## Features

Authentication

* Register
* Login
* Refresh Token

Users

* View profile
* Update profile

Categories

* Food
* Transport
* Shopping
* Salary
* etc.

Expenses

* Create
* Edit
* Delete
* List

Income

* Create
* Edit
* Delete

Dashboard

Return

```
Total Income

Total Expense

Current Balance

Expense by Category
```

---

# Technologies

You MUST use

✅ Go

✅ Gin

✅ PostgreSQL

✅ sqlc

✅ golang-migrate

✅ Docker

✅ Docker Compose

✅ JWT (or PASETO if you remember)

✅ Viper

✅ Validator

---

Optional

Redis ❌

gRPC ❌

AWS ❌

RabbitMQ ❌

Don't touch them.

---

# Folder Structure

```
expense-tracker

cmd/
    main.go

internal/

    api/

    db/

        migration/

        query/

        sqlc/

    token/

    util/

    config/

    middleware/

    validator/

    service/

    repository/

docs/

Dockerfile

docker-compose.yml

Makefile

README.md
```

---

# Database

Tables

```
users

categories

expenses

income
```

Relationships

```
User

    |

    |---- Expenses

    |

    |---- Income

    |

    |---- Categories
```

---

# Endpoints

Authentication

```
POST /users

POST /login

POST /refresh
```

User

```
GET /users/me

PATCH /users/me
```

Categories

```
POST /categories

GET /categories

PATCH /categories/:id

DELETE /categories/:id
```

Expenses

```
POST /expenses

GET /expenses

GET /expenses/:id

PATCH /expenses/:id

DELETE /expenses/:id
```

Income

Same CRUD.

Dashboard

```
GET /dashboard
```

---

# Testing

Write tests.

Not because they're fun.

Because otherwise you won't remember how.

Test

* sqlc queries
* token maker
* validators
* handlers
* services

Use testify.

---

# Docker

By the end

```
docker compose up
```

should start

```
Postgres

API
```

No manual setup.

---

# Git Workflow

Every feature

```
git checkout -b feature/login

git commit

git push
```

Merge later.

---

# Schedule

### Day 1

Database

sqlc

migrations

docker

---

### Day 2

Authentication

JWT/PASETO

middleware

---

### Day 3

Expense CRUD

Income CRUD

---

### Day 4

Dashboard

Validation

Testing

---

### Day 5

Refactor

README

Bug fixes

---

### Day 6

Finish remaining tests

Polish

---

### Day 7

Deploy locally

Ask a friend to test it using Postman or Bruno

---

# After This

Only after you've finished this project do we move on to:

1. Library Management API
2. URL Shortener
3. Chat API
4. Job Board
5. File Storage API
6. E-commerce Backend
7. Real-time Chat
8. gRPC rewrite
9. Redis caching
10. Event-driven version

Notice the pattern.

Every project adds **one** new concept.

Not ten.

---

# One more thing

I noticed something from our conversations over the last few weeks.

Every few days, your attention shifts to something shiny:

* "Should I learn Rust?"
* "What about AI?"
* "Maybe AWS?"
* "Maybe Kubernetes?"
* "Could I build this huge F1 system?"

Curiosity is a strength, but if it constantly changes your direction, you end up with broad exposure and very little depth.

For the next **7 days**, your only goal is this:

> **Finish one backend project.**

Not "start."

Not "redesign."

Not "research."

**Finish it.**

When you finish it, you'll have something much more valuable than another completed tutorial: a project you assembled yourself. That's the skill that compounds over time.

And here's the commitment I'd like you to make to yourself:

> **For the next 12 weeks, I will finish one project every week before chasing a new technology.**

If you can keep that promise, you'll be surprised how much more confident you'll feel reading *and* writing Go code by the end of those three months.
