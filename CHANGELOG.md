
<a name="HEAD"></a>
## [HEAD](https://github.com/xoctopus/x/compare/v0.0.6...HEAD)

> 2026-01-28

### Doc

* add readme


<a name="v0.0.6"></a>
## [v0.0.6](https://github.com/xoctopus/x/compare/v0.0.5...v0.0.6)

> 2026-01-28

### BREAKING CHANGE


WHAT:
modify: IDGen.ID() int64 => IDGen.ID() (int64, error)
added: IDGen.MustID() int64
WHY: users can handle clock backwards
HOW: use MustID or handler error


<a name="v0.0.5"></a>
## [v0.0.5](https://github.com/xoctopus/x/compare/v0.0.4...v0.0.5)

> 2026-01-27

### Chore

* add exports for internal/factory


<a name="v0.0.4"></a>
## [v0.0.4](https://github.com/xoctopus/x/compare/v0.0.3...v0.0.4)

> 2026-01-27

### Doc

* update CHANGELOG

### Feat

* **factory:** imporve id generator


<a name="v0.0.3"></a>
## [v0.0.3](https://github.com/xoctopus/x/compare/v0.0.2...v0.0.3)

> 2025-12-25

### Feat

* add default IDGen


<a name="v0.0.2"></a>
## [v0.0.2](https://github.com/xoctopus/x/compare/v0.0.1...v0.0.2)

> 2025-12-25

### Ci

* github ci workflow

### Feat

* **sfid:** add IDGen context injector


<a name="v0.0.1"></a>
## v0.0.1

> 2025-12-25

### Refact

* move from datatypex

