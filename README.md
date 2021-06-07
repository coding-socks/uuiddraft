# UUID Draft

Draft Prototype for UUIDv6 and beyond.

## Introduction

> A universally unique identifier (UUID) is a 128-bit label used for information in computer systems.
 
Source: https://en.wikipedia.org/wiki/Universally_unique_identifier

It is based on the 1st iteration of [draft-peabody-dispatch-new-uuid-format][draft-peabody-dispatch-new-uuid-format-01]. This document is only an [Internet-Draft](https://tools.ietf.org/html/rfc2026#section-2.2).

The goal is to provide implementation for these documents and during the implementation provide feedback for them.

## Production readiness

**This project is still in alpha phase.** In this stage the public API can change between days.

Beta version will be considered when the feature set covers most of the documents the implementation is based on, and the public API is reached a mature state.

Stable version will be considered only if enough positive feedback is gathered to lock the public API and all document the implementation is based on became ["Internet Standard"](https://tools.ietf.org/html/rfc2026#section-4.1.3).

## Documents

### IETF Documents

- [RFC4122](https://datatracker.ietf.org/doc/html/rfc4122)
- [draft-peabody-dispatch-new-uuid-format][draft-peabody-dispatch-new-uuid-format-01]

Huge thanks to Brad Peabody, Kyzer Davis, and other who contributed to these documents for their work. 

[draft-peabody-dispatch-new-uuid-format-01]: https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html

## Inspiration

Inspiration for the implementation comes from the following places:

- https://pkg.go.dev/github.com/google/uuid
- https://pkg.go.dev/github.com/gofrs/uuid
