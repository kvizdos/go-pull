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

## Production Readiness

As of now, this isn't really "prod" ready.. it technically could work, but its not super secure. Some ideas to make it prod ready are:
- Automatically run the built binary as a different user
- Sign binaries and put it in the Release description
    - Each time a new release is found, confirm the signature matches the binary present. If it doesn't, reject it.
    - This will at least confirm its *YOU* that is deploying stuff to your server. As usual, **secure your d@mn signing key!**
- In production mode, only pull Releases with the "Latest" tag- no pre-releases.
- Downtime concerns: currently, the process needs to be fully killed before proceeding. It'd be nice to somehow get both the old and new version running, and then swap the ports somehow (idk how this is possible)
