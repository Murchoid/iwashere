// // internal/repository/global.go
package repository

// import (
// 	"os"
// 	"path/filepath"
// 	"time"

// 	"githum.com/Murchoid/iwashere/internal/repository/jsonRepo"
// )

// type GlobalRepository struct {
//     *jsonRepo.JSONRepository // Embed your JSON repo
// }

// func NewGlobalRepository() (*GlobalRepository, error) {
//     globalDir := filepath.Join(utils.GetConfigDir(), "notes")

//     // Create global notes directory
//     if err := os.MkdirAll(globalDir, 0755); err != nil {
//         return nil, err
//     }

//     return &GlobalRepository{
//         &jsonRepo.JSONRepository{
//             NotesBasePath: globalDir,
//         },
//     }, nil
// }

// // Add a global note (not tied to any project)
// func (r *GlobalRepository) AddGlobalNote(msg string) error {
//     note := &models.Note{
//         Message:   msg,
//         CreatedAt: time.Now(),
//         // No ProjectPath, No Branch - it's global!
//     }
//     return r.SaveNote(note)
// }
