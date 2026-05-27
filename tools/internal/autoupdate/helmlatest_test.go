package autoupdate

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestHelmLatest(t *testing.T) {
	t.Run("Validate", func(t *testing.T) {
		type testCase struct {
			Message       string
			HelmLatest    *HelmLatest
			ExpectedError string
		}
		testCases := []testCase{
			{
				Message: "should return nil for a valid HelmLatest",
				HelmLatest: &HelmLatest{
					HelmRepo: "https://helm.cilium.io",
					Charts: map[string]map[string]Environment{
						"cilium": {
							"default": {},
						},
					},
				},
				ExpectedError: "",
			},
			{
				Message: "should return error for empty HelmRepo",
				HelmLatest: &HelmLatest{
					Charts: map[string]map[string]Environment{
						"cilium": {
							"default": {},
						},
					},
				},
				ExpectedError: "must specify HelmRepo",
			},
			{
				Message: "should return error for nil Charts",
				HelmLatest: &HelmLatest{
					HelmRepo: "https://helm.cilium.io",
				},
				ExpectedError: "must specify at least one chart in Charts",
			},
			{
				Message: "should return error for empty Charts",
				HelmLatest: &HelmLatest{
					HelmRepo: "https://helm.cilium.io",
					Charts:   map[string]map[string]Environment{},
				},
				ExpectedError: "must specify at least one chart in Charts",
			},
			{
				Message: "should return error for chart with no environments",
				HelmLatest: &HelmLatest{
					HelmRepo: "https://helm.cilium.io",
					Charts: map[string]map[string]Environment{
						"cilium": {},
					},
				},
				ExpectedError: `chart "cilium" must have at least one environment`,
			},
		}
		for _, testCase := range testCases {
			t.Run(testCase.Message, func(t *testing.T) {
				err := testCase.HelmLatest.Validate()
				if testCase.ExpectedError == "" {
					assert.Nil(t, err)
				} else {
					assert.EqualError(t, err, testCase.ExpectedError)
				}
			})
		}
	})
}

func TestHelmLatestExtractImagesFromYaml(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantErr     bool
		expectedMap map[string][]string
	}{
		{
			name:        "should return an empty map from an empty helm template",
			wantErr:     false,
			expectedMap: map[string][]string{},
			template:    `spec:`,
		},
		{
			name:    "should return a key value pair of charts and versions",
			wantErr: false,
			expectedMap: map[string][]string{
				"quay.io/cilium/cilium":       {"v1.19.4"},
				"quay.io/cilium/cilium-envoy": {"v1.36.6-1778235340"},
				"quay.io/cilium/operator-aws": {"v1.19.4"},
			},
			template: `spec:
  template:
    spec:
      containers:
      - name: cilium-agent
        image: "quay.io/cilium/cilium:v1.19.4@sha256:2eb67991eaa9368ba199c2fac2c573cb0ffdeb79184533344f42fc9a7ff6af3c"
      initContainers:
      - name: cilium-envoy
        image: "quay.io/cilium/cilium-envoy:v1.36.6-1778235340"
      - name: cilium-operator
        image: "quay.io/cilium/operator-aws:v1.19.4@sha256:9e41b3959d941a0b60ba187f5a2572305846248efb89ac59c18fd25a032f568d"
      - name: install-cni-binaries
        image: "quay.io/cilium/cilium:v1.19.4@sha256:2eb67991eaa9368ba199c2fac2c573cb0ffdeb79184533344f42fc9a7ff6af3c"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := yaml.NewDecoder(bytes.NewBufferString(tt.template))
			var dataDecoded any
			err := decoder.Decode(&dataDecoded)
			assert.NoError(t, err)

			imageMap := make(map[string][]string)

			var hl HelmLatest
			gotErr := hl.extractImagesFromYaml(dataDecoded, imageMap)
			assert.NoError(t, gotErr)
			assert.Equal(t, tt.expectedMap, imageMap)
		})
	}
}
