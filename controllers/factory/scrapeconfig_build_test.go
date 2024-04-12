package factory

import (
	"context"
	"testing"

	"github.com/VictoriaMetrics/operator/api/v1beta1"
	victoriametricsv1beta1 "github.com/VictoriaMetrics/operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGenerateScrapeConfig(t *testing.T) {
	type args struct {
		cr                    victoriametricsv1beta1.VMAgent
		m                     *victoriametricsv1beta1.VMScrapeConfig
		ssCache               *scrapesSecretsCache
		enforceNamespaceLabel string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic static cfg with basic auth",
			args: args{
				cr: victoriametricsv1beta1.VMAgent{
					Spec: victoriametricsv1beta1.VMAgentSpec{
						MinScrapeInterval: pointer.StringPtr("30s"),
						MaxScrapeInterval: pointer.String("5m"),
					},
				},
				m: &victoriametricsv1beta1.VMScrapeConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "static-1",
						Namespace: "default",
					},
					Spec: victoriametricsv1beta1.VMScrapeConfigSpec{
						ScrapeInterval: "10s",
						StaticConfigs: []victoriametricsv1beta1.StaticConfig{
							{
								Targets: []string{"http://test1.com", "http://test2.com"},
								Labels:  map[string]string{"bar": "baz"},
							},
						},
						BasicAuth: &victoriametricsv1beta1.BasicAuth{
							Username: v1.SecretKeySelector{Key: "username"},
							Password: v1.SecretKeySelector{Key: "password"},
						},
					},
				},
				ssCache: &scrapesSecretsCache{
					baSecrets: map[string]*BasicAuthCredentials{
						"scrapeConfig/default/static-1//0": {
							password: "dangerous",
							username: "admin",
						},
					},
				},
			},
			want: `job_name: scrapeConfig/default/static-1
honor_labels: false
scrape_interval: 30s
basic_auth:
  username: admin
  password: dangerous
relabel_configs: []
static_configs:
- targets:
  - http://test1.com
  - http://test2.com
  labels:
    bar: baz
`,
		},
		{
			name: "basic fileSDConfig",
			args: args{
				cr: victoriametricsv1beta1.VMAgent{
					Spec: victoriametricsv1beta1.VMAgentSpec{
						MinScrapeInterval: pointer.StringPtr("30s"),
						MaxScrapeInterval: pointer.String("5m"),
					},
				},
				m: &victoriametricsv1beta1.VMScrapeConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "file-1",
						Namespace: "default",
					},
					Spec: victoriametricsv1beta1.VMScrapeConfigSpec{
						ScrapeInterval: "10m",
						FileSDConfigs: []victoriametricsv1beta1.FileSDConfig{
							{
								Files: []string{"test1.json", "test2.json"},
							},
						},
						BasicAuth: &victoriametricsv1beta1.BasicAuth{
							Username:     v1.SecretKeySelector{Key: "username"},
							PasswordFile: "/var/run/secrets/password",
						},
					},
				},
				ssCache: &scrapesSecretsCache{
					baSecrets: map[string]*BasicAuthCredentials{
						"scrapeConfig/default/file-1//0": {
							username: "user",
						},
					},
				},
			},
			want: `job_name: scrapeConfig/default/file-1
honor_labels: false
scrape_interval: 5m
basic_auth:
  username: user
  password_file: /var/run/secrets/password
relabel_configs: []
file_sd_configs:
- files:
  - test1.json
  - test2.json
`,
		},
		{
			name: "basic httpSDConfig",
			args: args{
				cr: victoriametricsv1beta1.VMAgent{},
				m: &victoriametricsv1beta1.VMScrapeConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httpsd-1",
						Namespace: "default",
					},
					Spec: victoriametricsv1beta1.VMScrapeConfigSpec{
						HTTPSDConfigs: []victoriametricsv1beta1.HTTPSDConfig{
							{
								URL:      "http://www.test1.com",
								ProxyURL: pointer.String("http://www.proxy.com"),
							},
							{
								URL: "http://www.test2.com",
								Authorization: &victoriametricsv1beta1.Authorization{
									Type:        "Bearer",
									Credentials: &v1.SecretKeySelector{Key: "cred"},
								},
								TLSConfig: &victoriametricsv1beta1.TLSConfig{
									CA: v1beta1.SecretOrConfigMap{
										Secret: &v1.SecretKeySelector{
											Key: "ca",
											LocalObjectReference: v1.LocalObjectReference{
												Name: "tls-secret",
											},
										},
									},
									Cert: v1beta1.SecretOrConfigMap{Secret: &v1.SecretKeySelector{Key: "cert"}},
								},
							},
						},
					},
				},
				ssCache: &scrapesSecretsCache{
					baSecrets: map[string]*BasicAuthCredentials{
						"scrapeConfig/default/file-1//0": {
							username: "user",
						},
					},
					authorizationSecrets: map[string]string{
						"scrapeConfig/default/httpsd-1/httpsd/1": "auth-secret",
					},
				},
			},
			want: `job_name: scrapeConfig/default/httpsd-1
honor_labels: false
relabel_configs: []
http_sd_configs:
- url: http://www.test1.com
  proxy_url: http://www.proxy.com
- url: http://www.test2.com
  authorization:
    type: Bearer
    credentials: auth-secret
  tls_config:
    insecure_skip_verify: false
    ca_file: /etc/vmagent-tls/certs/default_tls-secret_ca
`,
		},
		{
			name: "basic kubernetesSDConfig",
			args: args{
				cr: victoriametricsv1beta1.VMAgent{},
				m: &victoriametricsv1beta1.VMScrapeConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kubernetesSDConfig-1",
						Namespace: "default",
					},
					Spec: victoriametricsv1beta1.VMScrapeConfigSpec{
						KubernetesSDConfigs: []victoriametricsv1beta1.KubernetesSDConfig{
							{
								APIServer:      pointer.String("http://127.0.0.1:6443"),
								Role:           "pod",
								AttachMetadata: victoriametricsv1beta1.AttachMetadata{Node: pointer.Bool(true)},
								Selectors: []victoriametricsv1beta1.K8SSelectorConfig{
									{
										Role:  "pod",
										Label: "app/instance",
										Field: "test",
									},
								},
								TLSConfig: &victoriametricsv1beta1.TLSConfig{
									InsecureSkipVerify: true,
								},
							},
							{
								APIServer: pointer.String("http://127.0.0.1:6443"),
								Role:      "node",
								Selectors: []victoriametricsv1beta1.K8SSelectorConfig{
									{
										Role:  "node",
										Label: "kubernetes.io/os",
										Field: "linux",
									},
								},
							},
						},
					},
				},
				ssCache: &scrapesSecretsCache{},
			},
			want: `job_name: scrapeConfig/default/kubernetesSDConfig-1
honor_labels: false
relabel_configs: []
kubernetes_sd_configs:
- api_server: http://127.0.0.1:6443
  role: pod
  tls_config:
    insecure_skip_verify: true
  attach_metadata: true
  selectors:
  - role: pod
    label: app/instance
    field: test
- api_server: http://127.0.0.1:6443
  role: node
  selectors:
  - role: node
    label: kubernetes.io/os
    field: linux
`,
		},
		{
			name: "mixed",
			args: args{
				cr: victoriametricsv1beta1.VMAgent{},
				m: &victoriametricsv1beta1.VMScrapeConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "mixconfigs-1",
						Namespace: "default",
					},
					Spec: victoriametricsv1beta1.VMScrapeConfigSpec{
						ConsulSDConfigs: []victoriametricsv1beta1.ConsulSDConfig{
							{
								Server:     "localhost:8500",
								TokenRef:   &v1.SecretKeySelector{Key: "consul_token"},
								Datacenter: pointer.String("dc1"),
								NodeMeta:   map[string]string{"worker": "1"},
							},
						},
						DNSSDConfigs: []victoriametricsv1beta1.DNSSDConfig{
							{
								Names: []string{"vmagent-0.vmagent.default.svc.cluster.local"},
								Port:  pointer.Int(8429),
							},
						},
						EC2SDConfigs: []victoriametricsv1beta1.EC2SDConfig{
							{
								Region: pointer.String("us-west-2"),
								Port:   pointer.Int(9404),
								Filters: []*v1beta1.EC2Filter{{
									Name:   "instance-id",
									Values: []string{"i-98765432109876543", "i-12345678901234567"},
								}},
							},
						},
						AzureSDConfigs: []victoriametricsv1beta1.AzureSDConfig{
							{
								Environment:    pointer.String("AzurePublicCloud"),
								SubscriptionID: "1",
								TenantID:       pointer.String("u1"),
								ResourceGroup:  pointer.String("rg1"),
								Port:           pointer.Int(80),
							},
						},
						GCESDConfigs: []victoriametricsv1beta1.GCESDConfig{
							{
								Project:      "eu-project",
								Zone:         "zone-a",
								TagSeparator: pointer.String("/"),
							},
						},
						OpenStackSDConfigs: []victoriametricsv1beta1.OpenStackSDConfig{
							{
								Role:             "instance",
								IdentityEndpoint: pointer.String("http://localhost:5000/v3"),
								Username:         pointer.String("user1"),
								UserID:           pointer.String("1"),
								Password:         &v1.SecretKeySelector{Key: "pass"},
								ProjectName:      pointer.String("poc"),
								AllTenants:       pointer.Bool(true),
								DomainName:       pointer.String("default"),
							},
						},
						DigitalOceanSDConfigs: []victoriametricsv1beta1.DigitalOceanSDConfig{
							{
								OAuth2: &victoriametricsv1beta1.OAuth2{
									Scopes:         []string{"scope-1"},
									TokenURL:       "http://some-token-url",
									EndpointParams: map[string]string{"timeout": "5s"},
									ClientID: victoriametricsv1beta1.SecretOrConfigMap{
										Secret: &v1.SecretKeySelector{
											Key:                  "bearer",
											LocalObjectReference: v1.LocalObjectReference{Name: "access-secret"},
										},
									},
								},
							},
						},
					},
				},
				ssCache: &scrapesSecretsCache{
					oauth2Secrets: map[string]*oauthCreds{
						"scrapeConfig/default/mixconfigs-1/digitaloceansd/0": {clientSecret: "some-secret", clientID: "some-id"},
					},
				},
			},
			want: `job_name: scrapeConfig/default/mixconfigs-1
honor_labels: false
relabel_configs: []
consul_sd_configs:
- server: localhost:8500
  datacenter: dc1
  node_meta:
    worker: "1"
dns_sd_configs:
- names:
  - vmagent-0.vmagent.default.svc.cluster.local
  port: 8429
ec2_sd_configs:
- region: us-west-2
  port: 9404
  filters:
  - name: instance-id
    values:
    - i-98765432109876543
    - i-12345678901234567
azure_sd_configs:
- environment: AzurePublicCloud
  subscription_id: "1"
  tenant_id: u1
  resource_group: rg1
  port: 80
gce_sd_configs:
- project: eu-project
  zone: zone-a
  tag_separator: /
openstack_sd_configs:
- role: instance
  region: ""
  identity_endpoint: http://localhost:5000/v3
  username: user1
  userid: "1"
  domain_name: default
  project_name: poc
  all_tenants: true
digitalocean_sd_configs:
- oauth2:
    client_id: some-id
    scopes:
    - scope-1
    endpoint_params:
      timeout: 5s
    token_url: http://some-token-url
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateScrapeConfig(context.Background(), &tt.args.cr, tt.args.m, tt.args.ssCache, tt.args.enforceNamespaceLabel)
			gotBytes, err := yaml.Marshal(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !assert.Equal(t, tt.want, string(gotBytes)) {
				t.Errorf("generateScrapeConfig() = \n%v, want \n%v", string(gotBytes), tt.want)
			}
		})
	}
}
