# Kustomize Transformer Job Hasher

Simple Kubernetes Fn to transform Job names to append a hash of the `spec`.
Aims to support a similar logic to `generateName` in Kubernetes Job kind, which is not supported by Kustomize.
