package zabbix

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestRoleGetServiceAccessTags(t *testing.T) {
	c := queryFixtureClient(t, "role.get", `{"roleids":["42","43"],"selectRules":"extend"}`,
		`[{"roleid":"42","rules":{"services.read.tag":{"tag":"env","value":"prod"},"services.write.tag":{"tag":"","value":""}}},{"roleid":"43","rules":{"api.access":"1"}}]`)
	roles, err := c.Role().Get(context.Background(), RoleGetParams{RoleIDs: []string{"42", "43"}, SelectRules: "extend"})
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 2 || roles[0].Rules.Value == nil || roles[1].Rules.Value == nil {
		t.Fatalf("selected roles = %#v", roles)
	}
	if got := roles[0].Rules.Value.ServicesReadTag; got == nil || *got != (ServiceAccessTag{Tag: "env", Value: "prod"}) {
		t.Fatalf("read tag = %#v", got)
	}
	if got := roles[0].Rules.Value.ServicesWriteTag; got == nil || *got != (ServiceAccessTag{}) {
		t.Fatalf("cleared write tag = %#v", got)
	}
	if rules := roles[1].Rules.Value; rules.ServicesReadTag != nil || rules.ServicesWriteTag != nil {
		t.Fatalf("omitted tag conditions = %#v", rules)
	}
}

func TestRoleCreateServiceAccessTags(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value *string
		want  string
	}{
		{
			name: "omitted values",
			want: `[{"name":"operators","type":1,"rules":{"services.read.mode":0,"services.read.tag":{"tag":"read"},"services.write.mode":0,"services.write.tag":{"tag":"write"}}}]`,
		},
		{
			name:  "empty values",
			value: writeTestPointer(""),
			want:  `[{"name":"operators","type":1,"rules":{"services.read.mode":0,"services.read.tag":{"tag":"read","value":""},"services.write.mode":0,"services.write.tag":{"tag":"write","value":""}}}]`,
		},
		{
			name:  "nonempty values",
			value: writeTestPointer("prod"),
			want:  `[{"name":"operators","type":1,"rules":{"services.read.mode":0,"services.read.tag":{"tag":"read","value":"prod"},"services.write.mode":0,"services.write.tag":{"tag":"write","value":"prod"}}}]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "role.create", tc.want, `{"roleids":["42"]}`)
			ids, err := c.Role().Create(context.Background(), RoleCreateParams{
				Name: "operators", Type: 1,
				Rules: &RoleRulesInput{
					ServicesReadMode:  writeTestPointer(0),
					ServicesReadTag:   &ServiceAccessTagInput{Tag: "read", Value: tc.value},
					ServicesWriteMode: writeTestPointer(0),
					ServicesWriteTag:  &ServiceAccessTagInput{Tag: "write", Value: tc.value},
				},
			})
			if err != nil || !reflect.DeepEqual(ids, []string{"42"}) {
				t.Fatalf("created role IDs = %v, error = %v", ids, err)
			}
		})
	}
}

func TestRoleUpdateClearsServiceAccessTags(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rules RoleRulesInput
		want  string
	}{
		{
			name:  "read tag with omitted value",
			rules: RoleRulesInput{ServicesReadTag: &ServiceAccessTagInput{Tag: ""}},
			want:  `[{"roleid":"42","rules":{"services.read.tag":{"tag":""}}}]`,
		},
		{
			name:  "read tag with empty value",
			rules: RoleRulesInput{ServicesReadTag: &ServiceAccessTagInput{Tag: "", Value: writeTestPointer("")}},
			want:  `[{"roleid":"42","rules":{"services.read.tag":{"tag":"","value":""}}}]`,
		},
		{
			name:  "write tag with omitted value",
			rules: RoleRulesInput{ServicesWriteTag: &ServiceAccessTagInput{Tag: ""}},
			want:  `[{"roleid":"42","rules":{"services.write.tag":{"tag":""}}}]`,
		},
		{
			name:  "write tag with empty value",
			rules: RoleRulesInput{ServicesWriteTag: &ServiceAccessTagInput{Tag: "", Value: writeTestPointer("")}},
			want:  `[{"roleid":"42","rules":{"services.write.tag":{"tag":"","value":""}}}]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "role.update", tc.want, `{"roleids":["42"]}`)
			ids, err := c.Role().Update(context.Background(), RoleUpdateParams{RoleID: "42", Rules: &tc.rules})
			if err != nil {
				t.Fatalf("clear service access tag: %v", err)
			}
			if len(ids) != 1 || ids[0] != "42" {
				t.Fatalf("updated role IDs = %v, want [42]", ids)
			}
		})
	}
}

func TestRoleUpdateOmitsServiceAccessTags(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params RoleUpdateParams
		want   string
	}{
		{
			name:   "omitted rules",
			params: RoleUpdateParams{RoleID: "42", Name: writeTestPointer("operators")},
			want:   `[{"roleid":"42","name":"operators"}]`,
		},
		{
			name:   "unrelated rule",
			params: RoleUpdateParams{RoleID: "42", Rules: &RoleRulesInput{APIAccess: writeTestPointer(0)}},
			want:   `[{"roleid":"42","rules":{"api.access":0}}]`,
		},
		{
			name: "explicit modes and empty lists",
			params: RoleUpdateParams{RoleID: "42", Rules: &RoleRulesInput{
				ServicesReadMode: writeTestPointer(0), ServicesReadList: []ServiceRef{},
				ServicesWriteMode: writeTestPointer(0), ServicesWriteList: []ServiceRef{},
			}},
			want: `[{"roleid":"42","rules":{"services.read.mode":0,"services.read.list":[],"services.write.mode":0,"services.write.list":[]}}]`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "role.update", tc.want, `{"roleids":["42"]}`)
			ids, err := c.Role().Update(context.Background(), tc.params)
			if err != nil || !reflect.DeepEqual(ids, []string{"42"}) {
				t.Fatalf("updated role IDs = %v, error = %v", ids, err)
			}
		})
	}
}

func TestRoleServiceAccessTagValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rules RoleRulesInput
	}{
		{"read tag", RoleRulesInput{ServicesReadTag: &ServiceAccessTagInput{Tag: "", Value: writeTestPointer("prod")}}},
		{"write tag", RoleRulesInput{ServicesWriteTag: &ServiceAccessTagInput{Tag: "", Value: writeTestPointer("prod")}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				t.Errorf("invalid tag condition reached upstream: %s", req.Method)
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"roleids":["42"]}}`, req.ID)
			})
			ids, err := c.Role().Create(context.Background(), RoleCreateParams{Name: "operators", Type: 1, Rules: &tc.rules})
			assertQueryError(t, err, "role.create", ErrorInvalidParams)
			if ids != nil {
				t.Fatalf("invalid create returned IDs: %v", ids)
			}
			ids, err = c.Role().Update(context.Background(), RoleUpdateParams{RoleID: "42", Rules: &tc.rules})
			assertQueryError(t, err, "role.update", ErrorInvalidParams)
			if ids != nil {
				t.Fatalf("invalid update returned IDs: %v", ids)
			}
		})
	}

	t.Run("entity tag still requires a name", func(t *testing.T) {
		c := responseServer(t, func(req rpcRequest) string {
			t.Errorf("empty entity tag reached upstream: %s", req.Method)
			return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"hostids":["42"]}}`, req.ID)
		})
		ids, err := c.Host().Update(context.Background(), HostUpdateParams{HostID: "42", Tags: []TagInput{{Tag: ""}}})
		assertQueryError(t, err, "host.update", ErrorInvalidParams)
		if ids != nil {
			t.Fatalf("invalid host update returned IDs: %v", ids)
		}
	})
}
