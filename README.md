# mmstats

Retrieve user stats from a Mattermost instance.

## Develop

```bash
# Start the preview container
docker run --name mattermost-preview -d \
    --publish 8065:8065 --add-host dockerhost:127.0.0.1 \
    mattermost/mattermost-preview

# Generate sample data
docker exec -it mattermost-preview \
    mattermost sampledata --seed 10 --teams 4 --users 30
```

The following table shows sample user credentials.

| Username | Password         |
| -------- | ---------------- |
| sysadmin | Sys@dmin-sample1 |
| user-1   | SampleUs@r-1     |
| user-2   | SampleUs@r-2     |
| ......   | ............     |

## Reference

- [Local Machine Setup using Docker][1]
- [`mattermost sampledata`][2]

[1]: https://docs.mattermost.com/install/docker-local-machine.html
[2]: https://docs.mattermost.com/administration/command-line-tools.html#mattermost-sampledata

