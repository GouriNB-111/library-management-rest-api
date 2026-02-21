# System Design Document  
Library Management REST API with Reservation and Checkout Workflow

---

## 1. Design Objectives

The primary objective of this system is to model a real-world library workflow using a RESTful backend service in Go.

Key goals:

- Support multi-copy book management
- Maintain clear state transitions
- Ensure consistent checkout and return workflow
- Implement FIFO-based reservation handling
- Maintain data integrity across operations
- Keep the system simple, modular, and extensible

---

## 2. Architectural Approach

This system follows a layered REST API design:

Client → HTTP Routes (Gin) → Business Logic → GORM ORM → SQLite Database

The application is implemented as a backend service without a frontend layer.  
It exposes REST endpoints that can be consumed by any client (Postman, frontend UI, etc.).

---

## 3. Entity Design

### 3.1 Book

Represents a logical book entry.

Fields:
- ID
- Title
- Author
- ISBN

A Book does not represent a physical unit. Instead, physical availability is handled by BookCopy.

---

### 3.2 BookCopy

Represents an individual physical copy of a book.

Fields:
- ID
- BookID
- Status ("available", "checked_out")

Each book may have multiple copies.
Copy-level tracking enables accurate availability management.

---

### 3.3 User

Represents a system user.

Fields:
- ID
- Name
- Role ("student", "librarian")

Only users with role "student" are allowed to checkout or reserve books.

---

### 3.4 Checkout

Represents a borrowing transaction.

Fields:
- ID
- UserID
- BookCopyID
- DueDate
- Returned (boolean)

This ensures full lifecycle tracking of borrowed books.

---

### 3.5 Reservation

Represents a waiting queue entry for a book.

Fields:
- ID
- UserID
- BookID
- CreatedAt
- Active (boolean)

Reservations are processed in FIFO order using CreatedAt.

---

## 4. State Transitions

The system enforces strict state transitions at the BookCopy level.

Available → Checked Out → Available

Checkout:
- Copy status changes from "available" to "checked_out"
- Checkout record created

Return:
- Checkout marked as returned
- Copy status changes to "available"
- Reservation queue checked
- If reservation exists → auto checkout

---

## 5. Checkout Workflow Design

Steps:

1. Validate user existence.
2. Ensure user role is "student".
3. Query for available book copy.
4. If available:
   - Mark copy as "checked_out"
   - Create checkout record
   - Set due date (7 days)

If no copies available:
- Inform user to reserve.

This ensures no over-allocation of book copies.

---

## 6. Reservation Workflow Design

If all copies are checked out:

1. User submits reservation.
2. Reservation stored with timestamp.
3. Active flag set to true.

Upon book return:

1. System checks active reservations.
2. Orders by CreatedAt (FIFO).
3. Selects earliest reservation.
4. Auto assigns returned copy.
5. Creates new checkout record.
6. Marks reservation inactive.

This guarantees fairness and deterministic allocation.

---

## 7. Fine Calculation Logic

Fine rules:

- Due date = Checkout Date + 7 days
- Fine = ₹10 per day after due date
- If returned on or before due date → no fine

Fine is computed dynamically during return operation.

---

## 8. Validation Strategy

The system includes the following validations:

- User must exist
- Book must exist
- Only students can checkout/reserve
- Cannot return already returned checkout
- Cannot checkout if no copies available

These validations ensure data consistency.

---

## 9. Concurrency Considerations

Potential race condition:

If two users attempt to checkout the last available copy simultaneously, both may read the copy as available before it is updated.

In production systems, this would be handled using:

- Database transactions
- Row-level locking (SELECT FOR UPDATE)
- Optimistic locking strategies
- Distributed locking mechanisms

SQLite provides limited concurrency control.  
For this implementation, the focus is logical correctness under normal operation, while acknowledging concurrency risks in high-traffic environments.

---

## 10. Scalability Considerations

For production-grade deployment:

- Replace SQLite with PostgreSQL or MySQL
- Introduce transaction management
- Implement authentication and authorization
- Separate service layer from route layer
- Introduce logging and monitoring

The current implementation is structured to allow these enhancements without major redesign.

---

## 11. Design Principles Followed

- Separation of concerns
- Clear entity modeling
- Deterministic state transitions
- RESTful endpoint design
- Minimal but sufficient validation
- FIFO fairness in reservation handling

---

## 12. Conclusion

This system models a realistic library workflow with:

- Multi-copy tracking
- Checkout lifecycle management
- Automated reservation queue handling
- Fine calculation
- State consistency

The design prioritizes clarity, correctness, and extensibility while remaining lightweight and easy to understand.