# 📜 gobash Example Scripts

Ready-to-use scripts that enhance your gobash experience!

## 🚀 Quick Start

```bash
# Make scripts executable
chmod +x *.sh

# Use inside gobash
./gobash
[IPTP-1] ~$ ./workspace.sh
[IPTP-1] ~$ ./qf.sh myproject
```

## 📚 Available Scripts

### 1. **workspace.sh** - Quick Workspace Selector
Navigate to your common workspaces quickly.

```bash
[IPTP-1] ~$ ./workspace.sh
🗂️  Quick Workspace Selector
============================

Current: /Users/you

1) projects         - Main projects
2) personal         - Personal projects
3) work             - Work projects
4) Documents        - Documents
5) Downloads        - Downloads

Select (1-5, or path): 1

✓ Selected: Main projects
  Current: /Users/you/projects
💡 Tip: Use name "Main projects" to name this session
```

**Customize it:** Edit the `WORKSPACES` array at the top of the script:
```bash
WORKSPACES=(
    "$HOME/mydir:My Description"
    "$HOME/other:Other Description"
)
```

### 2. **newproject.sh** - Project Initializer
Create and set up new projects with templates.

```bash
[IPTP-1] ~$ ./newproject.sh myapp go
🚀 Creating project: myapp (go)
==================================================

→ Setting up Go project...
✓ Go module initialized
✓ Git repository initialized

✅ Project created successfully!

💡 You're now in: /Users/you/myapp
   Try: name "working on myapp"

[IPTP-1] myapp$ name "working on myapp"
✓ Shell named: myapp
[myapp] myapp$ 
```

**Supported types:**
- `go` - Go project with module
- `node` - Node.js project with package.json
- `python` - Python project with venv setup
- `rust` - Rust project with Cargo
- `simple` - Basic directory structure

### 3. **gitstatus.sh** - Multi-Repo Status Checker
Check git status across multiple repositories.

```bash
[IPTP-1] projects$ ./gitstatus.sh
📊 Git Repository Status
========================

Subdirectories:
  ✓ auth-service              [main]           clean
  ⚠ web-frontend              [dev]            modified (3 files), 2 unpushed
  ✓ data-pipeline             [main]           clean
  ⚠ mobile-app                [feature/login]  modified (1 files), 1 untracked

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Summary: 4 repositories, 2 with changes

💡 Tip: Use 'cd <repo>' to navigate to a repository
```

### 4. **qf.sh** - Quick Find
Fast directory finder with fuzzy matching.

```bash
[IPTP-1] ~$ ./qf.sh golang
🔍 Searching for: golang

Searching locally...

Found 3 matching directories:

 1) ./projects/golang-api
 2) ./projects/golang-tools
 3) ./learning/golang-book

Select (1-3, or 'q' to quit): 1

✓ Selected: ./projects/golang-api
  Current: /Users/you/projects/golang-api
💡 Tip: Use save to bookmark this location
```

**Features:**
- Searches current directory first (fast!)
- Falls back to home directory if needed
- Auto-navigates if only one match
- Works inside and outside gobash

## 🎯 Usage Patterns

### Pattern 1: Quick Navigation
```bash
[IPTP-1] ~$ ./qf.sh myproject
# Instantly jump to myproject directory

[myapp] myproject$ name "working on myproject"
# Name your session

[myapp] myproject$ save
# Save the location
```

### Pattern 2: New Project Setup
```bash
[IPTP-1] ~$ ./newproject.sh authapi go
# Creates and navigates to new project

[IPTP-1] authapi$ name "auth api development"
# Set meaningful name

[auth_api_development] authapi$ save
# Save for later
```

### Pattern 3: Multi-Project Management
```bash
[IPTP-1] ~/projects$ ./gitstatus.sh
# Check all repos

[IPTP-1] projects$ cd web-frontend
# Jump to repo with changes

[IPTP-1] web-frontend$ git status
# Check details

[IPTP-1] web-frontend$ git add .
[IPTP-1] web-frontend$ git commit -m "Update"
```

### Pattern 4: Workspace Switching
```bash
[auth] myproject$ ./workspace.sh
# Select different workspace

[auth] projects$ cd other-project
# Continue working
```

## 🛠️ Customization

### Adding Your Own Workspaces

Edit `workspace.sh`:
```bash
WORKSPACES=(
    "$HOME/dev/backend:Backend services"
    "$HOME/dev/frontend:Frontend apps"
    "$HOME/research:Research projects"
    "/mnt/data:Data analysis"
)
```

### Creating New Project Types

Edit `newproject.sh`, add to the case statement:
```bash
mytype)
    echo "→ Setting up my custom project..."
    mkdir -p src tests docs
    # Your setup commands here
    echo "✓ Custom project initialized"
    ;;
```

### Adjusting Search Depth

Edit `qf.sh`:
```bash
# Local search depth
find . -maxdepth 5 ...  # Change from 3 to 5

# Home search depth
find ~ -maxdepth 7 ...  # Change from 5 to 7
```

## 💡 Tips

1. **Create a scripts directory:**
   ```bash
   mkdir ~/gobash-scripts
   cp *.sh ~/gobash-scripts/
   # Add to PATH or create aliases
   ```

2. **Make scripts global:**
   ```bash
   sudo cp *.sh /usr/local/bin/
   # Now use: workspace.sh anywhere
   ```

3. **Create shell aliases:**
   ```bash
   # Add to ~/.zshrc or ~/.bashrc
   alias ws='~/gobash-scripts/workspace.sh'
   alias qf='~/gobash-scripts/qf.sh'
   alias newp='~/gobash-scripts/newproject.sh'
   alias gits='~/gobash-scripts/gitstatus.sh'
   ```

4. **Combine scripts:**
   ```bash
   [IPTP-1] ~$ ./qf.sh myproject && ./gitstatus.sh
   # Find and check git status in one go
   ```

## 📖 See Also

- [SCRIPTING_GUIDE.md](../SCRIPTING_GUIDE.md) - Full guide to writing gobash scripts
- [README.md](../../README.md) - Main gobash documentation

## 🤝 Contributing

Got a useful script? Share it! These examples are meant to inspire and be customized for your workflow.

---

**Remember:** These scripts work both inside and outside gobash, but they're optimized for use within gobash for automatic navigation and state management! 🚀
