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
