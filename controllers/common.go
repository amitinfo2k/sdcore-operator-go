package controllers

const (
	ConfigMapVersionAnnotation = "workload.nephio.org/configMapVersion"

	// TODO(jbelamaric): Update to use ImageConfig spec.ImagePaths["upf"]
	MMEInitImage = "docker.io/amd64/busybox:stable"
	MMEImage     = "omecproject/nucleus:master-a8002eb"
	SPGWCImage   = "omecproject/spgw:master-e419062"
	HssDbImage   = "omecproject/c3po-hssdb:master-df54425"
	HSSImage     = "omecproject/c3po-hss:master-df54425"
	CurlImage    = "curlimages/curl:7.77.0"
)
