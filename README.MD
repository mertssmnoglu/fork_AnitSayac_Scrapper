# Anıt Sayaç Crawler for Golang

**🌐 Language Options:** **English** | [Türkçe](README_tr.MD)

It is recommended that people who will use it for data processing use the JSON file. The JSON file is updated every every day at midnight (00:00) for UTC. The data is taken from the site [Anıt Sayaç](https://anitsayac.com/).

## ❗ Important Information

- It is known that errors occur due to incorrect data entry in data with IDs 70, 69, 114, 78, 71, 511, 605, 607, 3580, 1051, 934, 675, 617, 616, 611, 599, 138, 4185.
- Blank data appears that way because it is not available on the site.

## Installation

```bash
go mod tidy
go build -o crawler crawler.go
```

## Usage

```bash
go run .
```

## Git Installation and Usage for Data Research

### Installing Git

#### Windows

1. Download Git from [https://git-scm.com/download/win](https://git-scm.com/download/win)
2. Run the installer and follow the setup wizard
3. Verify installation: `git --version`

#### macOS

```bash
# Using Homebrew
brew install git

# Or download from https://git-scm.com/download/mac
```

#### Linux (Ubuntu/Debian)

```bash
sudo apt update
sudo apt install git
```

### Git Commands for Data Research

#### Search for specific names in commit history

```bash
# Search for "Fidan Çakır" in all commits across all branches
git log --oneline -S "Fidan Çakır" --source --all

# Search for "İkbal Uzuner" in commit history
git log --oneline -S "İkbal Uzuner" --source --all

# Search with grep in commit messages
git log --all --grep="Fidan Çakır" --oneline
```

#### View detailed changes for specific commits

```bash
# Show files changed in a specific commit
git show --name-only COMMIT_HASH

# Show full diff for a commit
git show COMMIT_HASH

# Show changes with patch details
git log -S "SEARCH_TERM" --patch --oneline
```

**Note:** Replace `COMMIT_HASH` with the actual commit hash (e.g., `38b5d7f` or `c3e1dfb`). You can find commit hashes using `git log --oneline` command.

#### Search in file contents

```bash
# Search current files for a term
git grep "Fidan Çakır"

# Search in specific file types
git grep "search_term" -- "*.json"
git grep "search_term" -- "*.csv"
```

#### Browse repository history

```bash
# View all commits in one line
git log --oneline

# View commits with file changes
git log --stat

# View commits in a specific time range
git log --since="2024-01-01" --until="2024-12-31"
```

#### Examples for this repository

```bash
# Find all commits that mention specific victims
git log --oneline -S "Fidan Çakır" --source --all
git log --oneline -S "İkbal Uzuner" --source --all

# Search for specific locations
git log --oneline -S "İzmir" --source --all
git log --oneline -S "İstanbul" --source --all

# Find commits by date patterns
git log --oneline -S "16/10/2024" --source --all

# Search in data files only
git log --oneline -- data.json data.csv
```

## Data Source

- [Anıt Sayaç](https://anitsayac.com/)

## License

[MIT](LISENCE)

## Author

[Ramazan Sancar](https://www.github.com/ramazansancar)
