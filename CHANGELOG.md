##  (2022-09-06)


### Features

*  support multiple event acknowledgement levels  ([#39](https://github.com/AkbaraliShaikh/raccoon/issues/39)) ([8bcf4f7](https://github.com/AkbaraliShaikh/raccoon/commit/8bcf4f7a919cd698d9aed7af3a29b5177e500764))
* add collector based design ([#17](https://github.com/AkbaraliShaikh/raccoon/issues/17)) ([193ba4d](https://github.com/AkbaraliShaikh/raccoon/commit/193ba4d68fd2ee41fe05acde11ee6fdc155fdaee))
* message ack sync changes improvements ([#43](https://github.com/AkbaraliShaikh/raccoon/issues/43)) ([ef219b6](https://github.com/AkbaraliShaikh/raccoon/commit/ef219b6ecaefa88852cd0bd842143063c2196b6f))
* parameterize user id header ([4453a6c](https://github.com/AkbaraliShaikh/raccoon/commit/4453a6c314378bc38811abcc07f835fe1c1ff94d))
* pull proto from proton ([b212119](https://github.com/AkbaraliShaikh/raccoon/commit/b212119469463ed2c35cce54d1ebca3ec322237d))
* support connection type ([#15](https://github.com/AkbaraliShaikh/raccoon/issues/15)) ([a562cd5](https://github.com/AkbaraliShaikh/raccoon/commit/a562cd5f2b9726e2a241da17f19d9ef7e0211f34))


### Bug Fixes

* application.yaml dependency during unit test ([24a73b2](https://github.com/AkbaraliShaikh/raccoon/commit/24a73b26ee6ba0aed451df4b9e1236ba2e40df87))
* checkorigin toggle ([b381101](https://github.com/AkbaraliShaikh/raccoon/commit/b381101a868595bb3adedf343383e0634c10b622))
* copy config on dockerfile ([99be932](https://github.com/AkbaraliShaikh/raccoon/commit/99be93252a5065fb69b585f2e83015cb42313db0))
* cyclic dependency when initializing log ([f8e9f9e](https://github.com/AkbaraliShaikh/raccoon/commit/f8e9f9e9577d5ea46ea40bfc57064aeadd6291bc))
* flaky handler test ([82a2f37](https://github.com/AkbaraliShaikh/raccoon/commit/82a2f3711887a885bb6baeec0a0f825919e81b77))
* handler test doesn't clean properly ([5855bd2](https://github.com/AkbaraliShaikh/raccoon/commit/5855bd28ba7914866797a9c19d5b58c2b6428674))
* host tag got overridden ([#26](https://github.com/AkbaraliShaikh/raccoon/issues/26)) ([b361657](https://github.com/AkbaraliShaikh/raccoon/commit/b3616571ad45623b1bb503722b734574e798b305))
* max connection error code and added missing reason ([#41](https://github.com/AkbaraliShaikh/raccoon/issues/41)) ([b895a6c](https://github.com/AkbaraliShaikh/raccoon/commit/b895a6c3104dba313614ceb3ba0780b7c6109f11))
* parse level on logger ([51b0a5e](https://github.com/AkbaraliShaikh/raccoon/commit/51b0a5e88bee045cab5b5f58e7d5b7d6d221566d))


### Performance Improvements

* **channel:** performance improvements  ([#36](https://github.com/AkbaraliShaikh/raccoon/issues/36)) ([2e34f07](https://github.com/AkbaraliShaikh/raccoon/commit/2e34f07bfac84a77f07d0c19b44162477b8c055f))
* **websockets:** performance optimisations for websockets ([#25](https://github.com/AkbaraliShaikh/raccoon/issues/25)) ([56977e9](https://github.com/AkbaraliShaikh/raccoon/commit/56977e93e0795c72dfa626a4feeacec4e806389b))


### Reverts

* Revert "[Rasyid] Fix Check origin failing when configuration is false" ([59a7e92](https://github.com/AkbaraliShaikh/raccoon/commit/59a7e921295ae39a477893e1be27d765ccd546bb))

