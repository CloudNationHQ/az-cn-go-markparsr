# Changelog

## [1.4.1](https://github.com/CloudNationHQ/az-cn-go-markparsr/compare/v1.4.0...v1.4.1) (2025-10-10)


### Bug Fixes

* add http timeout to url validator to prevent hanging ([#13](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/13)) ([c234588](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/c234588b30a7620ec9db58145639e5326c8ed69c))
* cleanup dead code ([#17](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/17)) ([5e08862](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/5e0886203cc8e2fa11b232bff0d5883b2f248239))
* cleanup unneeded comments ([#18](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/18)) ([ec96d3b](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/ec96d3b85591db16d13751fe51b3f439c2a81c3d))
* correct error handling in resource extraction fallback ([#16](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/16)) ([ce2cdda](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/ce2cdda4bfa1c20885a6063a1e0d4deff5e620fd))
* ignore missing resources section when module has none ([#21](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/21)) ([7cc5dfd](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/7cc5dfd62a380526909d4c419128f0a2dae61a16))
* include all terraform files when extracting items ([#20](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/20)) ([6857331](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/685733196a23e6bf6ab23036303f67f72475921e))
* keep checking terraform items when sections missing ([#19](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/19)) ([e43c91b](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/e43c91bbe03ee6988821997b6c5b4e9be7c8da76))
* limit concurrent url checks in validator ([#22](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/22)) ([bfc2b34](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/bfc2b347a9785ce7457f971ebe5f78de4e61c950))
* remove hcl parser poolinh to prevent state pollution ([#15](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/15)) ([c6b9377](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/c6b937775cafae643ef278e96a3032722116c4c8))

## [1.4.0](https://github.com/CloudNationHQ/az-cn-go-markparsr/compare/v1.3.0...v1.4.0) (2025-09-18)


### Features

* add interfaces and eliminated some duplicate logic ([#11](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/11)) ([5ddfe82](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/5ddfe82eff0bb9daca451bed8820ba8d8ffa2134))
* **deps:** bump github.com/hashicorp/hcl/v2 from 2.23.0 to 2.24.0 ([#10](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/10)) ([26494ea](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/26494eaf8afd2d2761d4c4474e0cb01872cfa0fe))

## [1.3.0](https://github.com/CloudNationHQ/az-cn-go-markparsr/compare/v1.2.0...v1.3.0) (2025-09-18)


### Features

* complete refactor ([#8](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/8)) ([5579db0](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/5579db09c4ce9cbb9700e89793479a5fbcbcfff8))

## [1.2.0](https://github.com/CloudNationHQ/az-cn-go-markparsr/compare/v1.1.0...v1.2.0) (2025-04-03)


### Features

* add provider prefix configuration via options pattern ([#6](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/6)) ([40453fc](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/40453fc7841c709f0c913c26d14347e2fcba89d5))

## [1.1.0](https://github.com/CloudNationHQ/az-cn-go-markparsr/compare/v1.0.0...v1.1.0) (2025-03-22)


### Features

* add random provider support and fix deduplication in symbolic names ([#4](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/4)) ([d94d2a1](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/d94d2a1c0557d6aadf3bdb3f856b22f4d2358d67))

## 1.0.0 (2025-03-18)


### Features

* add initial package ([#1](https://github.com/CloudNationHQ/az-cn-go-markparsr/issues/1)) ([d04cb7b](https://github.com/CloudNationHQ/az-cn-go-markparsr/commit/d04cb7bd91d4231c2f4843e3cd20dbef9611ae04))
