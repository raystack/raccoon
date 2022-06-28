## [v0.1.4](https://github.com/odpf/raccoon/compare/v0.1.3...v) (2022-06-28)


### Performance Improvements

* **channel:** performance improvements  ([#36](https://github.com/odpf/raccoon/issues/36)) ([2e34f07](https://github.com/odpf/raccoon/commit/2e34f07bfac84a77f07d0c19b44162477b8c055f))


## [v0.1.3](https://github.com/odpf/raccoon/compare/v0.1.2...v0.1.3) (2022-05-18)

### Chore
- **metrics:** add event level metrics ([#27](https://github.com/odpf/raccoon/issues/27))
- **proto:** buf lint fix ([#29](https://github.com/odpf/raccoon/issues/29))

### Docs
- adding protocol agnostic documentation ([#20](https://github.com/odpf/raccoon/issues/20))

### Fix
- host tag got overridden ([#26](https://github.com/odpf/raccoon/issues/26))

### Perf
- **websockets:** performance optimisations for websockets ([#25](https://github.com/odpf/raccoon/issues/25))

### Refactor
- update package names ([#34](https://github.com/odpf/raccoon/issues/34))
- introduces bootstrapper to orchestrate server initialization ([#23](https://github.com/odpf/raccoon/issues/23))





## [v0.1.2](https://github.com/odpf/raccoon/compare/v0.1.1...v0.1.2) (2021-11-23)


### Features

* Refactor codebase to Collector based design ([#19 ](https://github.com/odpf/raccoon/issues/19)) ([193ba4d](https://github.com/odpf/raccoon/commit/193ba4d68fd2ee41fe05acde11ee6fdc155fdaee ))
* Added support for JSON in websocket ([#19 ](https://github.com/odpf/raccoon/issues/19)) ([193ba4d](https://github.com/odpf/raccoon/commit/193ba4d68fd2ee41fe05acde11ee6fdc155fdaee ))
* Added support for REST endpoint with JSON and protobuf as supported data formats.  ([#19 ](https://github.com/odpf/raccoon/issues/19)) ([193ba4d](https://github.com/odpf/raccoon/commit/193ba4d68fd2ee41fe05acde11ee6fdc155fdaee )) 
* Added support for GRPC. ([#19 ](https://github.com/odpf/raccoon/issues/19)) ([193ba4d](https://github.com/odpf/raccoon/commit/193ba4d68fd2ee41fe05acde11ee6fdc155fdaee ))


### Bug Fixes

NA

## [v0.1.1](https://github.com/odpf/raccoon/compare/v0.1.0...v) (2021-10-28)


### Features

* support multitenancy ([#15](https://github.com/odpf/raccoon/issues/15)) ([a562cd5](https://github.com/odpf/raccoon/commit/a562cd5f2b9726e2a241da17f19d9ef7e0211f34))


### Bug Fixes

* checkorigin toggle rejects every connection when set to false ([b381101](https://github.com/odpf/raccoon/commit/b381101a868595bb3adedf343383e0634c10b622))
