package commands

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"githum.com/Murchoid/iwashere/internal/utils"
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
    if len(ctx.Args) > 0 {
        fmt.Println("Unrecognized arguments")
        fmt.Println()
        utils.PrintCommandHelp(c.Name(), c.Description(), c.Usage(), c.Examples())
        return nil
    }
    
    fmt.Println("Checking for updates...")
    
    // Get current executable path
    exe, err := os.Executable()
    if err != nil {
        return fmt.Errorf("failed to get executable path: %w", err)
    }
    
    // Get latest version from GitHub
    latest, err := getLatestVersion()
    if err != nil {
        return fmt.Errorf("failed to check latest version: %w", err)
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
        return fmt.Errorf("failed to create temp directory: %w", err)
    }
    defer os.RemoveAll(tmpDir)
    
    // Get platform-specific file info
    fileInfo, err := getPlatformFileInfo(latest)
    if err != nil {
        return err
    }
    
    fmt.Printf("Downloading for %s/%s...\n", fileInfo.goos, fileInfo.arch)
    fmt.Printf("URL: %s\n", fileInfo.url)
    
    // Download the archive
    archivePath := filepath.Join(tmpDir, "iwashere"+fileInfo.archiveExt)
    if err := downloadFile(fileInfo.url, archivePath); err != nil {
        return fmt.Errorf("failed to download: %w", err)
    }
    
    // Verify it's a valid archive
    if err := verifyArchive(archivePath, fileInfo.archiveExt); err != nil {
        return fmt.Errorf("downloaded file is not valid: %w", err)
    }
    
    // Extract the binary
    binaryPath, err := extractBinary(archivePath, tmpDir, fileInfo)
    if err != nil {
        return fmt.Errorf("failed to extract: %w", err)
    }
    
    // Make it executable
    if err := os.Chmod(binaryPath, 0755); err != nil {
        return fmt.Errorf("failed to set permissions: %w", err)
    }
    
    // Replace current executable
    backup := exe + ".bak"
    if err := os.Rename(exe, backup); err != nil {
        return fmt.Errorf("failed to create backup: %w", err)
    }
    
    if err := copyFile(binaryPath, exe); err != nil {
        // Restore backup
        os.Rename(backup, exe)
        return fmt.Errorf("failed to replace binary: %w", err)
    }
    
    os.Remove(backup)
    os.RemoveAll(tmpDir)
    
    fmt.Println("Update complete!")
    fmt.Printf("Run 'iwashere --version' to verify\n")
    return nil
}

type PlatformFileInfo struct {
    goos       string
    osName     string  
    arch       string
    binaryName string
    archiveExt string
    url        string
}

func getPlatformFileInfo(version string) (*PlatformFileInfo, error) {
    info := &PlatformFileInfo{
        goos: runtime.GOOS,
    }
    
    // Map Go arch to GoReleaser arch naming
    switch runtime.GOARCH {
    case "amd64":
        info.arch = "amd64"  // GoReleaser uses "amd64", not "x86_64"
    case "386":
        info.arch = "386"     // GoReleaser uses "386", not "i386"
    case "arm64":
        info.arch = "arm64"   // GoReleaser uses "arm64"
    default:
        info.arch = runtime.GOARCH
    }
    
    // Set OS-specific values
    switch runtime.GOOS {
    case "windows":
        info.binaryName = "iwashere.exe"
        info.archiveExt = ".zip"
        info.osName = "windows"  // lowercase!
    case "linux":
        info.binaryName = "iwashere"
        info.archiveExt = ".tar.gz"
        info.osName = "linux"    // lowercase!
    case "darwin":
        info.binaryName = "iwashere"
        info.archiveExt = ".tar.gz"
        info.osName = "darwin"   // lowercase! (not macOS)
    default:
        return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
    }
    
    // Format: iwashere_0.2.0_linux_amd64.tar.gz
    cleanVersion := strings.TrimPrefix(version, "v")
    
    info.url = fmt.Sprintf(
        "https://github.com/Murchoid/iwashere/releases/download/%s/iwashere_%s_%s_%s%s",
        version,           // v0.2.0
        cleanVersion,      // 0.2.0
        info.osName,       // linux (lowercase!)
        info.arch,         // amd64
        info.archiveExt,    // .tar.gz
    )
    
    return info, nil
}

func verifyArchive(path, ext string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Read first few bytes to check magic numbers
    header := make([]byte, 10)
    if _, err := file.Read(header); err != nil {
        return err
    }
    
    switch ext {
    case ".zip":
        // ZIP files start with "PK"
        if header[0] != 'P' || header[1] != 'K' {
            return fmt.Errorf("not a valid zip file (starts with %x)", header[:2])
        }
    case ".tar.gz":
        // gzip files start with 0x1F 0x8B
        if header[0] != 0x1F || header[1] != 0x8B {
            return fmt.Errorf("not a valid gzip file (starts with %x)", header[:2])
        }
    }
    
    return nil
}

func extractBinary(archivePath, destDir string, info *PlatformFileInfo) (string, error) {
    switch info.archiveExt {
    case ".zip":
        return extractFromZip(archivePath, destDir, info.binaryName)
    case ".tar.gz":
        return extractFromTarGz(archivePath, destDir, info.binaryName)
    default:
        return "", fmt.Errorf("unsupported archive format: %s", info.archiveExt)
    }
}

func extractFromZip(zipPath, destDir, binaryName string) (string, error) {
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return "", err
    }
    defer r.Close()
    
    for _, f := range r.File {
        // Look for the binary (might be in a subdirectory)
        if strings.HasSuffix(f.Name, binaryName) && !f.FileInfo().IsDir() {
            rc, err := f.Open()
            if err != nil {
                return "", err
            }
            defer rc.Close()
            
            binaryPath := filepath.Join(destDir, binaryName)
            out, err := os.Create(binaryPath)
            if err != nil {
                return "", err
            }
            defer out.Close()
            
            if _, err := io.Copy(out, rc); err != nil {
                return "", err
            }
            
            return binaryPath, nil
        }
    }
    
    return "", fmt.Errorf("binary %s not found in archive", binaryName)
}

func extractFromTarGz(tarGzPath, destDir, binaryName string) (string, error) {
    file, err := os.Open(tarGzPath)
    if err != nil {
        return "", err
    }
    defer file.Close()
    
    gzr, err := gzip.NewReader(file)
    if err != nil {
        return "", err
    }
    defer gzr.Close()
    
    tr := tar.NewReader(gzr)
    
    for {
        header, err := tr.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return "", err
        }
        
        // Look for the binary (might be in a subdirectory)
        if header.Typeflag == tar.TypeReg && strings.HasSuffix(header.Name, binaryName) {
            binaryPath := filepath.Join(destDir, binaryName)
            out, err := os.Create(binaryPath)
            if err != nil {
                return "", err
            }
            defer out.Close()
            
            if _, err := io.Copy(out, tr); err != nil {
                return "", err
            }
            
            return binaryPath, nil
        }
    }
    
    return "", fmt.Errorf("binary %s not found in archive", binaryName)
}

func getLatestVersion() (string, error) {
    resp, err := http.Get("https://api.github.com/repos/Murchoid/iwashere/releases/latest")
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return "", fmt.Errorf("GitHub API returned %s", resp.Status)
    }
    
    var release struct {
        TagName string `json:"tag_name"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return "", err
    }
    
    return release.TagName, nil
}

func downloadFile(url, path string) error {
    fmt.Printf("Downloading from: %s\n", url)
    
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        // Read error page
        body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
        return fmt.Errorf("download failed (HTTP %d): %s", resp.StatusCode, string(body))
    }
    
    
    out, err := os.Create(path)
    if err != nil {
        return err
    }
    defer out.Close()
    
    _, err = io.Copy(out, resp.Body)
    return err
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