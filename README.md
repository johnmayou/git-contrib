# git-contrib

CLI tool to visualize local commit contributions with a GitHub-style heatmap.

## Usage

### 1. Scan for Repos

```bash
git-contrib scan run ~/Projects
```

### 2. See Found Repos

```bash
git-contrib scan show
```

Example output:

```plaintext
~/Projects/vault
~/Projects/consul
~/Projects/vagrant
~/Projects/terraform
```

### 3. Generate Heatmap

```bash
git-contrib stats --email me@email.com
```

Example output:

![heatmap](heatmap.png)

## License

[MIT](./LICENSE)
