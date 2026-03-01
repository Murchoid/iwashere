// // internal/commands/global.go
package commands

// type GlobalCommand struct{}

// func (c *GlobalCommand) Name() string { return "global" }
// func (c *GlobalCommand) Description() string { return "Manage global (non-project) notes" }

// func (c *GlobalCommand) Execute(ctx *Context) error {
//     // Create global repository
//     globalRepo, err := repository.NewGlobalRepository()
//     if err != nil {
//         return err
//     }

//     if len(ctx.Args) == 0 {
//         // List global notes
//         notes, _ := globalRepo.ListNotes(nil)
//         for _, note := range notes {
//             fmt.Printf("[%s] %s\n",
//                 note.CreatedAt.Format("2006-01-02"),
//                 note.Message)
//         }
//         return nil
//     }

//     // Add global note
//     return globalRepo.AddGlobalNote(strings.Join(ctx.Args, " "))
// }
