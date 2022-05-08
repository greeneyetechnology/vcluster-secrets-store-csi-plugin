package syncers

import (
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/controller-runtime/pkg/log"
	secretsstorev1 "sigs.k8s.io/secrets-store-csi-driver/apis/v1"
	"github.com/loft-sh/vcluster-sdk/plugin"
	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	"github.com/loft-sh/vcluster-sdk/translate"
	"k8s.io/apimachinery/pkg/api/equality"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	// Make sure our scheme is registered
	_ = secretsstorev1.AddToScheme(plugin.Scheme)
}

func NewSecretStoreSyncer(ctx *synccontext.RegisterContext) syncer.Base {
	return &secretStoreSyncer{
		NamespacedTranslator: translator.NewNamespacedTranslator(ctx, "spc", &secretsstorev1.SecretProviderClass{}),
	}
}

type secretStoreSyncer struct {
	translator.NamespacedTranslator
}

var _ syncer.Initializer = &secretStoreSyncer{}

func (s *secretStoreSyncer) Init(ctx *synccontext.RegisterContext) error {
	return translate.EnsureCRDFromPhysicalCluster(ctx.Context, ctx.PhysicalManager.GetConfig(), ctx.VirtualManager.GetConfig(), secretsstorev1.SchemeGroupVersion.WithKind("SecretProviderClass"))
}

func (s *secretStoreSyncer) SyncDown(ctx *synccontext.SyncContext, vObj client.Object) (ctrl.Result, error) {
	return s.SyncDownCreate(ctx, vObj, s.translate(vObj.(*secretsstorev1.SecretProviderClass)))
}

func (s *secretStoreSyncer) translate(vObj *secretsstorev1.SecretProviderClass) *secretsstorev1.SecretProviderClass {
	newObj := vObj.DeepCopy()
	newObj.Spec.SecretObjects = s.translateSecretsToPhysical(vObj)
	return s.TranslateMetadata(newObj).(*secretsstorev1.SecretProviderClass)
}

func (s *secretStoreSyncer) translateSecretsToPhysical(vObj *secretsstorev1.SecretProviderClass) []*secretsstorev1.SecretObject {
	var sos []*secretsstorev1.SecretObject
	for _, secretObject := range vObj.Spec.SecretObjects {
		so := secretObject.DeepCopy()
		so.SecretName = translate.PhysicalName(so.SecretName, vObj.Namespace)
		sos = append(sos, so)
	}
	return sos
}

func (s *secretStoreSyncer) Sync(ctx *synccontext.SyncContext, pObj client.Object, vObj client.Object) (ctrl.Result, error) {
	return s.SyncDownUpdate(ctx, vObj, s.translateUpdate(pObj.(*secretsstorev1.SecretProviderClass), vObj.(*secretsstorev1.SecretProviderClass)))
}

func (s *secretStoreSyncer) translateUpdate(pObj, vObj *secretsstorev1.SecretProviderClass) *secretsstorev1.SecretProviderClass {
	var updated *secretsstorev1.SecretProviderClass

	// check annotations & labels
	changed, updatedAnnotations, updatedLabels := s.TranslateMetadataUpdate(vObj, pObj)
	if changed {
		updated = newIfNil(updated, pObj)
		updated.Labels = updatedLabels
		updated.Annotations = updatedAnnotations
	}

	newObj := vObj.DeepCopy()
	newObj.Spec.SecretObjects = s.translateSecretsToPhysical(vObj)

	// check spec
	if !equality.Semantic.DeepEqual(newObj.Spec, pObj.Spec) {
		updated = newIfNil(updated, pObj)
		updated.Spec = newObj.Spec
	}

	return updated
}

func newIfNil(updated *secretsstorev1.SecretProviderClass, pObj *secretsstorev1.SecretProviderClass) *secretsstorev1.SecretProviderClass {
	if updated == nil {
		return pObj.DeepCopy()
	}
	return updated
}
