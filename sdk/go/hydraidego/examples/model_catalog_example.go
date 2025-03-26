//go:build ignore
// +build ignore

// This file provides a detailed example of a catalog-style model used with CatalogCreate().
// It explains required fields, supported types, optional metadata, and best practices.

package examples

import "time"

// Example: CatalogCreditLog — a catalog model for logging credit operations per user.
//
// This struct demonstrates how to define a model for CatalogCreate.
// Each field uses `hydraide` tags to indicate its role within the KeyValuePair.
// All values will be transformed into HydrAIDE-compatible binary format at runtime.

type CatalogCreditLog struct {
	// 🔑 REQUIRED
	// This will be used as the Treasure key.
	// Must be a non-empty string.
	UserUUID string `hydraide:"key"`

	// 📦 OPTIONAL — The value of the Treasure.
	// Can be:
	// - Primitive types: string, bool, int8–64, uint8–64, float32, float64
	// - Structs (encoded via GOB)
	// - Pointer to struct (also GOB-encoded)
	// ⚠️ Use the SMALLEST type possible for space efficiency.
	Log struct {
		Amount   int16  // ✅ Small integer: better memory & disk usage than int
		Reason   string // Reason for the credit log (e.g. "bonus")
		Currency string // Currency ISO code (e.g. "HUF", "EUR")
	} `hydraide:"value"`

	// ⏳ OPTIONAL
	// The logical expiration of this entry.
	// Must be a valid non-zero time.Time.
	ExpireAt time.Time `hydraide:"expireAt"`

	// 🧾 OPTIONAL METADATA — useful for tracking/audit purposes
	// If omitted, these fields will not be included in the stored record.

	CreatedBy string    `hydraide:"createdBy"` // Who created the record
	CreatedAt time.Time `hydraide:"createdAt"` // When it was created
	UpdatedBy string    `hydraide:"updatedBy"` // Who last updated it
	UpdatedAt time.Time `hydraide:"updatedAt"` // When it was last updated
}
