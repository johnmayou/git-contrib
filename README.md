# git-contrib

CLI tool that scans for local git repositories and prints a GitHub-style
contribution heatmap based on daily commit volume.

## Usage

### Scan for git repositories

```bash
git-contrib scan run ~/Projects
```

### See found git repositories

```bash
git-contrib scan show
```

Example output:

```text
~/Projects/vault
~/Projects/consul
~/Projects/vagrant
~/Projects/terraform
```

### Generate Heatmap

```bash
git-contrib stats --email me@email.com
```

Example output:

![heatmap](heatmap.png)
