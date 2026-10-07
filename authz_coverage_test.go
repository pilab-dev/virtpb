package virtpb_test

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	_ "go.pilab.hu/cloud/virtpb/v3/pilab/agent/v2"
	_ "go.pilab.hu/cloud/virtpb/v3/pilab/cli/v1"
	commonv1 "go.pilab.hu/cloud/virtpb/v3/pilab/common/v1"
	_ "go.pilab.hu/cloud/virtpb/v3/pilab/director/v2"
	_ "go.pilab.hu/cloud/virtpb/v3/pilab/frontend/v1"
	_ "go.pilab.hu/cloud/virtpb/v3/pilab/network/v2"
	_ "go.pilab.hu/cloud/virtpb/v3/pilab/pivirtd/v1"
	_ "go.pilab.hu/cloud/virtpb/v3/pilab/resource/v1"
)

// unguarded lists the proto packages whose services are not an authorisation
// boundary: pivirtd is the host executor, reachable only on a local unix
// socket, and never learns users or tenants (control-plane-architecture.md
// §2.1). Everything else must declare an authz rule on every method.
var unguarded = map[string]bool{
	"pilab.pivirtd.v1":  true,
	"pilab.resource.v1": true,
}

var verbs = map[string]bool{
	"get": true, "list": true, "watch": true, "create": true, "update": true,
	"delete": true, "power": true, "console": true, "migrate": true,
	"snapshot": true, "revert": true, "attach": true, "admin": true,
}

func TestEveryGuardedMethodDeclaresAuthz(t *testing.T) {
	var checked int
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if !strings.HasPrefix(string(fd.Package()), "pilab.") || unguarded[string(fd.Package())] {
			return true
		}
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			for j := 0; j < svc.Methods().Len(); j++ {
				m := svc.Methods().Get(j)
				name := string(m.FullName())
				if !proto.HasExtension(m.Options(), commonv1.E_Authz) {
					t.Errorf("%s: missing (pilab.common.v1.authz) option", name)
					continue
				}
				rule, _ := proto.GetExtension(m.Options(), commonv1.E_Authz).(*commonv1.AuthzRule)
				checked++
				switch rule.GetAuthn() {
				case commonv1.Authn_AUTHN_UNSPECIFIED:
					t.Errorf("%s: authn must be set explicitly", name)
				case commonv1.Authn_AUTHN_USER:
					if rule.GetKind() == "" || !verbs[rule.GetVerb()] {
						t.Errorf("%s: AUTHN_USER needs a kind and a known verb, got kind=%q verb=%q", name, rule.GetKind(), rule.GetVerb())
					}
				default:
					if rule.GetKind() != "" || rule.GetVerb() != "" {
						t.Errorf("%s: %s must not set kind/verb", name, rule.GetAuthn())
					}
				}
			}
		}
		return true
	})
	if checked == 0 {
		t.Fatal("no guarded methods found; the registry walk is broken")
	}
	t.Logf("checked %d methods", checked)
}
