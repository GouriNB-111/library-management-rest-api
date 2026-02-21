# AI Prompts Used During Development

This project was developed with assistance from an AI tool.  
Below is a summary of the prompts used during development.

---

## 1. Initial Project Setup

Prompt:
"Build a Library Management REST API in Go using Gin and GORM with SQLite."

Purpose:
To generate the base backend structure including models, database connection, and server setup.

---

## 2. Multi-Copy Book Management

Prompt:
"Implement book management where a single book can have multiple physical copies and track each copy separately."

Purpose:
To support real-world multi-copy tracking instead of single inventory count.

---

## 3. Checkout Workflow

Prompt:
"Add checkout functionality where only students can borrow books, mark copy as checked_out, and set a due date."

Purpose:
To implement borrowing workflow and state transition logic.

---

## 4. Return Workflow and Fine Calculation

Prompt:
"Add return API to mark checkout as returned, make copy available again, and calculate fine if overdue."

Purpose:
To implement lifecycle completion and fine logic.

---

## 5. FIFO Reservation System

Prompt:
"Implement a reservation system with FIFO logic and auto-assign the book to the next reserved user when returned."

Purpose:
To handle waiting queue logic when no copies are available.

---

## 6. Concurrency Awareness

Prompt:
"Explain how to handle race conditions when two users attempt to checkout the same book at the same time."

Purpose:
To document production-level concurrency considerations.

---

All AI-assisted prompts were reviewed, understood, and modified where necessary to ensure correctness and alignment with project requirements.