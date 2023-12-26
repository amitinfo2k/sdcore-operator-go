package controllers

const (
	ConfigMapVersionAnnotation = "workload.nephio.org/configMapVersion"

	// TODO(jbelamaric): Update to use ImageConfig spec.ImagePaths["upf"]
	MMEInitImage = "docker.io/amd64/busybox:stable"
	MMEImage     = "omecproject/nucleus:master-a8002eb"
	SPGWCImage   = "omecproject/spgw:master-e419062"
	HssDbImage   = "omecproject/c3po-hssdb:master-df54425"
	HSSImage     = "omecproject/c3po-hss:master-df54425"
	PCRFImage    = "omecproject/c3po-pcrf:pcrf-d58dd1c"
	PCRFDbImage  = "omecproject/c3po-pcrfdb:pcrf-d58dd1c"
	CurlImage    = "curlimages/curl:7.77.0"
)
