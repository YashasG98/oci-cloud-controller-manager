package oci

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/oracle/oci-go-sdk/core"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	clientretry "k8s.io/client-go/util/retry"
)

var UpdateNodeSpecBackoff = wait.Backoff{
	Steps:    20,
	Duration: 50 * time.Millisecond,
	Jitter:   1.0,
}

type FaultDomainController struct {
	nodeInformer coreinformers.NodeInformer
	kubeClient   clientset.Interface
	recorder     record.EventRecorder

	cloud *CloudProvider
}

// NewFaultDomainController creates a FaultDomainController object
func NewFaultDomainController(
	nodeInformer coreinformers.NodeInformer,
	kubeClient clientset.Interface,
	cloud *CloudProvider) *FaultDomainController {

	eventBroadcaster := record.NewBroadcaster()
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "fault-domain-controller"})
	eventBroadcaster.StartLogging(glog.Infof)
	if kubeClient != nil {
		glog.V(0).Infof("Sending events to api server.")
		eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	} else {
		glog.V(0).Infof("No api server defined - no events will be sent to API server.")
	}

	fdc := &FaultDomainController{
		nodeInformer: nodeInformer,
		kubeClient:   kubeClient,
		recorder:     recorder,
		cloud:        cloud,
	}

	// Use shared informer to listen to add nodes
	fdc.nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: fdc.AddCloudNode,
	})

	return fdc
}

func (cnc *FaultDomainController) Run(stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
}

// This adds the fault domain label to the node added
func (cnc *FaultDomainController) AddCloudNode(obj interface{}) {
	node := obj.(*v1.Node)

	err := clientretry.RetryOnConflict(UpdateNodeSpecBackoff, func() error {
		curNode, err := cnc.kubeClient.CoreV1().Nodes().Get(node.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		var instanceID string
		var instance *core.Instance
		instanceID, err = cnc.cloud.InstanceID(context.TODO(), types.NodeName(curNode.Name))
		if err != nil {
			glog.Errorf("Failed to map provider ID to instance ID", err)
		}
		instance, err = cnc.cloud.client.Compute().GetInstance(context.TODO(), instanceID)
		if err != nil {
			glog.Errorf("Failed to get instance from instance ID", err)
		}
		glog.V(2).Infof("Adding node label from cloud provider: %s=%s", "oke.oraclecloud.com/fault-domain", *instance.FaultDomain)
		curNode.ObjectMeta.Labels["oke.oraclecloud.com/fault-domain"] = *instance.FaultDomain

		_, err = cnc.kubeClient.CoreV1().Nodes().Update(curNode)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		utilruntime.HandleError(err)
		return
	}

	glog.Infof("Fault domain controller has successfully added label to node %s", node.Name)
}
