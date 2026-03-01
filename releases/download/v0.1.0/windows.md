# Installing iwashere on Windows

## Quick Install (2 minutes)

1. **Download** `iwashere-windows-amd64.exe` from the releases page
2. **Rename** it to `iwashere.exe`
3. **Create folder**: `C:\tools\iwashere\`
4. **Move** `iwashere.exe` into that folder
5. **Add to PATH**:
   - Press `Windows Key + X`
   - Select "System"
   - Click "Advanced system settings"
   - Click "Environment Variables"
   - Under "System variables", find "Path" and click "Edit"
   - Click "New" and add `C:\tools\iwashere`
   - Click "OK" on all windows
6. **Restart** your terminal
7. **Test** it:
   ```cmd
   iwashere --help

One-Line Install (if you have curl)
Open PowerShell as Administrator and run:


```md C:\tools\iwashere -Force
curl -L https://github.com/Murchoid/iwashere/releases/download/v0.1.0/iwashere-windows-amd64.exe -o C:\tools\iwashere\iwashere.exe
setx PATH "$env:PATH;C:\tools\iwashere" /M```