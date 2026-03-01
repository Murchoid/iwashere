# Installing iwashere on Linux

## Quick Install

# Download
wget https://github.com/Murchoid/iwashere/releases/download/v0.1.0/iwashere-linux-amd64

# Make executable and move to PATH
chmod +x iwashere-linux-amd64
sudo mv iwashere-linux-amd64 /usr/local/bin/iwashere

# Test it
```
iwashere --help
```

One-Liner Install

```
curl -L https://github.com/Murchoid/iwashere/releases/download/v0.1.0/iwashere-linux-amd64 -o /tmp/iwashere && chmod +x /tmp/iwashere && sudo mv /tmp/iwashere /usr/local/bin/iwashere 
