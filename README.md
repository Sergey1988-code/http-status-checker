# 🔍 HTTP Status Checker

**Lightning-fast command line tool for bulk HTTP status checking of your endpoints. Built in Go for maximum performance.**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Downloads](https://img.shields.io/github/downloads/Sergey1988-code/http-status-checker/total.svg)](https://github.com/Sergey1988-code/http-status-checker/releases)

**The Problem:** Need to check dozens of URLs? Manual browser checking takes hours, and `curl` loops are painful to write and maintain.

**The Solution:** `http-status-checker` - your reliable assistant. Single binary, instant results.

---

## 🚀 Features

- ✅ **Instant checking** of hundreds URLs simultaneously (thanks to goroutines!)
- 📊 **Beautiful color output** with intuitive visualization
- 📁 **File support** - simply list URLs in a text file
- ⚡ **Incredible speed** - checks dozens of URLs in seconds
- 🛡 **No dependencies** - single binary for Linux, Windows, macOS

---

## 📦 Installation

### Method 1: Download pre-built binary (recommended)

1. Go to [Releases](https://github.com/Sergey1988-code/http-status-checker/releases) page
2. Download version for your OS
3. Extract and add to PATH (or use directly from folder)

### Method 2: Build from source

```bash
git clone https://github.com/Sergey1988-code/http-status-checker.git
cd http-status-checker
go build -o http-status-checker main.go
```

---

## 🎯 Quick Start
1. **Create a file with your URLs:**

```bash
echo "https://httpstat.us/200
https://httpstat.us/404
https://httpstat.us/500
https://google.com" > urls.txt
```

2. **Run the check:**
```bash
./http-status-checker -file urls.txt
```

3. **Get instant results:**

```
🚀 Starting check of 4 URLs...

✅ 200 OK        https://httpstat.us/200
❌ 404 Not Found https://httpstat.us/404  
💥 500 Internal Server Error https://httpstat.us/500
✅ 200 OK        https://google.com

📊 Statistics:
• Total checked: 4
• Successful (2xx): 2
• Client errors (4xx): 1  
• Server errors (5xx): 1
• Execution time: 0.8s
```

---

##  Usage Examples

### 💡Basic usage
```bash
./http-status-checker -file my_urls.txt
```

### Show only problematic URLs
```bash
./http-status-checker -file urls.txt -only-errors
```

### Set timeout (in seconds)
```bash
./http-status-checker -file urls.txt -timeout 10
```

### Check single URL (without file)
```bash
./http-status-checker -url "https://example.com"
```

---

## 🛠 For Developers

### Usage in CI/CD pipelines
```yaml
# Example for GitHub Actions
- name: Check endpoints
  run: |
    ./http-status-checker -file production-urls.txt
    if [ $? -ne 0 ]; then
      echo "❌ Endpoint issues detected!"
      exit 1
    fi
```

### Usage in scripts
```bash
#!/bin/bash
echo "Checking service health..."

./http-status-checker -file services.txt --silent > results.json

# Further results processing...
```

---

## 🤝 Contributing

We welcome contributions to the project!

1. Fork the repository
2. Create your feature branch (git checkout -b feature/amazing-feature)
3. Commit your changes (git commit -m 'Add some amazing feature')
4. Push to the branch (git push origin feature/amazing-feature)
5. Open a Pull Request

---

## 📝 Roadmap

- [ ] Export results to JSON/CSV
- [ ] Check for keywords in response body
- [ ] Support for various HTTP methods (POST, PUT, DELETE)
- [ ] Configuration file for complex scenarios
- [ ] Real-time monitoring

---

## ⚠️ Troubleshooting
**Issue:** "Permission denied" when running
**Solution:** chmod +x http-status-checker

**Issue:** URLs checking very slowly
**Solution:** Make sure you're using the latest version. Older versions don't use concurrency.

---

## 📄 License
This project is distributed under MIT License. See  [LICENSE](LICENSE) file for details.

---

## 💬 Feedback
Found a bug or have a suggestion? Create an [Issue](https://github.com/Sergey1988-code/http-status-checker/issues)!

**Star the repository ⭐ if this tool saved you time!**

---

*Made with ❤️ in Go*