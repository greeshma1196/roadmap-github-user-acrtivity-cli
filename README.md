# roadmap-github-user-activity-cli

A simple [CLI](https://roadmap.sh/projects/github-user-activity) to fetch the recent activity of a GitHub user and display it in the terminal.

Events fetched:

1. CreateEvent
2. DeleteEvent
3. IssuesEvent
4. PullRequestEvent
5. PushEvent
6. ReleaseEvent

## How To Use

### build cli

```sh
make build
```

### run cli

```sh
./github-activity USER_NAME
```

### run tests

```sh
make test
```

### run lint

```sh
make lint
```
