Demo App

## Small tricks
- Version embedded into the binary
- Static HTML served by the binary, and optionally can be exported to a directory if you want it to be served by a web server or CDN

## Dependency injection
Demo app uses plugin mechanism to connect different packages.

## Dependency hierarchy
- Util 
- rtenv
- config
- Authn
- Authz
- Listeners (server)
- Data


## Data domains:
### Audit log
Append only

### Config
Typical backend: git, s3, file system dir 
Transaction log, with signer info
Cumulative state hash
Time machine function
No anonymization support
Concurrent modification isn't permitted
Can't reference to other domain
Changes delayed (periodically refresh)

### Generic data
Uses registry pattern for every major entity
Multiple version of a record
Anonymization support
Field level encryption support
Explicit previous record

