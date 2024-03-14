# Modify GitHub Issues

This is a GPTScript tool that uses the GitHub API to comment on or close issues.

Set your GitHub access token to the `GPTSCRIPT_GITHUB_TOKEN` environment variable in order to use this tool.
If the variable is not set, the tool will attempt to make unauthenticated API requests.

## Example

```yaml
tools: github.com/gptscript-ai/query-github-issues, github.com/gptscript-ai/modify-github-issues

Find the closed issue called "Test issue" in the repo g-linville/test and leave a comment on it.
```

## License

This tool is available under the Apache License 2.0. See [LICENSE](LICENSE) for more information.
