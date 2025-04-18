package autogen

import (
	"encoding/json"
	"testing"

	policiesv1alpha1 "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

func TestGenerateRuleForControllers(t *testing.T) {
	tests := []struct {
		name          string
		controllers   sets.Set[string]
		policySpec    []byte
		generatedRule map[string]policiesv1alpha1.ValidatingPolicyAutogen
	}{
		{
			name:        "autogen rule for deployments",
			controllers: sets.New("deployments"),
			policySpec: []byte(`{
				"matchConstraints": {
					"resourceRules": [
						{
							"apiGroups": [
								""
							],
							"apiVersions": [
								"v1"
							],
							"operations": [
								"CREATE",
								"UPDATE"
							],
							"resources": [
								"pods"
							]
						}
					]
				},
				"validations": [
					{
						"expression": "object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"
					}
				]
			}`),
			generatedRule: map[string]policiesv1alpha1.ValidatingPolicyAutogen{
				"": {
					Spec: &policiesv1alpha1.ValidatingPolicySpec{
						MatchConstraints: &admissionregistrationv1.MatchResources{
							ResourceRules: []admissionregistrationv1.NamedRuleWithOperations{
								{
									RuleWithOperations: admissionregistrationv1.RuleWithOperations{
										Operations: []admissionregistrationv1.OperationType{
											admissionregistrationv1.Create,
											admissionregistrationv1.Update,
										},
										Rule: admissionregistrationv1.Rule{
											APIGroups:   []string{"apps"},
											APIVersions: []string{"v1"},
											Resources:   []string{"deployments"},
										},
									},
								},
							},
						},
						Validations: []admissionregistrationv1.Validation{
							{
								Expression: "object.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)",
							},
						},
					},
				},
			},
		},
		{
			name:        "autogen rule for deployments and daemonsets",
			controllers: sets.New("deployments", "daemonsets"),
			policySpec: []byte(`{
				"matchConstraints": {
					"resourceRules": [
						{
							"apiGroups": [
								""
							],
							"apiVersions": [
								"v1"
							],
							"operations": [
								"CREATE",
								"UPDATE"
							],
							"resources": [
								"pods"
							]
						}
					]
				},
				"validations": [
					{
						"expression": "object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"
					}
				]
			}`),
			generatedRule: map[string]policiesv1alpha1.ValidatingPolicyAutogen{
				"": {
					Spec: &policiesv1alpha1.ValidatingPolicySpec{
						MatchConstraints: &admissionregistrationv1.MatchResources{
							ResourceRules: []admissionregistrationv1.NamedRuleWithOperations{
								{
									RuleWithOperations: admissionregistrationv1.RuleWithOperations{
										Operations: []admissionregistrationv1.OperationType{
											admissionregistrationv1.Create,
											admissionregistrationv1.Update,
										},
										Rule: admissionregistrationv1.Rule{
											APIGroups:   []string{"apps"},
											APIVersions: []string{"v1"},
											Resources:   []string{"daemonsets", "deployments"},
										},
									},
								},
							},
						},
						Validations: []admissionregistrationv1.Validation{
							{
								Expression: "object.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)",
							},
						},
					},
				},
			},
		},
		{
			name:        "autogen rule for deployments, daemonsets, statefulsets and replicasets",
			controllers: sets.New("deployments", "daemonsets", "statefulsets", "replicasets"),
			policySpec: []byte(`{
				"matchConstraints": {
					"resourceRules": [
						{
							"apiGroups": [
								""
							],
							"apiVersions": [
								"v1"
							],
							"operations": [
								"CREATE",
								"UPDATE"
							],
							"resources": [
								"pods"
							]
						}
					]
				},
				"matchConditions": [
					{
						"name": "only for production",
						"expression": "has(object.metadata.labels) && has(object.metadata.labels.prod) && object.metadata.labels.prod == 'true'"
					}
				],
				"validations": [
					{
						"expression": "object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"
					}
				]
			}`),
			generatedRule: map[string]policiesv1alpha1.ValidatingPolicyAutogen{
				"": {
					Spec: &policiesv1alpha1.ValidatingPolicySpec{
						MatchConstraints: &admissionregistrationv1.MatchResources{
							ResourceRules: []admissionregistrationv1.NamedRuleWithOperations{
								{
									RuleWithOperations: admissionregistrationv1.RuleWithOperations{
										Operations: []admissionregistrationv1.OperationType{
											admissionregistrationv1.Create,
											admissionregistrationv1.Update,
										},
										Rule: admissionregistrationv1.Rule{
											APIGroups:   []string{"apps"},
											APIVersions: []string{"v1"},
											Resources:   []string{"daemonsets", "deployments", "replicasets", "statefulsets"},
										},
									},
								},
							},
						},
						MatchConditions: []admissionregistrationv1.MatchCondition{
							{
								Name:       "autogen-only for production",
								Expression: "!((object.apiVersion == 'apps/v1' && object.kind =='DaemonSet') || (object.apiVersion == 'apps/v1' && object.kind =='Deployment') || (object.apiVersion == 'apps/v1' && object.kind =='ReplicaSet') || (object.apiVersion == 'apps/v1' && object.kind =='StatefulSet')) || (has(object.spec.template.metadata.labels) && has(object.spec.template.metadata.labels.prod) && object.spec.template.metadata.labels.prod == 'true')",
							},
						},
						Validations: []admissionregistrationv1.Validation{
							{
								Expression: "object.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)",
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var spec policiesv1alpha1.ValidatingPolicySpec
			err := json.Unmarshal(test.policySpec, &spec)
			assert.NoError(t, err)
			genRule, err := generateRuleForControllers(spec, test.controllers)
			assert.NoError(t, err)
			assert.Equal(t, test.generatedRule, genRule)
		})
	}
}

func TestGenerateCronJobRule(t *testing.T) {
	tests := []struct {
		policySpec    []byte
		generatedRule map[string]policiesv1alpha1.ValidatingPolicyAutogen
	}{
		{
			policySpec: []byte(`{
    "matchConstraints": {
        "resourceRules": [
            {
                "apiGroups": [
                    ""
                ],
                "apiVersions": [
                    "v1"
                ],
                "operations": [
                    "CREATE",
                    "UPDATE"
                ],
                "resources": [
                    "pods"
                ]
            }
        ]
    },
    "validations": [
        {
            "expression": "object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"
        }
    ]
}`),
			generatedRule: map[string]policiesv1alpha1.ValidatingPolicyAutogen{
				"cronjobs": {
					Spec: &policiesv1alpha1.ValidatingPolicySpec{
						MatchConstraints: &admissionregistrationv1.MatchResources{
							ResourceRules: []admissionregistrationv1.NamedRuleWithOperations{
								{
									RuleWithOperations: admissionregistrationv1.RuleWithOperations{
										Operations: []admissionregistrationv1.OperationType{
											admissionregistrationv1.Create,
											admissionregistrationv1.Update,
										},
										Rule: admissionregistrationv1.Rule{
											APIGroups:   []string{"batch"},
											APIVersions: []string{"v1"},
											Resources:   []string{"cronjobs"},
										},
									},
								},
							},
						},
						Validations: []admissionregistrationv1.Validation{
							{
								Expression: "object.spec.jobTemplate.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)",
							},
						},
					},
				},
			},
		},
		{
			policySpec: []byte(`{
    "matchConstraints": {
        "resourceRules": [
            {
                "apiGroups": [
                    ""
                ],
                "apiVersions": [
                    "v1"
                ],
                "operations": [
                    "CREATE",
                    "UPDATE"
                ],
                "resources": [
                    "pods"
                ]
            }
        ]
    },
	"matchConditions": [
		{
	        "name": "only for production",
			"expression": "has(object.metadata.labels) && has(object.metadata.labels.prod) && object.metadata.labels.prod == 'true'"
	    }
	],
    "validations": [
        {
            "expression": "object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"
        }
    ]
}`),
			generatedRule: map[string]policiesv1alpha1.ValidatingPolicyAutogen{
				"cronjobs": {
					Spec: &policiesv1alpha1.ValidatingPolicySpec{
						MatchConstraints: &admissionregistrationv1.MatchResources{
							ResourceRules: []admissionregistrationv1.NamedRuleWithOperations{
								{
									RuleWithOperations: admissionregistrationv1.RuleWithOperations{
										Operations: []admissionregistrationv1.OperationType{
											admissionregistrationv1.Create,
											admissionregistrationv1.Update,
										},
										Rule: admissionregistrationv1.Rule{
											APIGroups:   []string{"batch"},
											APIVersions: []string{"v1"},
											Resources:   []string{"cronjobs"},
										},
									},
								},
							},
						},
						MatchConditions: []admissionregistrationv1.MatchCondition{
							{
								Name:       "autogen-cronjobs-only for production",
								Expression: "!((object.apiVersion == 'batch/v1' && object.kind =='CronJob')) || (has(object.spec.jobTemplate.spec.template.metadata.labels) && has(object.spec.jobTemplate.spec.template.metadata.labels.prod) && object.spec.jobTemplate.spec.template.metadata.labels.prod == 'true')",
							},
						},
						Validations: []admissionregistrationv1.Validation{
							{
								Expression: "object.spec.jobTemplate.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)",
							},
						},
					},
				},
			},
		},
		{
			policySpec: []byte(`{
    "matchConstraints": {
        "resourceRules": [
            {
                "apiGroups": [
                    ""
                ],
                "apiVersions": [
                    "v1"
                ],
                "operations": [
                    "CREATE",
                    "UPDATE"
                ],
                "resources": [
                    "pods"
                ]
            }
        ]
    },
    "variables": [
        {
            "name": "environment",
            "expression": "has(object.metadata.labels) && 'env' in object.metadata.labels && object.metadata.labels['env'] == 'prod'"
        }
    ],
    "validations": [
        {
            "expression": "variables.environment == true",
            "message": "labels must be env=prod"
        }
    ]
}`),
			generatedRule: map[string]policiesv1alpha1.ValidatingPolicyAutogen{
				"cronjobs": {
					Spec: &policiesv1alpha1.ValidatingPolicySpec{
						MatchConstraints: &admissionregistrationv1.MatchResources{
							ResourceRules: []admissionregistrationv1.NamedRuleWithOperations{
								{
									RuleWithOperations: admissionregistrationv1.RuleWithOperations{
										Operations: []admissionregistrationv1.OperationType{
											admissionregistrationv1.Create,
											admissionregistrationv1.Update,
										},
										Rule: admissionregistrationv1.Rule{
											APIGroups:   []string{"batch"},
											APIVersions: []string{"v1"},
											Resources:   []string{"cronjobs"},
										},
									},
								},
							},
						},
						Variables: []admissionregistrationv1.Variable{
							{
								Name:       "environment",
								Expression: "has(object.spec.jobTemplate.spec.template.metadata.labels) && 'env' in object.spec.jobTemplate.spec.template.metadata.labels && object.spec.jobTemplate.spec.template.metadata.labels['env'] == 'prod'",
							},
						},
						Validations: []admissionregistrationv1.Validation{
							{
								Expression: "variables.environment == true",
								Message:    "labels must be env=prod",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var spec policiesv1alpha1.ValidatingPolicySpec
			err := json.Unmarshal(tt.policySpec, &spec)
			assert.NoError(t, err)
			genRule, err := generateRuleForControllers(spec, sets.New("cronjobs"))
			assert.NoError(t, err)
			assert.Equal(t, tt.generatedRule, genRule)
		})
	}
}

func TestUpdateGenRuleByte(t *testing.T) {
	tests := []struct {
		pbyte  []byte
		config string
		want   []byte
	}{
		{
			pbyte:  []byte("object.spec"),
			config: "deployments",
			want:   []byte("object.spec.template.spec"),
		},
		{
			pbyte:  []byte("oldObject.spec"),
			config: "deployments",
			want:   []byte("oldObject.spec.template.spec"),
		},
		{
			pbyte:  []byte("object.spec"),
			config: "cronjobs",
			want:   []byte("object.spec.jobTemplate.spec.template.spec"),
		},
		{
			pbyte:  []byte("oldObject.spec"),
			config: "cronjobs",
			want:   []byte("oldObject.spec.jobTemplate.spec.template.spec"),
		},
		{
			pbyte:  []byte("object.metadata"),
			config: "deployments",
			want:   []byte("object.spec.template.metadata"),
		},
		{
			pbyte:  []byte("object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"),
			config: "cronjobs",
			want:   []byte("object.spec.jobTemplate.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"),
		},
		{
			pbyte:  []byte("object.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"),
			config: "deployments",
			want:   []byte("object.spec.template.spec.containers.all(container, has(container.securityContext) && has(container.securityContext.allowPrivilegeEscalation) && container.securityContext.allowPrivilegeEscalation == false)"),
		},
	}
	for _, tt := range tests {
		got := updateFields(tt.pbyte, replacementsMap[configsMap[tt.config].replacementsRef]...)
		assert.Equal(t, tt.want, got)
	}
}
