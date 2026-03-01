## 3. Simple Linux Install Instructions


# Installing iwashere on Linux

## Quick Install

```bash
# Download
wget https://github.com/Murchoid/iwashere/releases/download/v0.1.0/iwashere-linux-amd64

# Make executable and move to PATH
chmod +x iwashere-linux-amd64
sudo mv iwashere-linux-amd64 /usr/local/bin/iwashere

# Test it
iwashere --help
One-Liner Install
bash
curl -L https://github.com/Murchoid/iwashere/releases/download/v0.1.0/iwashere-linux-amd64 -o /tmp/iwashere && chmod +x /tmp/iwashere && sudo mv /tmp/iwashere /usr/local/bin/iwashere
text

## 4. GitHub Release Template

When creating your release on GitHub, use this template:

```markdown
# iwashere v0.1.0 - Never lose your context again! 🎯

A CLI tool that remembers where you left off in your coding projects.

##  Features
- **`init`** - Start tracking a project
- **`add`** - Save notes with git context
- **`show`/`list`** - View your notes
- **`edit`/`delete`** - Manage notes
- **`branch`** - Branch-specific notes
- **`session`** - Track work sessions
- **`tag`** - Organize notes

## Quick Start
```bash
cd your-project
iwashere init
iwashere add "Starting work on feature X"
iwashere list