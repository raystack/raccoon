# Release Process

For maintainers, please read the sections below as a guide to create a new release.

## Create A New Release

Please follow these steps to create a new release:

* create a new tag of the form `vM.m.p`, where:
  * `M` = Major version, indicates there are breaking changes from the last Major version.
  * `m` = Minor version, indicates there are backward-compatible changes.
  * `p` = Patch version, indicates there are backward-compatible bug-fixes.

For example:
``` bash
$ git tag v1.2.0
```

* push the tags to trigger a release.
```bash
$ git push --tags
```

 Raccoon uses Goreleaser under the hood for release management. Each release pushes:
* A [github release](https://github.com/raystack/raccoon/releases/)
* A docker image to [raystack/raccoon](https://hub.docker.com/r/raystack/raccoon)
* Updates raystack's [homebrew-tap](https://github.com/raystack/homebrew-tap)
* Updates raystack's [scoop-bucket](https://github.com/raystack/scoop-bucket)

Additionally, the Github release will also contain with pre-built binaries for:
* `linux`
* `darwin` (macOS)
* `windows`

## Important Notes

* Raccoon release tags follow [SEMVER](https://semver.org/) convention.
* Github workflow is used to build and push the built docker image to Docker hub.
* A release is triggered when a github tag of format `vM.m.p` is pushed.
* Release tags should only point to main branch

