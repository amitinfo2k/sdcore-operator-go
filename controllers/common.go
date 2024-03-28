package controllers

const (
	ConfigMapVersionAnnotation = "workload.nephio.org/configMapVersion"
	Config4GAnnotation         = "field.cattle.io/workloadMetrics"
	// TODO(jbelamaric): Update to use ImageConfig spec.ImagePaths["upf"]
	MMEInitImage = "docker.io/amd64/busybox:stable"
	MMEImage     = "omecproject/nucleus:master-a8002eb"
	//SPGWCImage    = "omecproject/spgw:master-e419062"
	SPGWCImage    = "amitinfo2k/spgw:1.0.1"
	HssDbImage    = "omecproject/c3po-hssdb:master-df54425"
	HSSImage      = "omecproject/c3po-hss:master-df54425"
	PCRFImage     = "omecproject/c3po-pcrf:pcrf-d58dd1c"
	PCRFDbImage   = "omecproject/c3po-pcrfdb:pcrf-d58dd1c"
	CurlImage     = "curlimages/curl:7.77.0"
	Config4GImage = "omecproject/webui"
)
