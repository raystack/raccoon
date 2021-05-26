# Release
For maintainers, please read the sections below as a guide to create a new release.
## Create A New Release
Please follow these steps to create a new release:
- Update `version.txt` file
- Generate changelog from commits by using [conventional-changelog-cli](https://www.npmjs.com/package/conventional-changelog-cli#quick-start)
  ```sh
  $ conventional-changelog -s -p conventionalcommits -i CHANGELOG.md
  ```
- Commit `version.txt` and `CHANGELOG.md` together and mark the commit with the release tag. Make sure the release tag and `version.txt` are the same.
  ```sh
  $ git add version.txt CHANGELOG.md
  $ git commit -m "docs: update changelog and version for vM.m.p"
  $ git tag vM.m.p
  ```
- Push the commit and the tag. Release action will trigger to publish docker image and create GitHub release.

## Important Notes
- Raccoon release tags follow [SEMVER](https://semver.org/) convention.
- Github workflow is used to build and push the built docker image to Docker hub.
- A release is triggered when a github tag of format `vM.m.p` is pushed. After the release job is succeeded, a docker image of
format `M.m.p` is pushed to docker hub.
- Release tags should only point to main branch