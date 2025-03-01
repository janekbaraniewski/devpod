// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	context "context"

	managementv1 "github.com/loft-sh/api/v4/pkg/apis/management/v1"
	scheme "github.com/loft-sh/api/v4/pkg/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// WorkspaceAccessKeysGetter has a method to return a WorkspaceAccessKeyInterface.
// A group's client should implement this interface.
type WorkspaceAccessKeysGetter interface {
	WorkspaceAccessKeys() WorkspaceAccessKeyInterface
}

// WorkspaceAccessKeyInterface has methods to work with WorkspaceAccessKey resources.
type WorkspaceAccessKeyInterface interface {
	Create(ctx context.Context, workspaceAccessKey *managementv1.WorkspaceAccessKey, opts metav1.CreateOptions) (*managementv1.WorkspaceAccessKey, error)
	Update(ctx context.Context, workspaceAccessKey *managementv1.WorkspaceAccessKey, opts metav1.UpdateOptions) (*managementv1.WorkspaceAccessKey, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, workspaceAccessKey *managementv1.WorkspaceAccessKey, opts metav1.UpdateOptions) (*managementv1.WorkspaceAccessKey, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*managementv1.WorkspaceAccessKey, error)
	List(ctx context.Context, opts metav1.ListOptions) (*managementv1.WorkspaceAccessKeyList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *managementv1.WorkspaceAccessKey, err error)
	WorkspaceAccessKeyExpansion
}

// workspaceAccessKeys implements WorkspaceAccessKeyInterface
type workspaceAccessKeys struct {
	*gentype.ClientWithList[*managementv1.WorkspaceAccessKey, *managementv1.WorkspaceAccessKeyList]
}

// newWorkspaceAccessKeys returns a WorkspaceAccessKeys
func newWorkspaceAccessKeys(c *ManagementV1Client) *workspaceAccessKeys {
	return &workspaceAccessKeys{
		gentype.NewClientWithList[*managementv1.WorkspaceAccessKey, *managementv1.WorkspaceAccessKeyList](
			"workspaceaccesskeys",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *managementv1.WorkspaceAccessKey { return &managementv1.WorkspaceAccessKey{} },
			func() *managementv1.WorkspaceAccessKeyList { return &managementv1.WorkspaceAccessKeyList{} },
		),
	}
}
