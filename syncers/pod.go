package syncers

import (
	"context"
	"fmt"
	"github.com/loft-sh/vcluster-sdk/hook"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func NewPodHook() hook.ClientHook {
	return &podHook{}
}

type podHook struct{}

func (p *podHook) Name() string {
	return "pod-hook"
}

func (p *podHook) Resource() client.Object {
	return &corev1.Pod{}
}

var _ hook.MutateCreatePhysical = &podHook{}

func (p *podHook) MutateCreatePhysical(ctx context.Context, obj client.Object) (client.Object, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object %v is not a pod", obj)
	}

	log.Log.Info("pod created", "name", pod.Name)
	volumes := pod.Spec.Volumes
	for i, volume := range volumes {
		if volume.CSI == nil {
			continue
		}
		if volume.CSI.Driver == "secrets-store.csi.k8s.io" {
			vclusterName := pod.Labels["vcluster.loft.sh/managed-by"]
			vclusterNs := pod.Annotations["vcluster.loft.sh/namespace"]
			secretProviderClass := volume.CSI.VolumeAttributes["secretProviderClass"]
			pod.Spec.Volumes[i].CSI.VolumeAttributes["secretProviderClass"] = fmt.Sprintf(
				"%s-x-%s-x-%s",
				secretProviderClass,
				vclusterNs,
				vclusterName,
			)
		}
	}

	return pod, nil
}
