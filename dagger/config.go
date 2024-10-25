package main

const DaggerEngineConfig = `
debug = true
insecure-entitlements = ["security.insecure"]

[registry."docker.io"]
mirrors = ["%s"]
http = true
insecure = true

[registry."%s"]
http = true
insecure = true
`

const RegistryConfig = `
version: 0.1
log:
  accesslog:
    disabled: true
  level: warn
  formatter: text

storage:
  filesystem:
    rootdirectory: /var/lib/registry

http:
  addr: :5000
  relativeurls: false
  draintimeout: 60s

proxy:
  remoteurl: https://registry-1.docker.io
`
