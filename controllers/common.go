package controllers

const (
	ConfigMapVersionAnnotation = "workload.nephio.org/configMapVersion"

	// TODO(jbelamaric): Update to use ImageConfig spec.ImagePaths["upf"]
	MMEInitImage = "docker.io/amd64/busybox:stable"
	MMEImage     = "omecproject/nucleus:master-a8002eb"
	SPGWCImage   = "omecproject/spgw:master-e419062"
)
