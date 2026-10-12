# Info

## Repository

### Worktree Path

`{{ .RepoWorktreePath }}`

### Status

```
{{ .RepoStatus }}
```

### Diff

| Filename | Diff |
| --- | --- |
{{- range $f := .RepoDiffFiles }}
| {{ escapeMarkdown $f.Filename }} | {{ markdownTableCode $f.Diff }} |
{{- end }}

## Git Executable

- Version:  `{{ .GitExeVersion }}`
- Path: `{{ .GitExePath }}`