# :recycle: Contributing

Contributions to this project are very much welcome, please make sure that your code changes are tested and that they follow
Go best-practices.

## Getting started

### Required environmental variables

The tests require at least two environment variables: `STREAM_CHAT_API_KEY` and `STREAM_CHAT_API_SECRET`. There are multiple ways to provide that:
- simply set it in your current shell (`export STREAM_CHAT_API_KEY=xyz`)
- you could use [direnv](https://direnv.net/)
- if you debug the tests in VS Code, you can set up an env file there as well: `"go.testEnvFile": "${workspaceFolder}/.env"`.

### Code formatting & linter

We enforce code formatting with [`gufumpt`](https://github.com/mvdan/gofumpt) (a stricter `gofmt`). If you use VS Code, it's recommended to set this setting there for auto-formatting:

```json
"editor.formatOnSave": true,
"gopls": {
    "formatting.gofumpt": true
}
```

Gofumpt will mostly take care of your linting issues as well.
