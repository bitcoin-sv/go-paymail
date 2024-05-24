# Changelog

## [0.15.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.14.0...v0.15.0) (2024-05-24)


### Features

* **SPV-789:** extend PIKE capability ([#91](https://github.com/bitcoin-sv/go-paymail/issues/91)) ([6b1a07c](https://github.com/bitcoin-sv/go-paymail/commit/6b1a07c7cd68b492d7fa7e590a9faeb7b269ba2e))

## [0.14.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.13.0...v0.14.0) (2024-03-26)


### Features

* **BUX-395:** add pike capabilitiy ([#77](https://github.com/bitcoin-sv/go-paymail/issues/77)) ([3c6789d](https://github.com/bitcoin-sv/go-paymail/commit/3c6789d1352cba6297076b7c11650b5c06334bc6))

## [0.13.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.12.1...v0.13.0) (2024-02-26)


### Features

* **223:** replace router with gin ([#74](https://github.com/bitcoin-sv/go-paymail/issues/74)) ([9d9f3e7](https://github.com/bitcoin-sv/go-paymail/commit/9d9f3e7dc2706c0fceba7399f6afaadc504cf947))

## [0.12.1](https://github.com/bitcoin-sv/go-paymail/compare/v0.12.0...v0.12.1) (2024-01-22)


### Bug Fixes

* **BUX-497:** Routes are hardcoded instead of initialized by configured capabilities ([#71](https://github.com/bitcoin-sv/go-paymail/issues/71)) ([8e14dd0](https://github.com/bitcoin-sv/go-paymail/commit/8e14dd09fe732b3b27de1bf7303e2cee777ffac2))

## [0.12.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.11.0...v0.12.0) (2024-01-05)


### Features

* **BUX-437:** go-resty version update to avoid security alert ([#69](https://github.com/bitcoin-sv/go-paymail/issues/69)) ([c9d0558](https://github.com/bitcoin-sv/go-paymail/commit/c9d0558040f8853609a33d98f65a1f351dad085d))

## [0.11.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.10.0...v0.11.0) (2023-12-21)


### Features

* **BUX-420:** go version and workflows update ([#64](https://github.com/bitcoin-sv/go-paymail/issues/64)) ([64a7bd4](https://github.com/bitcoin-sv/go-paymail/commit/64a7bd4122342794ad57535a52583f5acdc47670))

## [0.10.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.9.4...v0.10.0) (2023-12-21)


### Features

* **BUX-206:** unify logs ([#63](https://github.com/bitcoin-sv/go-paymail/issues/63)) ([849f4db](https://github.com/bitcoin-sv/go-paymail/commit/849f4dbea0de1d66ee89c7cb677f87536e89404a))

## [0.9.4](https://github.com/bitcoin-sv/go-paymail/compare/v0.9.3...v0.9.4) (2023-12-05)


### Bug Fixes

* Update nLocktime nSequence handling ([#57](https://github.com/bitcoin-sv/go-paymail/issues/57)) ([92eff98](https://github.com/bitcoin-sv/go-paymail/commit/92eff9847d23c805e910588b90abe67baa7b1c02))

## [0.9.1](https://github.com/bitcoin-sv/go-paymail/compare/v0.9.0...v0.9.1) (2023-11-24)


### Bug Fixes

* fix v0.9.1 ([#54](https://github.com/bitcoin-sv/go-paymail/issues/54)) ([d996d9c](https://github.com/bitcoin-sv/go-paymail/commit/d996d9c4424aee32eb001526dfcd7ab0cfade8d4))

## [0.9.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.8.0...v0.9.0) (2023-11-24)


### Features

* **BUX-172:** verify merkle root for unmined inputs ([#48](https://github.com/bitcoin-sv/go-paymail/issues/48)) ([f719ffa](https://github.com/bitcoin-sv/go-paymail/commit/f719ffa5fffddaa327c0ff0f79cd5c17845eb2f3))

## [0.8.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.7.2...v0.8.0) (2023-11-24)


### Features

* **BUX-322:** add decoded beef to P2PTransaction ([#51](https://github.com/bitcoin-sv/go-paymail/issues/51)) ([94efda0](https://github.com/bitcoin-sv/go-paymail/commit/94efda042d8b29e6aa589f1a12f549eade023f30))

## [0.7.2](https://github.com/bitcoin-sv/go-paymail/compare/v0.7.1...v0.7.2) (2023-11-15)


### Bug Fixes

* **BUX-250:** fix decode;add spv and decode tests for corrupted/invalid beef ([#45](https://github.com/bitcoin-sv/go-paymail/issues/45)) ([6fb0fc9](https://github.com/bitcoin-sv/go-paymail/commit/6fb0fc9ee537b519a4145a3c3f781830e136cfe4))

## [0.7.1](https://github.com/bitcoin-sv/go-paymail/compare/v0.7.0...v0.7.1) (2023-11-14)


### Bug Fixes

* **BUX-250:** SPV; decoding ([#43](https://github.com/bitcoin-sv/go-paymail/issues/43)) ([5e253d6](https://github.com/bitcoin-sv/go-paymail/commit/5e253d6ef1c259a45752d85475bd5f3db633e8c4))

## [0.7.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.6.0...v0.7.0) (2023-11-06)


### Features

* **BUX-250:** logging and returning more descriptive errors from tx processing ([#40](https://github.com/bitcoin-sv/go-paymail/issues/40)) ([10d1da7](https://github.com/bitcoin-sv/go-paymail/commit/10d1da75f1c210d0c55d9f5138509e60911fb9ac))

## [0.6.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.5.1...v0.6.0) (2023-11-03)


### Features

* **BUX-290:** updated calculation and verification of merkle roots ([#37](https://github.com/bitcoin-sv/go-paymail/issues/37)) ([c658f96](https://github.com/bitcoin-sv/go-paymail/commit/c658f964d0d5afd14b49fe26458b9674ed776a96))
* **BUX-296:** adjust BEEF tx decoding to BUMP structure ([#34](https://github.com/bitcoin-sv/go-paymail/issues/34)) ([2ae3790](https://github.com/bitcoin-sv/go-paymail/commit/2ae3790e077891111538555624cba1fd7c877e2f))

## [0.5.1](https://github.com/bitcoin-sv/go-paymail/compare/v0.5.0...v0.5.1) (2023-10-13)


### Bug Fixes

* **BUX-272:** Change the order of elements in decoded CMP ([d391319](https://github.com/bitcoin-sv/go-paymail/commit/d3913191a30b3d3c44009d73730cad2d8dd260cf))

## [0.5.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.4.0...v0.5.0) (2023-10-12)


### Features

* **BUX-164:** verify beef tx ([#22](https://github.com/bitcoin-sv/go-paymail/issues/22)) ([90b9efb](https://github.com/bitcoin-sv/go-paymail/commit/90b9efb72caa70df217c078c1d282e7fa53fb1c3))

## [0.4.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.3.0...v0.4.0) (2023-10-04)


### Features

* adds new version of codeowner/codestandards/template files ([#24](https://github.com/bitcoin-sv/go-paymail/issues/24)) ([74f16e0](https://github.com/bitcoin-sv/go-paymail/commit/74f16e0d9c9f700a77181d32b2b925baf0d9d6b6))

## [0.3.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.2.1...v0.3.0) (2023-09-29)


### Features

* delete ExecuteSimplifiedPaymentVerification from  PaymailServiceProvider ([#20](https://github.com/bitcoin-sv/go-paymail/issues/20)) ([14e4b69](https://github.com/bitcoin-sv/go-paymail/commit/14e4b6901537d5e807fed37c8f84f54bebe9d873))

## [0.2.1](https://github.com/bitcoin-sv/go-paymail/compare/v0.2.0...v0.2.1) (2023-09-28)


### Bug Fixes

* **BUX-242:** fix go-paymail routes ([#17](https://github.com/bitcoin-sv/go-paymail/issues/17)) ([8bb7077](https://github.com/bitcoin-sv/go-paymail/commit/8bb7077ff7092acc8f3eebcb24e78a0dac10097b))

## [0.2.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.1.0...v0.2.0) (2023-09-25)


### Features

* **BUX-207:** decode beef transaction ([#13](https://github.com/bitcoin-sv/go-paymail/issues/13)) ([e1a8273](https://github.com/bitcoin-sv/go-paymail/commit/e1a8273d79bb3753aa41a86fd0433c1be90f9f5b))

## [0.1.0](https://github.com/bitcoin-sv/go-paymail/compare/v0.0.1...v0.1.0) (2023-09-11)


### Features

* **BUX-000:** unifies github actions workflow ([#7](https://github.com/bitcoin-sv/go-paymail/issues/7)) ([b53b3a0](https://github.com/bitcoin-sv/go-paymail/commit/b53b3a04e02c152532b50b5eeff4456d64f28814))
