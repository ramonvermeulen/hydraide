# `name` Package – Deterministic Naming in HydrAIDE

---

📘 This is part of the HydrAIDE Go SDK documentation.
Full SDK documentation: https://github.com/hydraide/hydraide/blob/main/docs/sdk/go/README.md

---

The `name` package defines how Swamps are **logically grouped and uniquely identified** using a three-level hierarchy:

```
Sanctuary / Realm / Swamp
```

This naming convention is foundational in HydrAIDE. It is used to:
- Organize data into **meaningful, navigable structures**
- Enable **stateless, hash-based routing** in distributed environments
- Ensure consistent, reproducible folder or server assignments

By enforcing a consistent structure for all Swamps, HydrAIDE eliminates the need for central registries or coordinators. This makes it possible to locate, access, and route data **with zero external metadata**.

---

## ✨ Example

```go
n := name.New().
    Sanctuary("users").
    Realm("profiles").
    Swamp("alice123")

fmt.Println(n.Get())                // "users/profiles/alice123"
fmt.Println(n.GetServerNumber(100)) // → e.g. 42
```

---

## 📚 Functions

```go
New() Name
```
Creates an empty Name instance to begin building the hierarchy.

```go
Sanctuary(id string) Name
```
Sets the top-level logical grouping (e.g. "users", "domains", "products").

```go
Realm(name string) Name
```
Sets a mid-level scope to organize within a Sanctuary.

```go
Swamp(name string) Name
```
Defines the final segment — the unique Swamp identifier.

```go
Get() string
```
Returns the full canonical path ("sanctuary/realm/swamp").
🔒 Internal use only.

```go
GetServerNumber(allServers int) uint16
```
Returns the 1-based server number where this Name belongs.
🔒 Internal use only.

```go
Load(path string) Name
```
Reconstructs a Name from a path string.
🔒 Internal use only.

---

## 🧠 Best Practices

### 🧩 Real-World Grouping Patterns

Here are some practical examples of how to structure Names in real-world use cases:

```go
// User-related data
name.New().Sanctuary("users").Realm("profiles").Swamp("u123")
name.New().Sanctuary("users").Realm("sessions").Swamp("u123-active")

// Product catalog
name.New().Sanctuary("products").Realm("info").Swamp("sku-8932")
name.New().Sanctuary("products").Realm("inventory").Swamp("warehouse-5")

// Analytics logs
name.New().Sanctuary("logs").Realm("search").Swamp("2025-03-25")
name.New().Sanctuary("logs").Realm("clicks").Swamp("session-8fk29x")

// Platform usage
name.New().Sanctuary("system").Realm("metrics").Swamp("cpu")
name.New().Sanctuary("system").Realm("events").Swamp("node-42")
```

> 💡 Tip: Reuse consistent Sanctuary/Realm schemes across modules to keep your data model predictable and scalable.

- Names should use short, compact, readable words that clearly represent logical units.
- Avoid ambiguous or overly technical terms. Prioritize clarity and consistency.
- For a deeper understanding of naming strategies, conventions, and real-world patterns, see the extended guide: [Naming Convention](https://github.com/hydraide/hydraide/blob/main/docs/thinking-in-hydraide/naming-convention.md)

- Always group Swamps under a logical Sanctuary/Realm structure
- Avoid using raw strings like "flat/swampX" — use semantic nesting instead
- Never call `Get()` or `GetServerNumber()` directly from application code — let the SDK route based on Name

---

## 📄 License Notice
This document is part of the HydrAIDE knowledge base and is licensed under a **custom restrictive license**.  
You may not use its contents to build or assist in building alternative engines, architectures, or competing systems.  
See: [LICENSE.md](https://github.com/hydraide/hydraide/blob/main/LICENSE.md)
