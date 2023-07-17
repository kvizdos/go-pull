# Go Pull

A super basic deployment system based on GitHub Releases (not meant for production).

Every hour, this will pull the latest GitHub Releases from your project. If there is a new release out, it will pull it and restart the running process to reflect the updates.

**THIS IS NOT MEANT FOR PRODUCTION.**

Setting it up is easy. First, create a `.env` file with the following items filled:

```
PIPELINE_GITHUB_TOKEN=<GITHUB PAT>
PIPELINE_RELEASES_API=<GITHUB PROJECT URL>/releases
PIPELINE_BUILD_OUT=./app
```

The system will automatically make `./app` group `0700`, giving the owner (whoever runs the Go script) access to READ / WRITE / EXECUTE the app, and restricts all other access from group users or regular users.
