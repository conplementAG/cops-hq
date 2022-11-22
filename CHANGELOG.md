# Changelog

## [2.1.0](https://github.com/conplementAG/cops-hq/compare/v2.0.2...v2.1.0) (2022-11-22)


### Features

* **executor:** direct os/exec command support for complex cases ([95013a1](https://github.com/conplementAG/cops-hq/commit/95013a1923bcf5f27599710169dc78da074f697f))

## [2.0.2](https://github.com/conplementAG/cops-hq/compare/v2.0.1...v2.0.2) (2022-11-17)


### Bug Fixes

* **executor:** prevent escaping for quotations which are parts of the arguments ([80be3ed](https://github.com/conplementAG/cops-hq/commit/80be3eda20f6c622559f628d6d74beea6410545c))
* **terraform-recipe:** remove quotations which were anyways always removed before ([dc49955](https://github.com/conplementAG/cops-hq/commit/dc49955337a5ecd994dd69afa8d60006cfdb028b))

## [2.0.1](https://github.com/conplementAG/cops-hq/compare/v2.0.0...v2.0.1) (2022-10-25)


### Bug Fixes

* correct v2 package versioning ([e355a7c](https://github.com/conplementAG/cops-hq/commit/e355a7c4be56c85704182928c6218c839425f9d2))

## [2.0.0](https://github.com/conplementAG/cops-hq/compare/v1.0.0...v2.0.0) (2022-10-24)


### ⚠ BREAKING CHANGES

* terraform - auto approve parameter in deploy and destroy flows

### Features

* configurable RequireInfrastructureEncryption setting for terraform state storage ([35e7759](https://github.com/conplementAG/cops-hq/commit/35e7759ef9102c2154d37128abe62eda1f044c03))
* terraform - auto approve parameter in deploy and destroy flows ([d745cbe](https://github.com/conplementAG/cops-hq/commit/d745cbe0ad3d189be640d50dd90a1fabdc4a2e29))

## [1.0.0](https://github.com/conplementAG/cops-hq/compare/v0.12.0...v1.0.0) (2022-09-20)


### ⚠ BREAKING CHANGES

* **cli:** drop shorthand flag "s" from global cli flags

### Features

* **cli:** drop shorthand flag "s" from global cli flags ([6fbf977](https://github.com/conplementAG/cops-hq/commit/6fbf977281ff2c071638fb290b279974288ba4cd))

## [0.12.0](https://github.com/conplementAG/cops-hq/compare/v0.11.1...v0.12.0) (2022-09-12)


### Features

* **configuration:** make decrypted config available ([173af36](https://github.com/conplementAG/cops-hq/commit/173af363f80dcb2a3a8761807e189be91a87cd11))

## [0.11.1](https://github.com/conplementAG/cops-hq/compare/v0.11.0...v0.11.1) (2022-09-02)


### Bug Fixes

* **check-dependencies:** downgrade min. kubectl version to v1.23.9 ([4844daf](https://github.com/conplementAG/cops-hq/commit/4844daf556bee022535e6699e8150d453db883a5))
* undo delete .idea\misc.xml file ([15cdabb](https://github.com/conplementAG/cops-hq/commit/15cdabb9069db9896d2785b54e8e2ad3bfe44347))

## [0.11.0](https://github.com/conplementAG/cops-hq/compare/v0.10.0...v0.11.0) (2022-09-01)


### Features

* Dockerfile update with copsctl 0.8.2 ([b1a5596](https://github.com/conplementAG/cops-hq/commit/b1a55967c712ed22ad2a49e48f72c245a9d6b4b5))
* tooling upgrade ([fae3351](https://github.com/conplementAG/cops-hq/commit/fae3351a10d61db26d1563038df582280997ab18))


### Bug Fixes

* **terraform-recipe:** adapt IP firewall rule handling backend storage ([efeaef5](https://github.com/conplementAG/cops-hq/commit/efeaef5f7e3d4a7d5668612dd4091800db897256))

## [0.10.0](https://github.com/conplementAG/cops-hq/compare/v0.9.1...v0.10.0) (2022-08-19)


### Features

* **terraform-recipe:** add support to define IP firewall rules for backend storage ([0102fb8](https://github.com/conplementAG/cops-hq/commit/0102fb87914c4da8195fab7265b5492044d30963))

## [0.9.1](https://github.com/conplementAG/cops-hq/compare/v0.9.0...v0.9.1) (2022-08-08)


### Bug Fixes

* **terraform-recipe:** renaming optional tags ([7e11d6d](https://github.com/conplementAG/cops-hq/commit/7e11d6d20d05403a35fc4ff6ef291dd6ddd73820))

## [0.9.0](https://github.com/conplementAG/cops-hq/compare/v0.8.0...v0.9.0) (2022-07-18)


### Features

* **terraform-recipe:** resource group tagging ([16e1a0d](https://github.com/conplementAG/cops-hq/commit/16e1a0dc85cec683cfc02a1129eb0fdf68796db6))

## [0.8.0](https://github.com/conplementAG/cops-hq/compare/v0.7.0...v0.8.0) (2022-07-13)


### Features

* cluster info extensions for navigating config sub-objects ([a98e307](https://github.com/conplementAG/cops-hq/commit/a98e30755488f53f70d1a27a2ac0b08f3f347111))

## [0.7.0](https://github.com/conplementAG/cops-hq/compare/v0.6.0...v0.7.0) (2022-07-12)


### Features

* **cli:** default command and global init functionality ([ae9b4d9](https://github.com/conplementAG/cops-hq/commit/ae9b4d972bf5679118ebc6da93ee0ebed3776477))
* **hq:** new with custom options, possibility to disable logging to file ([dd0d31c](https://github.com/conplementAG/cops-hq/commit/dd0d31cd6990e0084473aa5152372344f761d12b))
* private endpoint resource type ([90987cd](https://github.com/conplementAG/cops-hq/commit/90987cd5b6ee62075a88cc96217aeba47d02caa2))

## [0.6.0](https://github.com/conplementAG/cops-hq/compare/v0.5.0...v0.6.0) (2022-06-30)


### Features

* helm recipe ([f837994](https://github.com/conplementAG/cops-hq/commit/f837994ab6cfd6f8e6117439e607a0e2f9e01dba))


### Bug Fixes

* adapt terraform/helm recipe docu with variable autopopulate example ([ecca07c](https://github.com/conplementAG/cops-hq/commit/ecca07c43ae3fe538f2482a6565e41e6baf699b4))
* add helm recipe docu ([d9b241b](https://github.com/conplementAG/cops-hq/commit/d9b241b2839c924e0c2a737333429bd6fe91f4d6))
* terraform recipe unit test ([0975779](https://github.com/conplementAG/cops-hq/commit/0975779a78053f1c173b3812c8212607c47e3e51))

## [0.5.0](https://github.com/conplementAG/cops-hq/compare/v0.4.0...v0.5.0) (2022-06-14)


### Features

* plan analyzer ([230c3bb](https://github.com/conplementAG/cops-hq/commit/230c3bbdc61da55d369cf41052e7b08d938bc5e9))


### Bug Fixes

* **copsctl:** query to prompt connect with less noise and no chance to fail ([6a74d5a](https://github.com/conplementAG/cops-hq/commit/6a74d5a92f2f168237b9e7d84901dd735c2f6824))

## [0.4.0](https://github.com/conplementAG/cops-hq/compare/v0.3.0...v0.4.0) (2022-06-13)


### Features

* copsctl recipe ([4864b1e](https://github.com/conplementAG/cops-hq/commit/4864b1ee765a8747b925f09f3c7b5b353a3bacc0))

## [0.3.0](https://github.com/conplementAG/cops-hq/compare/v0.2.0...v0.3.0) (2022-06-01)


### Features

* terraform files now saved to separate directory ([de5b09e](https://github.com/conplementAG/cops-hq/commit/de5b09e2cdd637ec47cf33bba5869a4f5ada35f4))

## [0.2.0](https://github.com/conplementAG/cops-hq/compare/v0.1.0...v0.2.0) (2022-06-01)


### Features

* add github action cibuild ([ccdefd2](https://github.com/conplementAG/cops-hq/commit/ccdefd2b9a50d39d4bd5f50dd50af0dcca021858))
* LoadConfigFile override, AskUserToConfirmWithKeyword ([fe7864f](https://github.com/conplementAG/cops-hq/commit/fe7864f7522348ab1a11a37722e41d16f7a8b335))


### Bug Fixes

* add fileending for workflow cibuild.yml ([a01469c](https://github.com/conplementAG/cops-hq/commit/a01469c1c57460d323bd724d21c3b5b502bcd272))
* Dockerfile & unit tests ([582783a](https://github.com/conplementAG/cops-hq/commit/582783a44ee27fbfa6024519cee5f4c549aa8f48))
* remove linebreaks from azure cli output ([b71ea4d](https://github.com/conplementAG/cops-hq/commit/b71ea4d95638d168a9a3875db57bcc5f65ff471b))
* remove typo ([0b04ced](https://github.com/conplementAG/cops-hq/commit/0b04ced91cebe2d58e3d5195b243f1550c6d231a))
* support missing module in naming service ([9decd03](https://github.com/conplementAG/cops-hq/commit/9decd03c29db36b78c7594418128ca681f296691))

## [0.1.0](https://github.com/conplementAG/cops-hq/compare/v0.0.3...v0.1.0) (2022-05-27)


### Features

* adding major and minor tags during release ([88fbc7e](https://github.com/conplementAG/cops-hq/commit/88fbc7e98acddaf53779162b7a9467b5b3eea1e1))
* trigger releaser on every commit ([0cb84bc](https://github.com/conplementAG/cops-hq/commit/0cb84bce17ed06dac3d39dd0d9144426d213eac6))
