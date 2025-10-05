# git-contrib

CLI tool to visualize your local commit contributions with a GitHub-style heatmap.

## Usage

Scan for repos:

```bash
git-contrib scan run ~/Projects
```

See found repos:

```bash
git-contrib scan show
```

```plaintext
~/Projects/vault
~/Projects/consul
~/Projects/vagrant
~/Projects/terraform
```

Generate heatmap:

```bash
git-contrib stats --email me@email.com
```

![heatmap](heatmap.png)

## License

[MIT](./LICENSE)
