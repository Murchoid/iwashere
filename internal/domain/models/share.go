package models

import (
	"time"
)

// EncryptedPayload is what gets stored in the git-tracked shared folder
type EncryptedPayload struct {
	NoteID        string    `json:"note_id"`
	EncryptedNote string    `json:"encrypted_note"` // Base64 encrypted note JSON
	EncryptedKey  string    `json:"encrypted_key"`  // Base64 encrypted symmetric key
	IV            string    `json:"iv"`             // Base64 IV
	SharedBy      string    `json:"shared_by"`
	SharedAt      time.Time `json:"shared_at"`
	NotePreview   string    `json:"note_preview"` // Plaintext preview for listing
}

// SharedNote is the actual note content after decryption
type SharedNote struct {
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	Tags        []string  `json:"tags,omitempty"`
	Branch      string    `json:"branch,omitempty"`
	SessionName string    `json:"session_name,omitempty"`
	Author      string    `json:"author"`
}

// EncryptedPayloadIndex helps list available shares
type EncryptedPayloadIndex struct {
	Shares map[string][]string `json:"shares"` // recipient -> list of note IDs
}
