# Library Management REST API with Reservation and Checkout Workflow

## Project Overview

This project is a RESTful backend service built using Go (Golang) to manage a library system. It supports complete book lifecycle management including multi-copy tracking, student checkout workflow, return handling, fine calculation, and FIFO-based reservation management.

The system models real-world transactional behavior and ensures consistent state transitions across operations.

---

## Tech Stack

- Language: Go (Golang)
- Web Framework: Gin
- ORM: GORM
- Database: SQLite
- Architecture: REST API

---

## Core Features

### Book Management
- Add books with multiple copies
- Track each copy independently
- View all books and their copy status

### User Management
- Register users (student or librarian)
- Only students are allowed to checkout or reserve books
- View all registered users

### Checkout Workflow
- Assigns only available copies
- Marks selected copy as "checked_out"
- Creates a checkout record
- Automatically sets due date (7 days)

### Return Workflow
- Marks checkout record as returned
- Changes book copy status back to "available"
- Calculates fine if overdue (₹10 per day)
- Automatically assigns book to next reserved user (FIFO)

### Reservation System (FIFO Queue)
- If no copies are available, students can reserve the book
- Reservations are processed in First-In-First-Out order
- Upon return, the earliest active reservation is fulfilled automatically

---

## System Workflow

1. A book is added with multiple copies.
2. Students checkout available copies.
3. If no copies remain, students may reserve the book.
4. When a copy is returned:
   - The system checks the reservation queue.
   - The earliest reservation is fulfilled.
   - A new checkout record is created automatically.

---

## Fine Calculation

- Due Date = Checkout Date + 7 days
- Fine = ₹10 per day after due date
- No fine if returned on or before due date

---

## API Endpoints

### Books
- POST /books – Add a new book with copies
- GET /books – View all books

### Users
- POST /users – Register a new user
- GET /users – View all users

### Checkout
- POST /checkout – Borrow a book
- GET /checkouts – View all checkout records

### Reservation
- POST /reserve – Reserve a book
- GET /reservations – View all reservations

### Return
- POST /return – Return a borrowed book

---

## Database Design

Entities:
- Book
- BookCopy
- User
- Checkout
- Reservation

Relationships:
- One Book has many BookCopies
- One User can have many Checkouts
- One Book can have many Reservations
- One Checkout references one BookCopy

Each book copy is tracked independently to support multi-copy management.

---

## Concurrency Considerations

If multiple users attempt to checkout the last available copy simultaneously, a race condition may occur.

In production systems, this is typically handled using:
- Database transactions
- Row-level locking (e.g., SELECT FOR UPDATE)
- Optimistic locking strategies
- Distributed locking mechanisms

Due to SQLite limitations and project scope, this implementation focuses on logical correctness while acknowledging concurrency considerations.

---

## How to Run

1. Install Go.
2. Navigate to the project folder.
3. Run the following commands:

```
go mod tidy
go run main.go
```

The server will start locally at:

http://localhost:8080

Note: This URL works only after running `go run main.go` on your local machine. The API is not deployed online, so accessing this link from GitHub without running the server will show "Unable to connect".

---

## Assumptions

- Only students can checkout or reserve books.
- Fine rate is fixed at ₹10 per day.
- SQLite is used for simplicity.
- Authentication and authorization layers are not implemented.
