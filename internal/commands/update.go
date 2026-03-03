package commands

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type UpdateCommand struct{}

func NewUpdateCommand() Command {
	return &UpdateCommand{}
}

func (c *UpdateCommand) Name() string { return "update" }
func (c *UpdateCommand) Description() string { return "Update iwashere to latest version" }
func (c *UpdateCommand) Usage() string { return "iwashere update" }
func (c *UpdateCommand) Examples() []string { return []string {"iwashere update"} }


func (c *UpdateCommand) Execute(ctx *Context) error {
    fmt.Println("Checking for updates...")
    
    // Get current executable path
    exe, err := os.Executable()
    if err != nil {
        return err
    }
    
    // Get latest version from GitHub
    latest, err := getLatestVersion()
    if err != nil {
        return err
    }
    
    current := GetVersion()
    if current == latest {
        fmt.Println("Already up to date!")
        return nil
    }
    
    fmt.Printf("Updating from %s to %s...\n", current, latest)
    
    // Download latest version
    tmpDir, err := os.MkdirTemp("", "iwashere-update")
    if err != nil {
        return err
    }
    defer os.RemoveAll(tmpDir)
    
    // Determine architecture
    arch := "x86_64"
    if strings.Contains(runtime.GOARCH, "386") {
        arch = "386"
    }
    
    // Download zip
    url := fmt.Sprintf("https://github.com/Murchoid/iwashere/releases/download/%s/iwashere_%s_Windows_%s.zip", 
        latest, latest, arch)
    
    zipPath := filepath.Join(tmpDir, "iwashere.zip")
    if err := downloadFile(url, zipPath); err != nil {
        return err
    }
    
    // Extract
    if err := unzip(zipPath, tmpDir); err != nil {
        return err
    }
    
    // Find the new executable
    newExe := filepath.Join(tmpDir, "iwashere.exe")
    
    // Replace current executable
    backup := exe + ".bak"
    os.Rename(exe, backup)
    
    if err := copyFile(newExe, exe); err != nil {
        // Restore backup
        os.Rename(backup, exe)
        return err
    }
    
    os.Remove(backup)
    
    fmt.Println("Update complete!")
    return nil
}

func getLatestVersion() (string, error) {
    resp, err := http.Get("https://api.github.com/repos/Murchoid/iwashere/releases/latest")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var release struct {
        TagName string `json:"tag_name"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return "", err
    }
    
    return release.TagName, nil
}

func downloadFile(url, path string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()
    
    _, err = io.Copy(out, resp.Body)
    return err
}

func unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()
    
    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()
        
        path := filepath.Join(dest, f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            out, err := os.Create(path)
            if err != nil {
                return err
            }
            defer out.Close()
            
            _, err = io.Copy(out, rc)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

func copyFile(src, dst string) error {
    in, err := os.Open(src)
    if err != nil {
        return err
    }
    defer in.Close()
    
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()
    
    _, err = io.Copy(out, in)
    return err
}

func init() {
	Register("update", NewUpdateCommand)
}