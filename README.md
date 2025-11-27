# paramsyouseek  
### **Advanced Parameter Discovery & Recon Framework**  
**Author:** Jai Rathee  
**Version:** v1.0

`paramsyouseek` is a powerful, multi-feature parameter discovery and crawling tool for bug bounty hunters, pentesters, and security researchers.

It extracts parameters from traditional HTML + URLs + JS + sitemaps and even generates new parameters using a guessing engine.  
It also supports JS-rendered pages via **Chromedp**, making it highly effective against modern web apps (React, Angular, Vue, etc.).

---

# 🧩 Features

### ✅ URL Parameter Extraction  
Extracts all query parameters from incoming URLs and crawled URLs.

### ✅ HTML Form Extraction  
Finds:
- `<input>` names  
- `<select>` names  
- `<textarea>` names  
- Hidden fields  

### ✅ JavaScript Extraction  
Finds:
- JS variable names  
- Parameter-like strings  
- Words in scripts  

### ✅ Headless JS Rendering (Chromedp)  
Renders dynamic JS-heavy pages and extracts parameters from:
- Ajax responses  
- React/Angular router URLs  
- Lazy-loaded components  
- API calls in network logs (future upgrade)  

### ✅ High-Speed Multi-threaded Crawler  
- No deadlocks  
- Atomic job counter  
- Duplicate avoidance  
- Depth control  
- Domain restriction  
- URL normalization  
- Asset skipping (.css, .js, images, etc.)

### ✅ Smart Parameter Guessing Engine  
Generates new parameter names using:
- JS variables  
- Page words  
- Embedded seed wordlist  
- Normalized tokens  

### ✅ Guess Testing (Lightweight & Safe)  
Detects:
- Reflected parameters  
- Server behavior changes  

### ✅ robots.txt & Sitemap Crawler  
Parses:
- `robots.txt`  
- `<loc>` URLs  
- `<sitemap>` references  

---

# 🏗️ Installation

### Build from source

```bash
git clone https://github.com/jairathee/paramsyouseek
cd paramsyouseek
go mod tidy
go build -o paramsyouseek
