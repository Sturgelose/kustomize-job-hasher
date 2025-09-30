// Package main implements a job hasher function for job to append the name with a hash of the job's spec
// Only append the hash if the job has the annotation "job-hasher" set to true
package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strconv"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {
		for _, item := range items {
			if err := hasher(item); err != nil {
				return nil, err
			}
		}
		return items, nil
	}
	p := framework.SimpleProcessor{Config: nil, Filter: kio.FilterFunc(fn)}
	cmd := command.Build(p, command.StandaloneDisabled, false)
	command.AddGenerateDockerfile(cmd)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// hasher renames the job with a hash of the job
func hasher(r *yaml.RNode) error {

	// filter for Job items
	if r.GetKind() != "Job" {
		return nil
	}

	// If it has annotation "job-hasher" not set to true, return
	annotation := r.GetAnnotations()["job-hasher"]
	shouldHash, err := strconv.ParseBool(annotation)
	if err != nil || !shouldHash {
		return nil
	}

	// lookup the name field of Job items
	name := r.GetName()

	// hash the contents of the spec of the job
	spec, err := r.Pipe(yaml.Lookup("spec"))
	if err != nil {
		s, _ := r.String()
		return fmt.Errorf("%v: %s", err, s)
	}

	// marshal the spec node to bytes for hashing
	// NOTE: if the spec content's order is not consistent, the hash will be different
	specBytes, err := spec.MarshalJSON()
	if err != nil {
		s, _ := r.String()
		return fmt.Errorf("failed to marshal spec for hashing: %v: %s", err, s)
	}

	sum := sha256.Sum256(specBytes)
	hash := fmt.Sprintf("%x", sum[:])[:10]

	// append the hash to the name
	r.SetName(name + "-" + hash)

	return nil
}
