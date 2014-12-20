# Azul3D - semver #

This package provides semantic versioning for Go packages on custom domains.

[![GoDoc](https://godoc.org/azul3d.org/semver.v1?status.svg)](https://godoc.org/azul3d.org/semver.v1)

What is it?

* [Semantic Versioning](http://semver.org/) for Go packages.
* Like [gopkg.in](http://gopkg.in), but it runs in your own Go HTTP server.
* Folder-based packages (e.g. `mydomain/my/pkg.v1` -> `github.com/myuser/my-pkg`).
* Git tags and branches (e.g. `v1 -> tag/branch v1.3.2`).
* Development branches (e.g. `import "pkg.v2-dev"`).

## Version 1.0.1

* Documentation
 * [azul3d.org/semver.v1](http://azul3d.org/semver.v1)
 * `import "azul3d.org/semver.v1"`
* Changes
 * Fixed a bug that caused branches to resolve incorrectly (see [#2](https://github.com/azul3d/semver/issues/2)).

## Version 1 #

* Documentation
 * [azul3d.org/semver.v1](http://azul3d.org/semver.v1)
 * `import "azul3d.org/semver.v1"`


