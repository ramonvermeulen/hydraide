## 🧠 How HydrAIDE Works in Practice – Let’s Start Here

I know you're probably curious about how this works in practice, so let’s start right there.

Since I’m a Go developer, I’ll show you a Go-based example to make it concrete.

--- 

## 🧱 What’s a Swamp?

HydrAIDE stores all data in structures called **Swamps**.
A Swamp is basically a **storage unit** (or container), and each one lives in its own folder with its own file structure.

There are two types of Swamps:

* `Catalog` → for many lightweight key–value pairs
* `Profile` → for complex, structured, entity-level data

Why does this matter?

> Because if you don’t separate your data correctly into Swamps, then **every call will load ALL data into memory** – which is a huge waste.

---

### 🧩 Two Core Data Types: `Catalog` and `Profile`

In HydrAIDE, we primarily work with **two types of data models**:

* **Catalog** – Think of this like a database table where you store *key–value pairs*. A lot of them. Really, a lot.
* **Profile** – This is a **complex structured model**, for example all the data about a user: name, email, permissions, etc.

Here’s the key difference from traditional databases:

> In HydrAIDE, each `Profile` is stored in a **separate container** – almost like each user (or company) has their own database or table.

---

### 🧠 Quick Recap

* **Swamp** = a container that stores data in its own folder
* **Profile** = complex data like a full user or company object
* **Catalog** = key-value data like a company’s permissions list

Behavior:

* When you load a **Profile**, **only that profile** is loaded into memory.
* When you use a **Catalog**, **all its entries** are loaded into memory at once.

So: keep Catalogs **small and focused**. Don’t overload them with unnecessary data!

---

## 📦 Example: A `Catalog` for Registered Companies

Let’s say you want to store the registered companies and their permissions.
We’ll use a `Catalog` to keep track of:

* which companies are registered
* when they registered
* and what permissions they have

In this example, we’re only storing the **company ID** and their **permissions** — nothing more.
The goal is to **quickly check which companies are in the system and what they’re allowed to do**.

---

### 🔧 How to Use It with the SDK

Whether you’re working with a `Catalog` or a `Profile`, the SDK lets you handle the data using plain Go structs.

You can perform CRUD operations using:

* `Load()`
* `Save()`
* `Delete()`
* `Destroy()`

These methods behave like in a regular database, but you can also inject your own logic if needed.

You’ll also use the:

* `RegisterPattern()` function to **register the swamp’s structure** in HydrAIDE (only once, at app start)
* `createModelCatalogName()` to **generate the storage key**, so the SDK knows **which server and folder** to work with

---

### 🧷 One More Thing: Tags Inside the Struct

Notice how we decorate the fields in the Go struct with tags. This is how HydrAIDE knows:

* which field is the key
* which one is the value
* and how to handle them internally

---

### 🚀 Let’s Build: A Company Catalog Swamp

Now we’re ready to create our first real `Catalog` Swamp. One that stores company registrations.

In the methods, we’ll be using SDK calls like `CatalogRead` and others with the `Catalog` prefix to signal what kind of storage we’re dealing with.



```go
package example

import (
	"github.com/hydraide/hydraide/sdk/go/hydraidego"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/name"
...
"time"
)

// ModelCatalogCompany represents the company data, permissions, etc.
// The company ID is the domain name, like trendizz-com, etc.
type ModelCatalogCompany struct {
	// CompanyID is the unique identifier for the company
	CompanyID string `hydraide:"key"`
	// Payload is the value data associated with a company
	Payload *Payload `hydraide:"value"`
	// CreatedAt is the timestamp when the first permission was created
	CreatedAt time.Time `hydraide:"createdAt"`
	// UpdatedAt is the timestamp of the last modification of the company's permissions
	UpdatedAt time.Time `hydraide:"updatedAt"`
}

type Payload struct {
	// Rights contains the permission set assigned to the company
	Rights map[string]interface{}
}

// Load loads a single record from the Hydra database
func (m *ModelCatalogCompany) Load(repo repo.Repo) error {
	// create a timeout-limited context for the Hydra operation
	ctx, cancelFunc := hydraidehelper.CreateHydraContext()
	defer cancelFunc()
	// get the HydraIDE SDK client
	h := repo.GetHydraidego()
	// read data from HydraIDE database into the model
	return h.CatalogRead(ctx, m.createModelCatalogName(), m.CompanyID, m)
}

func (m *ModelCatalogCompany) Save(repo repo.Repo) error {
	// create a timeout-limited context for the Hydra operation
	ctx, cancelFunc := hydraidehelper.CreateHydraContext()
	defer cancelFunc()
	// get the Hydra client
	h := repo.GetHydraidego()
	// save the data into the Hydra database
	_, err := h.CatalogSave(ctx, m.createModelCatalogName(), m)
	return err
}

// Delete removes the model's data from the Hydra database
func (m *ModelCatalogCompany) Delete(repo repo.Repo) error {
	ctx, cancelFunc := hydraidehelper.CreateHydraContext()
	defer cancelFunc()
	// get the Hydra client
	h := repo.GetHydraidego()
	// delete the record from the Hydra database
	return h.CatalogDelete(ctx, m.createModelCatalogName(), m.CompanyID)
}

// Destroy deletes the entire swamp from the Hydra database
// ⚠️ Use for testing purposes only – dangerous in production!
func (m *ModelCatalogCompany) Destroy(repo repo.Repo) error {
	ctx, cancelFunc := hydraidehelper.CreateHydraContext()
	defer cancelFunc()
	// get the Hydra client
	h := repo.GetHydraidego()
	// destroy the entire swamp from the Hydra database
	return h.Destroy(ctx, m.createModelCatalogName())
}

// RegisterPattern registers the swamp pattern in the Hydra database
func (m *ModelCatalogCompany) RegisterPattern(repo repo.Repo) error {

	// get the Hydra client
	h := repo.GetHydraidego()

	ctx, cancelFunc := hydraidehelper.CreateHydraContext()
	defer cancelFunc()

	errorResponses := h.RegisterSwamp(ctx, &hydraidego.RegisterSwampRequest{
		SwampPattern: m.createModelCatalogName(),
		// CloseAfterIdle set to 6 hours – this is a central catalog, often read, so it should remain in memory
		CloseAfterIdle:  time.Second * 21600,
		IsInMemorySwamp: false, // this data needs to persist long-term, so it's not in-memory
		FilesystemSettings: &hydraidego.SwampFilesystemSettings{
			// WriteInterval set to 1 second – new data should be written quickly
			WriteInterval: time.Second * 1,
			MaxFileSize:   32768, // 32KB compressed file chunks on the server – allows for change tracking
		},
	})

	if errorResponses != nil {
		return hydraidehelper.ConcatErrors(errorResponses)
	}

	return nil
}

// createModelCatalogName constructs the name object used to locate the swamp
func (m *ModelCatalogCompany) createModelCatalogName() name.Name {
	return name.New().Sanctuary("trendizz").Realm("companies").Swamp("all")
}

```

Once you've created the model, working with it in Go is incredibly simple.

Need to load a single record? It’s this easy!
It feels a bit like using an ORM, but it’s **much faster**, **more efficient**, and **more flexible**.
It doesn’t tie your hands, but lets you move quickly and freely.

Since you defined `CompanyID` as the key in your model,
all you need to do is set the `CompanyID`, and the `Load()` method will fetch the data into the model for you.

Here’s how it looks:

```go
// ...

companyCatalog := &ModelCatalogCompany{
	CompanyID: "trendizz.com",
}
err := companyCatalog.Load(repo)
if err != nil {
    // handle the error
}
// now you can work directly with the loaded data from the struct

// ...
```

---

## Want to save data? 

This works even if the record already exists — meaning it can also overwrite existing data.

```go
@SuppressLint
// ...
companyCatalog := &ModelCatalogCompany{
    CompanyID: "trendizz.com",
    Payload: &Payload{
        Rights: map[string]interface{}{
            "admin": true,
            "user":  true,
        },
    },
}
err := companyCatalog.Save(repo)
if err != nil {
    // handle the error
}
// ...
```

--- 
## Want to delete a record?

```go
// ...
companyCatalog := &ModelCatalogCompany{
    CompanyID: "trendizz.com",
}
err := companyCatalog.Delete(repo)
if err != nil {
    // handle the error
}
// ...
```

--- 

## What if you need to work with a Profile instead?

In a Profile, we don’t store many entities in one swamp — instead, we store just one, but in very detailed form.

The following code snippet is a real-world example of how we manage data:
In this case, we’re storing AI-generated data for a single domain.

```go
package dizzletcompanyprofile

import (
	"time"

	"github.com/hydraide/hydraide/sdk/go/hydraidego"
	"github.com/hydraide/hydraide/sdk/go/hydraidego/name"
	"github.com/yourproject/internal/hydraidehelper"
	"github.com/yourproject/internal/repo"
)

const (
	RealmNameCompanyProfile = "companyProfile"
)

// For Profile models, you don't need to use decorators,
// but you can optionally use `hydraide:"omitempty"` to reduce file size.
// With this tag, only non-empty fields will be saved — empty fields will be skipped entirely.
type ModelCompanyProfile struct {
	Domain           string
	FoundedAt        time.Time          `hydraide:"omitempty"`
	Employees        int                `hydraide:"omitempty"`
	Departments      []*Department      `hydraide:"omitempty"` // Nested department structures
	TechnologiesUsed []string           `hydraide:"omitempty"` // Tech stack keywords
	Tags             map[string]bool    `hydraide:"omitempty"` // Label flags (e.g. "b2b": true)
	Offices          []*OfficeLocation  `hydraide:"omitempty"` // Physical office locations
	LastUpdated      time.Time          `hydraide:"omitempty"` // Last profile update timestamp
}

type Department struct {
	Name        string   `hydraide:"omitempty"`
	Head        string   `hydraide:"omitempty"`
	HeadEmail   string   `hydraide:"omitempty"`
	TeamMembers []string `hydraide:"omitempty"`
}

type OfficeLocation struct {
	City      string  `hydraide:"omitempty"`
	Country   string  `hydraide:"omitempty"`
	Latitude  float64 `hydraide:"omitempty"`
	Longitude float64 `hydraide:"omitempty"`
	Timezone  string  `hydraide:"omitempty"`
	Active    bool    `hydraide:"omitempty"`
}

func (m *ModelCompanyProfile) Save(r repo.Repo) error {
	ctx, cancel := hydraidehelper.CreateHydraContext()
	defer cancel()

	m.LastUpdated = time.Now()
	return r.GetHydraidego().ProfileSave(ctx, m.createName(m.Domain), m)
}

func (m *ModelCompanyProfile) Load(r repo.Repo) error {
	ctx, cancel := hydraidehelper.CreateHydraContext()
	defer cancel()

	return r.GetHydraidego().ProfileRead(ctx, m.createName(m.Domain), m)
}

func (m *ModelCompanyProfile) Destroy(r repo.Repo) error {
	ctx, cancel := hydraidehelper.CreateHydraContext()
	defer cancel()

	return r.GetHydraidego().Destroy(ctx, m.createName(m.Domain))
}

func (m *ModelCompanyProfile) RegisterPattern(r repo.Repo) error {
	ctx, cancel := hydraidehelper.CreateHydraContext()
	defer cancel()

	return r.GetHydraidego().RegisterSwamp(ctx, &hydraidego.RegisterSwampRequest{
		// The pattern is registered using a wildcard because we want the same settings
		// for all possible domains. The actual domain can be anything, but the config stays the same.
		// Note the use of .Swamp("*") — it means the same setup applies for all domain-specific swamps.
		SwampPattern:    name.New().Sanctuary(SanctuaryName).Realm(RealmNameCompanyProfile).Swamp("*"),
		CloseAfterIdle:  5 * time.Minute,
		IsInMemorySwamp: false,
		FilesystemSettings: &hydraidego.SwampFilesystemSettings{
			WriteInterval: 2 * time.Second,
			MaxFileSize:   2 * 1024 * 1024, // 2MB file chunks
		},
	})
}

func (m *ModelCompanyProfile) createName(domain string) name.Name {
	// It’s important that the .Swamp(domain) call uses the company’s own domain as its identifier.
	// This ensures the profile is stored under a consistent, domain-specific swamp.
	return name.New().Sanctuary(SanctuaryName).Realm(RealmNameCompanyProfile).Swamp(domain)
}

```





