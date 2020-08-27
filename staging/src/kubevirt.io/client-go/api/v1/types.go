/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2017 Red Hat, Inc.
 *
 */

package v1

//go:generate swagger-doc
//go:generate deepcopy-gen -i . --go-header-file ../../../../../../hack/boilerplate/boilerplate.go.txt
//go:generate defaulter-gen -i . --go-header-file ../../../../../../hack/boilerplate/boilerplate.go.txt
//go:generate openapi-gen -i kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1,k8s.io/apimachinery/pkg/util/intstr,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/runtime,k8s.io/api/core/v1,kubevirt.io/client-go/api/v1 --output-package=kubevirt.io/kubevirt/staging/src/kubevirt.io/client-go/api/v1  --go-header-file ../../../../../../hack/boilerplate/boilerplate.go.txt

/*
 ATTENTION: Rerun code generators when comments on structs or fields are modified.
*/

import (
	"encoding/json"
	"fmt"

	k8sv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"

	cdiv1 "kubevirt.io/containerized-data-importer/pkg/apis/core/v1alpha1"
)

const DefaultGracePeriodSeconds int64 = 30

// VirtualMachineInstance is *the* VirtualMachineInstance Definition. It represents a virtual machine in the runtime environment of kubernetes.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// VirtualMachineInstance Spec contains the VirtualMachineInstance specification.
	Spec VirtualMachineInstanceSpec `json:"spec" valid:"required"`
	// Status is the high level overview of how the VirtualMachineInstance is doing. It contains information available to controllers and users.
	Status VirtualMachineInstanceStatus `json:"status,omitempty"`
}

func (v *VirtualMachineInstance) MarshalBinary() (data []byte, err error) {
	return json.Marshal(*v)
}

func (v *VirtualMachineInstance) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, v)
}

// VirtualMachineInstanceList is a list of VirtualMachines
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineInstance `json:"items"`
}

// +k8s:openapi-gen=true
type EvictionStrategy string

// VirtualMachineInstanceSpec is a description of a VirtualMachineInstance.
//
// +k8s:openapi-gen=true
type VirtualMachineInstanceSpec struct {

	// If specified, indicates the pod's priority.
	// If not specified, the pod priority will be default or zero if there is no
	// default.
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`

	// Specification of the desired behavior of the VirtualMachineInstance on the host.
	Domain DomainSpec `json:"domain"`
	// NodeSelector is a selector which must be true for the vmi to fit on a node.
	// Selector which must match a node's labels for the vmi to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If affinity is specifies, obey all the affinity rules
	Affinity *k8sv1.Affinity `json:"affinity,omitempty"`
	// If specified, the VMI will be dispatched by specified scheduler.
	// If not specified, the VMI will be dispatched by default scheduler.
	// +optional
	SchedulerName string `json:"schedulerName,omitempty"`
	// If toleration is specified, obey all the toleration rules.
	Tolerations []k8sv1.Toleration `json:"tolerations,omitempty"`

	// EvictionStrategy can be set to "LiveMigrate" if the VirtualMachineInstance should be
	// migrated instead of shut-off in case of a node drain.
	//
	// +optional
	EvictionStrategy *EvictionStrategy `json:"evictionStrategy,omitempty"`

	// Grace period observed after signalling a VirtualMachineInstance to stop after which the VirtualMachineInstance is force terminated.
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// List of volumes that can be mounted by disks belonging to the vmi.
	Volumes []Volume `json:"volumes,omitempty"`
	// Periodic probe of VirtualMachineInstance liveness.
	// VirtualmachineInstances will be stopped if the probe fails.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	LivenessProbe *Probe `json:"livenessProbe,omitempty"`
	// Periodic probe of VirtualMachineInstance service readiness.
	// VirtualmachineInstances will be removed from service endpoints if the probe fails.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`
	// Specifies the hostname of the vmi
	// If not specified, the hostname will be set to the name of the vmi, if dhcp or cloud-init is configured properly.
	// +optional
	Hostname string `json:"hostname,omitempty"`
	// If specified, the fully qualified vmi hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>".
	// If not specified, the vmi will not have a domainname at all. The DNS entry will resolve to the vmi,
	// no matter if the vmi itself can pick up a hostname.
	// +optional
	Subdomain string `json:"subdomain,omitempty"`
	// List of networks that can be attached to a vm's virtual interface.
	Networks []Network `json:"networks,omitempty"`
	// Set DNS policy for the pod.
	// Defaults to "ClusterFirst".
	// Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.
	// DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy.
	// To have DNS options set along with hostNetwork, you have to specify DNS policy
	// explicitly to 'ClusterFirstWithHostNet'.
	// +optional
	DNSPolicy k8sv1.DNSPolicy `json:"dnsPolicy,omitempty" protobuf:"bytes,6,opt,name=dnsPolicy,casttype=DNSPolicy"`
	// Specifies the DNS parameters of a pod.
	// Parameters specified here will be merged to the generated DNS
	// configuration based on DNSPolicy.
	// +optional
	DNSConfig *k8sv1.PodDNSConfig `json:"dnsConfig,omitempty" protobuf:"bytes,26,opt,name=dnsConfig"`

	// Overwrite any system migration configuration
	// +optional
	MigrationConfiguration *PerVMIMigrationConfiguration `json:"migrationConfiguration,omitempty"`
}

// VirtualMachineInstanceStatus represents information about the status of a VirtualMachineInstance. Status may trail the actual
// state of a system.
//
// +k8s:openapi-gen=true
type VirtualMachineInstanceStatus struct {
	// NodeName is the name where the VirtualMachineInstance is currently running.
	NodeName string `json:"nodeName,omitempty"`
	// A brief CamelCase message indicating details about why the VMI is in this state. e.g. 'NodeUnresponsive'
	// +optional
	Reason string `json:"reason,omitempty"`
	// Conditions are specific points in VirtualMachineInstance's pod runtime.
	Conditions []VirtualMachineInstanceCondition `json:"conditions,omitempty"`
	// Phase is the status of the VirtualMachineInstance in kubernetes world. It is not the VirtualMachineInstance status, but partially correlates to it.
	Phase VirtualMachineInstancePhase `json:"phase,omitempty"`
	// Interfaces represent the details of available network interfaces.
	Interfaces []VirtualMachineInstanceNetworkInterface `json:"interfaces,omitempty"`
	// Guest OS Information
	GuestOSInfo VirtualMachineInstanceGuestOSInfo `json:"guestOSInfo,omitempty"`
	// Represents the status of a live migration
	MigrationState *VirtualMachineInstanceMigrationState `json:"migrationState,omitempty"`
	// Represents the method using which the vmi can be migrated: live migration or block migration
	MigrationMethod VirtualMachineInstanceMigrationMethod `json:"migrationMethod,omitempty"`
	// The Quality of Service (QOS) classification assigned to the virtual machine instance based on resource requirements
	// See PodQOSClass type for available QOS classes
	// More info: https://git.k8s.io/community/contributors/design-proposals/node/resource-qos.md
	// +optional
	QOSClass *k8sv1.PodQOSClass `json:"qosClass,omitempty"`

	// EvacuationNodeName is used to track the eviction process of a VMI. It stores the name of the node that we want
	// to evacuate. It is meant to be used by KubeVirt core components only and can't be set or modified by users.
	// +optional
	EvacuationNodeName string `json:"evacuationNodeName,omitempty"`

	// ActivePods is a mapping of pod UID to node name.
	// It is possible for multiple pods to be running for a single VMI during migration.
	ActivePods map[types.UID]string `json:"activePods,omitempty"`
}

func (v *VirtualMachineInstance) IsScheduling() bool {
	return v.Status.Phase == Scheduling
}

func (v *VirtualMachineInstance) IsScheduled() bool {
	return v.Status.Phase == Scheduled
}

func (v *VirtualMachineInstance) IsRunning() bool {
	return v.Status.Phase == Running
}

func (v *VirtualMachineInstance) IsMarkedForEviction() bool {
	return v.Status.EvacuationNodeName != ""
}

func (v *VirtualMachineInstance) IsMigratable() bool {
	for _, cond := range v.Status.Conditions {
		if cond.Type == VirtualMachineInstanceIsMigratable && cond.Status == k8sv1.ConditionTrue {
			return true
		}
	}
	return false
}

func (v *VirtualMachineInstance) IsEvictable() bool {
	return v.Spec.EvictionStrategy != nil && *v.Spec.EvictionStrategy == EvictionStrategyLiveMigrate
}

func (v *VirtualMachineInstance) IsFinal() bool {
	return v.Status.Phase == Failed || v.Status.Phase == Succeeded
}

func (v *VirtualMachineInstance) IsUnknown() bool {
	return v.Status.Phase == Unknown
}

func (v *VirtualMachineInstance) IsUnprocessed() bool {
	return v.Status.Phase == Pending || v.Status.Phase == VmPhaseUnset
}

// Checks if CPU pinning has been requested
func (v *VirtualMachineInstance) IsCPUDedicated() bool {
	return v.Spec.Domain.CPU != nil && v.Spec.Domain.CPU.DedicatedCPUPlacement
}

// WantsToHaveQOSGuaranteed checks if cpu and memoyr limits and requests are identical on the VMI.
// This is the indicator that people want a VMI with QOS of guaranteed
func (v *VirtualMachineInstance) WantsToHaveQOSGuaranteed() bool {
	resources := v.Spec.Domain.Resources
	return !resources.Requests.Memory().IsZero() && resources.Requests.Memory().Cmp(*resources.Limits.Memory()) == 0 &&
		!resources.Requests.Cpu().IsZero() && resources.Requests.Cpu().Cmp(*resources.Limits.Cpu()) == 0
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceConditionType string

// These are valid conditions of VMIs.
const (
	// VMIReady means the pod is able to service requests and should be added to the
	// load balancing pools of all matching services.
	VirtualMachineInstanceReady VirtualMachineInstanceConditionType = "Ready"

	// If there happens any error while trying to synchronize the VirtualMachineInstance with the Domain,
	// this is reported as false.
	VirtualMachineInstanceSynchronized VirtualMachineInstanceConditionType = "Synchronized"

	// If the VMI was paused by the user, this is reported as true.
	VirtualMachineInstancePaused VirtualMachineInstanceConditionType = "Paused"

	// Reflects whether the QEMU guest agent is connected through the channel
	VirtualMachineInstanceAgentConnected VirtualMachineInstanceConditionType = "AgentConnected"

	// Reflects whether the QEMU guest agent is connected through the channel
	VirtualMachineInstanceUnsupportedAgent VirtualMachineInstanceConditionType = "AgentVersionNotSupported"

	// Indicates whether the VMI is live migratable
	VirtualMachineInstanceIsMigratable VirtualMachineInstanceConditionType = "LiveMigratable"
	// Reason means that VMI is not live migratioable because of it's disks collection
	VirtualMachineInstanceReasonDisksNotMigratable = "DisksNotLiveMigratable"
	// Reason means that VMI is not live migratioable because of it's network interfaces collection
	VirtualMachineInstanceReasonInterfaceNotMigratable = "InterfaceNotLiveMigratable"
)

// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationConditionType string

// These are valid conditions of VMIs.
const (
	// VirtualMachineInstanceMigrationAbortRequested indicates that live migration abort has been requested
	VirtualMachineInstanceMigrationAbortRequested VirtualMachineInstanceMigrationConditionType = "migrationAbortRequested"
)

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceCondition struct {
	Type   VirtualMachineInstanceConditionType `json:"type"`
	Status k8sv1.ConditionStatus               `json:"status"`
	// +nullable
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationCondition struct {
	Type   VirtualMachineInstanceMigrationConditionType `json:"type"`
	Status k8sv1.ConditionStatus                        `json:"status"`
	// +nullable
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

// The migration phase indicates that the job has completed
func (m *VirtualMachineInstanceMigration) IsFinal() bool {
	return m.Status.Phase == MigrationFailed || m.Status.Phase == MigrationSucceeded
}

func (m *VirtualMachineInstanceMigration) IsRunning() bool {
	switch m.Status.Phase {
	case MigrationFailed, MigrationPending, MigrationPhaseUnset, MigrationSucceeded:
		return false
	}
	return true
}

// The migration phase indicates that the target pod should have already been created
func (m *VirtualMachineInstanceMigration) TargetIsCreated() bool {
	return m.Status.Phase != MigrationPhaseUnset &&
		m.Status.Phase != MigrationPending
}

// The migration phase indicates that job has been handed off to the VMI controllers to complete.
func (m *VirtualMachineInstanceMigration) TargetIsHandedOff() bool {
	return m.Status.Phase != MigrationPhaseUnset &&
		m.Status.Phase != MigrationPending &&
		m.Status.Phase != MigrationScheduling &&
		m.Status.Phase != MigrationScheduled
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceNetworkInterface struct {
	// IP address of a Virtual Machine interface. It is always the first item of
	// IPs
	IP string `json:"ipAddress,omitempty"`
	// Hardware address of a Virtual Machine interface
	MAC string `json:"mac,omitempty"`
	// Name of the interface, corresponds to name of the network assigned to the interface
	// TODO: remove omitempty, when api breaking changes are allowed
	Name string `json:"name,omitempty"`
	// List of all IP addresses of a Virtual Machine interface
	IPs []string `json:"ipAddresses,omitempty"`
	// The interface name inside the Virtual Machine
	InterfaceName string `json:"interfaceName,omitempty"`
}

// +k8s:openapi-gen=true
type VirtualMachineInstanceGuestOSInfo struct {
	// Name of the Guest OS
	Name string `json:"name,omitempty"`
	// Guest OS Kernel Release
	KernelRelease string `json:"kernelRelease,omitempty"`
	// Guest OS Version
	Version string `json:"version,omitempty"`
	// Guest OS Pretty Name
	PrettyName string `json:"prettyName,omitempty"`
	// Version ID of the Guest OS
	VersionID string `json:"versionId,omitempty"`
	// Kernel version of the Guest OS
	KernelVersion string `json:"kernelVersion,omitempty"`
	// Machine type of the Guest OS
	Machine string `json:"machine,omitempty"`
	// Guest OS Id
	ID string `json:"id,omitempty"`
}

// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationState struct {
	// The time the migration action began
	StartTimestamp *metav1.Time `json:"startTimestamp,omitempty"`
	// The time the migration action ended
	EndTimestamp *metav1.Time `json:"endTimestamp,omitempty"`
	// The Target Node has seen the Domain Start Event
	TargetNodeDomainDetected bool `json:"targetNodeDomainDetected,omitempty"`
	// The address of the target node to use for the migration
	TargetNodeAddress string `json:"targetNodeAddress,omitempty"`
	// The list of ports opened for live migration on the destination node
	TargetDirectMigrationNodePorts map[string]int `json:"targetDirectMigrationNodePorts,omitempty"`
	// The target node that the VMI is moving to
	TargetNode string `json:"targetNode,omitempty"`
	// The target pod that the VMI is moving to
	TargetPod string `json:"targetPod,omitempty"`
	// The source node that the VMI originated on
	SourceNode string `json:"sourceNode,omitempty"`
	// Indicates the migration completed
	Completed bool `json:"completed,omitempty"`
	// Indicates that the migration failed
	Failed bool `json:"failed,omitempty"`
	// Indicates that the migration has been requested to abort
	AbortRequested bool `json:"abortRequested,omitempty"`
	// Indicates the final status of the live migration abortion
	AbortStatus MigrationAbortStatus `json:"abortStatus,omitempty"`
	// The VirtualMachineInstanceMigration object associated with this migration
	MigrationUID types.UID `json:"migrationUid,omitempty"`
	// Lets us know if the vmi is currenly running pre or post copy migration
	Mode MigrationMode `json:"mode,omitempty"`
}

//
// +k8s:openapi-gen=true
type MigrationAbortStatus string

const (
	// MigrationAbortSucceeded means that the VirtualMachineInstance live migration has been aborted
	MigrationAbortSucceeded MigrationAbortStatus = "Succeeded"
	// MigrationAbortFailed means that the vmi live migration has failed to be abort
	MigrationAbortFailed MigrationAbortStatus = "Failed"
	// MigrationAbortInProgress mean that the vmi live migration is aborting
	MigrationAbortInProgress MigrationAbortStatus = "Aborting"
)

//
// +k8s:openapi-gen=true
type MigrationMode string

const (
	// MigrationPreCopy means the VMI migrations that is currenly running is in pre copy mode
	MigrationPreCopy MigrationMode = "PreCopy"
	// MigrationPostCopy means the VMI migrations that is currenly running is in post copy mode
	MigrationPostCopy MigrationMode = "PostCopy"
)

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationMethod string

const (
	// BlockMigration means that all VirtualMachineInstance disks should be copied over to the destination host
	BlockMigration VirtualMachineInstanceMigrationMethod = "BlockMigration"
	// LiveMigration means that VirtualMachineInstance disks will not be copied over to the destination host
	LiveMigration VirtualMachineInstanceMigrationMethod = "LiveMigration"
)

// VirtualMachineInstancePhase is a label for the condition of a VirtualMachineInstance at the current time.
//
// +k8s:openapi-gen=true
type VirtualMachineInstancePhase string

// These are the valid statuses of pods.
const (
	//When a VirtualMachineInstance Object is first initialized and no phase, or Pending is present.
	VmPhaseUnset VirtualMachineInstancePhase = ""
	// Pending means the VirtualMachineInstance has been accepted by the system.
	Pending VirtualMachineInstancePhase = "Pending"
	// A target Pod exists but is not yet scheduled and in running state.
	Scheduling VirtualMachineInstancePhase = "Scheduling"
	// A target pod was scheduled and the system saw that Pod in runnig state.
	// Here is where the responsibility of virt-controller ends and virt-handler takes over.
	Scheduled VirtualMachineInstancePhase = "Scheduled"
	// Running means the pod has been bound to a node and the VirtualMachineInstance is started.
	Running VirtualMachineInstancePhase = "Running"
	// Succeeded means that the VirtualMachineInstance stopped voluntarily, e.g. reacted to SIGTERM or shutdown was invoked from
	// inside the VirtualMachineInstance.
	Succeeded VirtualMachineInstancePhase = "Succeeded"
	// Failed means that the vmi crashed, disappeared unexpectedly or got deleted from the cluster before it was ever started.
	Failed VirtualMachineInstancePhase = "Failed"
	// Unknown means that for some reason the state of the VirtualMachineInstance could not be obtained, typically due
	// to an error in communicating with the host of the VirtualMachineInstance.
	Unknown VirtualMachineInstancePhase = "Unknown"
)

const (
	// This label marks resources that belong to KubeVirt. An optional value
	// may indicate which specific KubeVirt component a resource belongs to.
	AppLabel string = "kubevirt.io"
	// This annotation is used to match virtual machine instances represented as
	// libvirt XML domains with their pods. Among other things, the annotation is
	// used to detect virtual machines with dead pods. Used on Pod.
	DomainAnnotation string = "kubevirt.io/domain"
	// Represents the name of the migration job this target pod is associated with
	MigrationJobNameAnnotation                    string = "kubevirt.io/migrationJobName"
	ControllerAPILatestVersionObservedAnnotation  string = "kubevirt.io/latest-observed-api-version"
	ControllerAPIStorageVersionObservedAnnotation string = "kubevirt.io/storage-observed-api-version"
	// This label is used to match virtual machine instance IDs with pods.
	// Similar to kubevirt.io/domain. Used on Pod.
	// Internal use only.
	CreatedByLabel string = "kubevirt.io/created-by"
	// This label is used to indicate that this pod is the target of a migration job.
	MigrationJobLabel string = "kubevirt.io/migrationJobUID"
	// This label describes which cluster node runs the virtual machine
	// instance. Needed because with CRDs we can't use field selectors. Used on
	// VirtualMachineInstance.
	NodeNameLabel string = "kubevirt.io/nodeName"
	// This label describes which cluster node runs the target Pod for a Virtual
	// Machine Instance migration job. Needed because with CRDs we can't use field
	// selectors. Used on VirtualMachineInstance.
	MigrationTargetNodeNameLabel string = "kubevirt.io/migrationTargetNodeName"
	// This label declares whether a particular node is available for
	// scheduling virtual machine instances on it. Used on Node.
	NodeSchedulable string = "kubevirt.io/schedulable"
	// This annotation is regularly updated by virt-handler to help determine
	// if a particular node is alive and hence should be available for new
	// virtual machine instance scheduling. Used on Node.
	VirtHandlerHeartbeat string = "kubevirt.io/heartbeat"
	// Namespace recommended by Kubernetes for commonly recognized labels
	AppLabelPrefix = "app.kubernetes.io"
	// This label is commonly used by 3rd party management tools to identify
	// an application's name.
	AppNameLabel = AppLabelPrefix + "/name"
	// This label is commonly used by 3rd party management tools to identify
	// an application's version.
	AppVersionLabel = AppLabelPrefix + "/version"
	// This label is commonly used by 3rd party management tools to identify
	// a higher level application.
	AppPartOfLabel = AppLabelPrefix + "/part-of"
	// This label is commonly used by 3rd party management tools to identify
	// the component this application is a part of.
	AppComponentLabel = AppLabelPrefix + "/component"
	// This label identifies each resource as part of KubeVirt
	AppComponent = "kubevirt"
	// This label will be set on all resources created by the operator
	ManagedByLabel              = AppLabelPrefix + "/managed-by"
	ManagedByLabelOperatorValue = "kubevirt-operator"
	// This annotation represents the kubevirt version for an install strategy configmap.
	InstallStrategyVersionAnnotation = "kubevirt.io/install-strategy-version"
	// This annotation represents the kubevirt registry used for an install strategy configmap.
	InstallStrategyRegistryAnnotation = "kubevirt.io/install-strategy-registry"
	// This annotation represents the kubevirt deployment identifier used for an install strategy configmap.
	InstallStrategyIdentifierAnnotation = "kubevirt.io/install-strategy-identifier"
	// This annotation is a hash of all customizations that live under spec.CustomizeComponents
	KubeVirtCustomizeComponentAnnotationHash = "kubevirt.io/customizer-identifier"
	// This annotation represents the kubevirt generation that was used to create a resource
	KubeVirtGenerationAnnotation = "kubevirt.io/generation"
	// This annotation represents that this object is for temporary use during updates
	EphemeralBackupObject = "kubevirt.io/ephemeral-backup-object"

	// This label indicates the object is a part of the install strategy retrieval process.
	InstallStrategyLabel = "kubevirt.io/install-strategy"

	VirtualMachineInstanceFinalizer          string = "foregroundDeleteVirtualMachine"
	VirtualMachineInstanceMigrationFinalizer string = "kubevirt.io/migrationJobFinalize"
	CPUManager                               string = "cpumanager"
	// This annotation is used to inject ignition data
	// Used on VirtualMachineInstance.
	IgnitionAnnotation           string = "kubevirt.io/ignitiondata"
	PlacePCIDevicesOnRootComplex string = "kubevirt.io/placePCIDevicesOnRootComplex"

	VirtualMachineLabel        = AppLabel + "/vm"
	MemfdMemoryBackend  string = "kubevirt.io/memfd"
)

func NewVMI(name string, uid types.UID) *VirtualMachineInstance {
	return &VirtualMachineInstance{
		Spec: VirtualMachineInstanceSpec{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			UID:       uid,
			Namespace: k8sv1.NamespaceDefault,
		},
		Status: VirtualMachineInstanceStatus{},
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       VirtualMachineInstanceGroupVersionKind.Kind,
		},
	}
}

//
// +k8s:openapi-gen=true
type SyncEvent string

const (
	Created         SyncEvent = "Created"
	Deleted         SyncEvent = "Deleted"
	PresetFailed    SyncEvent = "PresetFailed"
	Override        SyncEvent = "Override"
	Started         SyncEvent = "Started"
	ShuttingDown    SyncEvent = "ShuttingDown"
	Stopped         SyncEvent = "Stopped"
	PreparingTarget SyncEvent = "PreparingTarget"
	Migrating       SyncEvent = "Migrating"
	Migrated        SyncEvent = "Migrated"
	SyncFailed      SyncEvent = "SyncFailed"
	Resumed         SyncEvent = "Resumed"
)

func (s SyncEvent) String() string {
	return string(s)
}

// TODO Namespace could be different, also store it somewhere in the domain, so that we can report deletes on handler startup properly
func NewVMIReferenceFromName(name string) *VirtualMachineInstance {
	return NewVMIReferenceFromNameWithNS(k8sv1.NamespaceDefault, name)
}

func NewVMIReferenceFromNameWithNS(namespace string, name string) *VirtualMachineInstance {
	vmi := &VirtualMachineInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			SelfLink:  fmt.Sprintf("/apis/%s/namespaces/%s/virtualmachineinstances/%s", GroupVersion.String(), namespace, name),
		},
	}
	vmi.SetGroupVersionKind(schema.GroupVersionKind{Group: GroupVersion.Group, Kind: "VirtualMachineInstance", Version: GroupVersion.Version})
	return vmi
}

func NewVMIReferenceWithUUID(namespace string, name string, uuid types.UID) *VirtualMachineInstance {
	vmi := NewVMIReferenceFromNameWithNS(namespace, name)
	vmi.UID = uuid
	return vmi
}

type VMISelector struct {
	// Name of the VirtualMachineInstance to migrate
	Name string `json:"name" valid:"required"`
}

// Given a VirtualMachineInstance, update all NodeSelectorTerms with anti-affinity for that VirtualMachineInstance's node.
// This is useful for the case when a migration away from a node must occur.
// This method returns the full Affinity structure updated the anti affinity terms
func UpdateAntiAffinityFromVMINode(pod *k8sv1.Pod, vmi *VirtualMachineInstance) *k8sv1.Affinity {
	if pod.Spec.Affinity == nil {
		pod.Spec.Affinity = &k8sv1.Affinity{}
	}

	if pod.Spec.Affinity.NodeAffinity == nil {
		pod.Spec.Affinity.NodeAffinity = &k8sv1.NodeAffinity{}
	}

	if pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &k8sv1.NodeSelector{}
	}

	selector := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	terms := selector.NodeSelectorTerms

	if len(terms) == 0 {
		selector.NodeSelectorTerms = append(terms, k8sv1.NodeSelectorTerm{})
		terms = selector.NodeSelectorTerms
	}

	for idx, term := range terms {
		if term.MatchExpressions == nil {
			term.MatchExpressions = []k8sv1.NodeSelectorRequirement{}
		}

		term.MatchExpressions = append(term.MatchExpressions, PrepareVMINodeAntiAffinitySelectorRequirement(vmi))
		selector.NodeSelectorTerms[idx] = term
	}

	return pod.Spec.Affinity
}

// Given a VirtualMachineInstance, create a NodeSelectorTerm with anti-affinity for that VirtualMachineInstance's node.
// This is useful for the case when a migration away from a node must occur.
func PrepareVMINodeAntiAffinitySelectorRequirement(vmi *VirtualMachineInstance) k8sv1.NodeSelectorRequirement {
	return k8sv1.NodeSelectorRequirement{
		Key:      "kubernetes.io/hostname",
		Operator: k8sv1.NodeSelectorOpNotIn,
		Values:   []string{vmi.Status.NodeName},
	}
}

// VirtualMachineInstance is *the* VirtualMachineInstance Definition. It represents a virtual machine in the runtime environment of kubernetes.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// VirtualMachineInstance Spec contains the VirtualMachineInstance specification.
	Spec VirtualMachineInstanceReplicaSetSpec `json:"spec" valid:"required"`
	// Status is the high level overview of how the VirtualMachineInstance is doing. It contains information available to controllers and users.
	// +nullable
	Status VirtualMachineInstanceReplicaSetStatus `json:"status,omitempty"`
}

// VMIList is a list of VMIs
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineInstanceReplicaSet `json:"items"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceReplicaSetSpec struct {
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this deployment.
	Selector *metav1.LabelSelector `json:"selector" valid:"required"`

	// Template describes the pods that will be created.
	Template *VirtualMachineInstanceTemplateSpec `json:"template" valid:"required"`

	// Indicates that the replica set is paused.
	// +optional
	Paused bool `json:"paused,omitempty" protobuf:"varint,7,opt,name=paused"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceReplicaSetStatus struct {
	// Total number of non-terminated pods targeted by this deployment (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`

	// The number of ready replicas for this replica set.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,4,opt,name=readyReplicas"`

	Conditions []VirtualMachineInstanceReplicaSetCondition `json:"conditions,omitempty" optional:"true"`

	// Canonical form of the label selector for HPA which consumes it through the scale subresource.
	LabelSelector string `json:"labelSelector,omitempty"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceReplicaSetCondition struct {
	Type   VirtualMachineInstanceReplicaSetConditionType `json:"type"`
	Status k8sv1.ConditionStatus                         `json:"status"`
	// +nullable
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceReplicaSetConditionType string

const (
	// VirtualMachineInstanceReplicaSetReplicaFailure is added in a replica set when one of its vmis
	// fails to be created due to insufficient quota, limit ranges, pod security policy, node selectors,
	// etc. or deleted due to kubelet being down or finalizers are failing.
	VirtualMachineInstanceReplicaSetReplicaFailure VirtualMachineInstanceReplicaSetConditionType = "ReplicaFailure"

	// VirtualMachineInstanceReplicaSetReplicaPaused is added in a replica set when the replica set got paused by the controller.
	// After this condition was added, it is safe to remove or add vmis by hand and adjust the replica count by hand.
	VirtualMachineInstanceReplicaSetReplicaPaused VirtualMachineInstanceReplicaSetConditionType = "ReplicaPaused"
)

//
// +k8s:openapi-gen=true
type DataVolumeTemplateSpec struct {
	// +nullable
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// DataVolumeSpec contains the DataVolume specification.
	Spec cdiv1.DataVolumeSpec `json:"spec"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceTemplateSpec struct {
	// +nullable
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`
	// VirtualMachineInstance Spec contains the VirtualMachineInstance specification.
	Spec VirtualMachineInstanceSpec `json:"spec,omitempty" valid:"required"`
}

// VirtualMachineInstanceMigration represents the object tracking a VMI's migration
// to another host in the cluster
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VirtualMachineInstanceMigrationSpec   `json:"spec" valid:"required"`
	Status            VirtualMachineInstanceMigrationStatus `json:"status,omitempty"`
}

// VirtualMachineInstanceMigrationList is a list of VirtualMachineMigrations
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineInstanceMigration `json:"items"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationSpec struct {
	// The name of the VMI to perform the migration on. VMI must exist in the migration objects namespace
	VMIName string `json:"vmiName,omitempty" valid:"required"`
}

// VirtualMachineInstanceMigration reprents information pertaining to a VMI's migration.
//
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationStatus struct {
	Phase      VirtualMachineInstanceMigrationPhase       `json:"phase,omitempty"`
	Conditions []VirtualMachineInstanceMigrationCondition `json:"conditions,omitempty"`
}

// VirtualMachineInstanceMigrationPhase is a label for the condition of a VirtualMachineInstanceMigration at the current time.
//
// +k8s:openapi-gen=true
type VirtualMachineInstanceMigrationPhase string

// These are the valid migration phases
const (
	MigrationPhaseUnset VirtualMachineInstanceMigrationPhase = ""
	// The migration is accepted by the system
	MigrationPending VirtualMachineInstanceMigrationPhase = "Pending"
	// The migration's target pod is being scheduled
	MigrationScheduling VirtualMachineInstanceMigrationPhase = "Scheduling"
	// The migration's target pod is running
	MigrationScheduled VirtualMachineInstanceMigrationPhase = "Scheduled"
	// The migration's target pod is being prepared for migration
	MigrationPreparingTarget VirtualMachineInstanceMigrationPhase = "PreparingTarget"
	// The migration's target pod is prepared and ready for migration
	MigrationTargetReady VirtualMachineInstanceMigrationPhase = "TargetReady"
	// The migration is in progress
	MigrationRunning VirtualMachineInstanceMigrationPhase = "Running"
	// The migration passed
	MigrationSucceeded VirtualMachineInstanceMigrationPhase = "Succeeded"
	// The migration failed
	MigrationFailed VirtualMachineInstanceMigrationPhase = "Failed"
)

//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstancePreset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// VirtualMachineInstance Spec contains the VirtualMachineInstance specification.
	Spec VirtualMachineInstancePresetSpec `json:"spec,omitempty" valid:"required"`
}

// VirtualMachineInstancePresetList is a list of VirtualMachinePresets
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstancePresetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineInstancePreset `json:"items"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineInstancePresetSpec struct {
	// Selector is a label query over a set of VMIs.
	// Required.
	Selector metav1.LabelSelector `json:"selector"`
	// Domain is the same object type as contained in VirtualMachineInstanceSpec
	Domain *DomainSpec `json:"domain,omitempty"`
}

func NewVirtualMachinePreset(name string, selector metav1.LabelSelector) *VirtualMachineInstancePreset {
	return &VirtualMachineInstancePreset{
		Spec: VirtualMachineInstancePresetSpec{
			Selector: selector,
			Domain:   &DomainSpec{},
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: k8sv1.NamespaceDefault,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: GroupVersion.String(),
			Kind:       VirtualMachineInstancePresetGroupVersionKind.Kind,
		},
	}
}

// VirtualMachine handles the VirtualMachines that are not running
// or are in a stopped state
// The VirtualMachine contains the template to create the
// VirtualMachineInstance. It also mirrors the running state of the created
// VirtualMachineInstance in its status.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Spec contains the specification of VirtualMachineInstance created
	Spec VirtualMachineSpec `json:"spec" valid:"required"`
	// Status holds the current state of the controller and brief information
	// about its associated VirtualMachineInstance
	Status VirtualMachineStatus `json:"status,omitempty"`
}

// Return the current runStrategy for the VirtualMachine
// if vm.spec.running is set, that will be mapped to runStrategy:
//   false: RunStrategyHalted
//   true: RunStrategyAlways
func (vm *VirtualMachine) RunStrategy() (VirtualMachineRunStrategy, error) {
	if vm.Spec.Running != nil && vm.Spec.RunStrategy != nil {
		return RunStrategyUnknown, fmt.Errorf("running and runstrategy are mutually exclusive")
	}
	RunStrategy := RunStrategyHalted
	if vm.Spec.Running != nil {
		if (*vm.Spec.Running) == true {
			RunStrategy = RunStrategyAlways
		}
	} else if vm.Spec.RunStrategy != nil {
		RunStrategy = *vm.Spec.RunStrategy
	}
	return RunStrategy, nil
}

// VirtualMachineList is a list of virtualmachines
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

// VirtualMachineRunStrategy is a label for the requested VirtualMachineInstance Running State at the current time.
//
// +k8s:openapi-gen=true
type VirtualMachineRunStrategy string

// These are the valid VMI run strategies
const (
	// Placeholder. Not a valid RunStrategy.
	RunStrategyUnknown VirtualMachineRunStrategy = ""
	// VMI should always be running.
	RunStrategyAlways VirtualMachineRunStrategy = "Always"
	// VMI should never be running.
	RunStrategyHalted VirtualMachineRunStrategy = "Halted"
	// VMI can be started/stopped using API endpoints.
	RunStrategyManual VirtualMachineRunStrategy = "Manual"
	// VMI will initially be running--and restarted if a failure occurs.
	// It will not be restarted upon successful completion.
	RunStrategyRerunOnFailure VirtualMachineRunStrategy = "RerunOnFailure"
)

// VirtualMachineSpec describes how the proper VirtualMachine
// should look like
//
// +k8s:openapi-gen=true
type VirtualMachineSpec struct {
	// Running controls whether the associatied VirtualMachineInstance is created or not
	// Mutually exclusive with RunStrategy
	Running *bool `json:"running,omitempty" optional:"true"`

	// Running state indicates the requested running state of the VirtualMachineInstance
	// mutually exclusive with Running
	RunStrategy *VirtualMachineRunStrategy `json:"runStrategy,omitempty" optional:"true"`

	// Template is the direct specification of VirtualMachineInstance
	Template *VirtualMachineInstanceTemplateSpec `json:"template"`

	// dataVolumeTemplates is a list of dataVolumes that the VirtualMachineInstance template can reference.
	// DataVolumes in this list are dynamically created for the VirtualMachine and are tied to the VirtualMachine's life-cycle.
	DataVolumeTemplates []DataVolumeTemplateSpec `json:"dataVolumeTemplates,omitempty"`
}

// StateChangeRequestType represents the existing state change requests that are possible
//
// +k8s:openapi-gen=true
type StateChangeRequestAction string

// These are the currently defined state change requests
const (
	StartRequest  StateChangeRequestAction = "Start"
	StopRequest   StateChangeRequestAction = "Stop"
	RenameRequest                          = "Rename"
)

// VirtualMachineStatus represents the status returned by the
// controller to describe how the VirtualMachine is doing
//
// +k8s:openapi-gen=true
type VirtualMachineStatus struct {
	// SnapshotInProgress is the name of the VirtualMachineSnapshot currently executing
	SnapshotInProgress *string `json:"snapshotInProgress,omitempty"`
	// Created indicates if the virtual machine is created in the cluster
	Created bool `json:"created,omitempty"`
	// Ready indicates if the virtual machine is running and ready
	Ready bool `json:"ready,omitempty"`
	// Hold the state information of the VirtualMachine and its VirtualMachineInstance
	Conditions []VirtualMachineCondition `json:"conditions,omitempty" optional:"true"`
	// StateChangeRequests indicates a list of actions that should be taken on a VMI
	// e.g. stop a specific VMI then start a new one.
	StateChangeRequests []VirtualMachineStateChangeRequest `json:"stateChangeRequests,omitempty" optional:"true"`
}

// +k8s:openapi-gen=true
type VirtualMachineStateChangeRequest struct {
	// Indicates the type of action that is requested. e.g. Start or Stop
	Action StateChangeRequestAction `json:"action"`
	// Provides additional data in order to perform the Action
	Data map[string]string `json:"data,omitempty" optional:"true"`
	// Indicates the UUID of an existing Virtual Machine Instance that this change request applies to -- if applicable
	UID *types.UID `json:"uid,omitempty" optional:"true" protobuf:"bytes,5,opt,name=uid,casttype=k8s.io/kubernetes/pkg/types.UID"`
}

// VirtualMachineCondition represents the state of VirtualMachine
//
// +k8s:openapi-gen=true
type VirtualMachineCondition struct {
	Type   VirtualMachineConditionType `json:"type"`
	Status k8sv1.ConditionStatus       `json:"status"`
	// +nullable
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

//
// +k8s:openapi-gen=true
type VirtualMachineConditionType string

const (
	// VirtualMachineFailure is added in a virtual machine when its vmi
	// fails to be created due to insufficient quota, limit ranges, pod security policy, node selectors,
	// etc. or deleted due to kubelet being down or finalizers are failing.
	VirtualMachineFailure VirtualMachineConditionType = "Failure"

	// VirtualMachineReady is copied to the virtual machine from its vmi
	VirtualMachineReady VirtualMachineConditionType = "Ready"

	// VirtualMachinePaused is added in a virtual machine when its vmi
	// signals with its own condition that it is paused.
	VirtualMachinePaused VirtualMachineConditionType = "Paused"

	// This condition indicates that the VM was renamed
	RenameConditionType VirtualMachineConditionType = "RenameOperation"
)

//
// +k8s:openapi-gen=true
type HostDiskType string

const (
	// if disk does not exist at the given path,
	// a disk image will be created there
	HostDiskExistsOrCreate HostDiskType = "DiskOrCreate"
	// a disk image must exist at given disk path
	HostDiskExists HostDiskType = "Disk"
)

//
// +k8s:openapi-gen=true
type NetworkInterfaceType string

const (
	// Virtual machine instance bride interface
	BridgeInterface NetworkInterfaceType = "bridge"
	// Virtual machine instance slirp interface
	SlirpInterface NetworkInterfaceType = "slirp"
	// Virtual machine instance masquerade interface
	MasqueradeInterface NetworkInterfaceType = "masquerade"
)

//
// +k8s:openapi-gen=true
type DriverCache string

//
// +k8s:openapi-gen=true
type DriverIO string

const (
	// CacheNone - I/O from the guest is not cached on the host, but may be kept in a writeback disk cache.
	CacheNone DriverCache = "none"
	// CacheWriteThrough - I/O from the guest is cached on the host but written through to the physical medium.
	CacheWriteThrough DriverCache = "writethrough"

	// IOThreads - User mode based threads with a shared lock that perform I/O tasks. Can impact performance but offers
	// more predictable behaviour. This method is also takes fewer CPU cycles to submit I/O requests.
	IOThreads DriverIO = "threads"
	// IONative - Kernel native I/O tasks (AIO) offer a better performance but can block the VM if the file is not fully
	// allocated so this method recommended only when the backing file/disk/etc is fully preallocated.
	IONative DriverIO = "native"
	// IODefault - Fallback to the default value from the kernel. With recent Kernel versions (for example RHEL-7) the
	// default is AIO.
	IODefault DriverIO = "default"
)

// Handler defines a specific action that should be taken
// TODO: pass structured data to these actions, and document that data here.
type Handler struct {
	// HTTPGet specifies the http request to perform.
	// +optional
	HTTPGet *k8sv1.HTTPGetAction `json:"httpGet,omitempty"`
	// TCPSocket specifies an action involving a TCP port.
	// TCP hooks not yet supported
	// TODO: implement a realistic TCP lifecycle hook
	// +optional
	TCPSocket *k8sv1.TCPSocketAction `json:"tcpSocket,omitempty"`
}

// Probe describes a health check to be performed against a VirtualMachineInstance to determine whether it is
// alive or ready to receive traffic.
// +k8s:openapi-gen=true
type Probe struct {
	// The action taken to determine the health of a VirtualMachineInstance
	Handler `json:",inline"`
	// Number of seconds after the VirtualMachineInstance has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness. Minimum value is 1.
	// +optional
	SuccessThreshold int32 `json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

// KubeVirt represents the object deploying all KubeVirt resources
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type KubeVirt struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KubeVirtSpec   `json:"spec" valid:"required"`
	Status            KubeVirtStatus `json:"status,omitempty"`
}

// KubeVirtList is a list of KubeVirts
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type KubeVirtList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeVirt `json:"items"`
}

// ---
// +k8s:openapi-gen=true
type KubeVirtSelfSignConfiguration struct {
	CARotateInterval   *metav1.Duration `json:"caRotateInterval,omitempty"`
	CertRotateInterval *metav1.Duration `json:"certRotateInterval,omitempty"`
	CAOverlapInterval  *metav1.Duration `json:"caOverlapInterval,omitempty"`
}

// ---
// +k8s:openapi-gen=true
type KubeVirtCertificateRotateStrategy struct {
	SelfSigned *KubeVirtSelfSignConfiguration `json:"selfSigned,omitempty"`
}

// ---
// +k8s:openapi-gen=true
type KubeVirtSpec struct {
	// The image tag to use for the continer images installed.
	// Defaults to the same tag as the operator's container image.
	ImageTag string `json:"imageTag,omitempty"`
	// The image registry to pull the container images from
	// Defaults to the same registry the operator's container image is pulled from.
	ImageRegistry string `json:"imageRegistry,omitempty"`

	// The ImagePullPolicy to use.
	ImagePullPolicy k8sv1.PullPolicy `json:"imagePullPolicy,omitempty" valid:"required"`

	// The namespace Prometheus is deployed in
	// Defaults to openshift-monitor
	MonitorNamespace string `json:"monitorNamespace,omitempty"`

	// The name of the Prometheus service account that needs read-access to KubeVirt endpoints
	// Defaults to prometheus-k8s
	MonitorAccount string `json:"monitorAccount,omitempty"`

	// Specifies if kubevirt can be deleted if workloads are still present.
	// This is mainly a precaution to avoid accidental data loss
	UninstallStrategy KubeVirtUninstallStrategy `json:"uninstallStrategy,omitempty"`

	CertificateRotationStrategy KubeVirtCertificateRotateStrategy `json:"certificateRotateStrategy,omitempty"`

	// Designate the apps.kubevirt.io/version label for KubeVirt components.
	// Useful if KubeVirt is included as part of a product.
	// If ProductVersion is not specified, KubeVirt's version will be used.
	ProductVersion string `json:"productVersion,omitempty"`

	// Designate the apps.kubevirt.io/part-of label for KubeVirt components.
	// Useful if KubeVirt is included as part of a product.
	// If ProductName is not specified, the part-of label will be omitted.
	ProductName string `json:"productName,omitempty"`

	// holds kubevirt configurations.
	// same as the virt-configMap
	Configuration KubeVirtConfiguration `json:"configuration,omitempty"`

	// selectors and tolerations that should apply to KubeVirt infrastructure components
	// +optional
	Infra *ComponentConfig `json:"infra,omitempty"`

	// selectors and tolerations that should apply to KubeVirt workloads
	// +optional
	Workloads *ComponentConfig `json:"workloads,omitempty"`

	CustomizeComponents CustomizeComponents `json:"customizeComponents,omitempty"`
}

// +k8s:openapi-gen=true
type CustomizeComponents struct {
	// +listType=atomic
	Patches []CustomizeComponentsPatch `json:"patches,omitempty"`
}

// +k8s:openapi-gen=true
type CustomizeComponentsPatch struct {
	ResourceName string    `json:"resourceName,omitempty"`
	ResourceType string    `json:"resourceType,omitempty"`
	Patch        string    `json:"patch,omitempty"`
	Type         PatchType `json:"type,omitempty"`
}

type PatchType string

const (
	JSONPatchType           PatchType = "json"
	MergePatchType          PatchType = "merge"
	StrategicMergePatchType PatchType = "strategic"
)

type KubeVirtUninstallStrategy string

const (
	KubeVirtUninstallStrategyRemoveWorkloads                KubeVirtUninstallStrategy = "RemoveWorkloads"
	KubeVirtUninstallStrategyBlockUninstallIfWorkloadsExist KubeVirtUninstallStrategy = "BlockUninstallIfWorkloadsExist"
)

// KubeVirtStatus represents information pertaining to a KubeVirt deployment.
//
// +k8s:openapi-gen=true
type KubeVirtStatus struct {
	Phase                    KubeVirtPhase       `json:"phase,omitempty"`
	Conditions               []KubeVirtCondition `json:"conditions,omitempty" optional:"true"`
	OperatorVersion          string              `json:"operatorVersion,omitempty" optional:"true"`
	TargetKubeVirtRegistry   string              `json:"targetKubeVirtRegistry,omitempty" optional:"true"`
	TargetKubeVirtVersion    string              `json:"targetKubeVirtVersion,omitempty" optional:"true"`
	TargetDeploymentConfig   string              `json:"targetDeploymentConfig,omitempty" optional:"true"`
	TargetDeploymentID       string              `json:"targetDeploymentID,omitempty" optional:"true"`
	ObservedKubeVirtRegistry string              `json:"observedKubeVirtRegistry,omitempty" optional:"true"`
	ObservedKubeVirtVersion  string              `json:"observedKubeVirtVersion,omitempty" optional:"true"`
	ObservedDeploymentConfig string              `json:"observedDeploymentConfig,omitempty" optional:"true"`
	ObservedDeploymentID     string              `json:"observedDeploymentID,omitempty" optional:"true"`
}

// KubeVirtPhase is a label for the phase of a KubeVirt deployment at the current time.
//
// +k8s:openapi-gen=true
type KubeVirtPhase string

// These are the valid KubeVirt deployment phases
const (
	// The deployment is processing
	KubeVirtPhaseDeploying KubeVirtPhase = "Deploying"
	// The deployment succeeded
	KubeVirtPhaseDeployed KubeVirtPhase = "Deployed"
	// The deletion is processing
	KubeVirtPhaseDeleting KubeVirtPhase = "Deleting"
	// The deletion succeeeded
	KubeVirtPhaseDeleted KubeVirtPhase = "Deleted"
)

// KubeVirtCondition represents a condition of a KubeVirt deployment
//
// +k8s:openapi-gen=true
type KubeVirtCondition struct {
	Type   KubeVirtConditionType `json:"type"`
	Status k8sv1.ConditionStatus `json:"status"`
	// +optional
	// +nullable
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +optional
	// +nullable
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Message            string      `json:"message,omitempty"`
}

//
// +k8s:openapi-gen=true
type KubeVirtConditionType string

// These are the valid KubeVirt condition types
const (
	// Whether the deployment or deletion was successful (only used if false)
	KubeVirtConditionSynchronized KubeVirtConditionType = "Synchronized"
	// Whether all resources were created and up-to-date
	KubeVirtConditionCreated KubeVirtConditionType = "Created"

	// Conditions for HCO, see https://github.com/kubevirt/hyperconverged-cluster-operator/blob/master/docs/conditions.md
	// Whether KubeVirt is functional and available in the cluster.
	KubeVirtConditionAvailable KubeVirtConditionType = "Available"
	// Whether the operator is actively making changes to KubeVirt
	KubeVirtConditionProgressing KubeVirtConditionType = "Progressing"
	// Whether KubeVirt is not functioning completely
	KubeVirtConditionDegraded KubeVirtConditionType = "Degraded"
)

const (
	EvictionStrategyLiveMigrate EvictionStrategy = "LiveMigrate"
)

// RestartOptions may be provided when deleting an API object.
//
// +k8s:openapi-gen=true
type RestartOptions struct {
	metav1.TypeMeta `json:",inline"`

	// The duration in seconds before the object should be force-restared. Value must be non-negative integer.
	// The value zero indicates, restart immediately. If this value is nil, the default grace period for deletion of the corresponding VMI for the
	// specified type will be used to determine on how much time to give the VMI to restart.
	// Defaults to a per object value if not specified. zero means restart immediately.
	// Allowed Values: nil and 0
	// +optional
	GracePeriodSeconds *int64 `json:"gracePeriodSeconds,omitempty" protobuf:"varint,1,opt,name=gracePeriodSeconds"`
}

// VirtualMachineInstanceGuestAgentInfo represents information from the installed guest agent
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceGuestAgentInfo struct {
	metav1.TypeMeta `json:",inline"`
	// GAVersion is a version of currently installed guest agent
	GAVersion string `json:"guestAgentVersion,omitempty"`
	// Hostname represents FQDN of a guest
	Hostname string `json:"hostname,omitempty"`
	// OS contains the guest operating system information
	OS VirtualMachineInstanceGuestOSInfo `json:"os,omitempty"`
	// Timezone is guest os current timezone
	Timezone string `json:"timezone,omitempty"`
	// UserList is a list of active guest OS users
	UserList []VirtualMachineInstanceGuestOSUser `json:"userList,omitempty"`
	// FSInfo is a guest os filesystem information containing the disk mapping and disk mounts with usage
	FSInfo VirtualMachineInstanceFileSystemInfo `json:"fsInfo,omitempty"`
}

// VirtualMachineInstanceGuestOSUserList comprises the list of all active users on guest machine
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceGuestOSUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineInstanceGuestOSUser `json:"items"`
}

// VirtualMachineGuestOSUser is the single user of the guest os
// +k8s:openapi-gen=true
type VirtualMachineInstanceGuestOSUser struct {
	UserName  string  `json:"userName"`
	Domain    string  `json:"domain,omitempty"`
	LoginTime float64 `json:"loginTime,omitempty"`
}

// VirtualMachineInstanceFileSystemInfo represents information regarding single guest os filesystem
// +k8s:openapi-gen=true
type VirtualMachineInstanceFileSystemInfo struct {
	Filesystems []VirtualMachineInstanceFileSystem `json:"disks"`
}

// VirtualMachineInstanceFileSystemList comprises the list of all filesystems on guest machine
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type VirtualMachineInstanceFileSystemList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineInstanceFileSystem `json:"items"`
}

// VirtualMachineInstanceFileSystem represents guest os disk
// +k8s:openapi-gen=true
type VirtualMachineInstanceFileSystem struct {
	DiskName       string `json:"diskName"`
	MountPoint     string `json:"mountPoint"`
	FileSystemType string `json:"fileSystemType"`
	UsedBytes      int    `json:"usedBytes"`
	TotalBytes     int    `json:"totalBytes"`
}

// Options for a rename operation
type RenameOptions struct {
	metav1.TypeMeta `json:",inline"`
	NewName         string  `json:"newName"`
	OldName         *string `json:"oldName,omitempty"`
}

// KubeVirtConfiguration holds all kubevirt configurations
// +k8s:openapi-gen=true
type KubeVirtConfiguration struct {
	CPUModel                    string                  `json:"cpuModel,omitempty"`
	CPURequest                  *resource.Quantity      `json:"cpuRequest,string,omitempty"`
	DeveloperConfiguration      *DeveloperConfiguration `json:"developerConfiguration,omitempty"`
	EmulatedMachines            []string                `json:"emulatedMachines,omitempty"`
	ImagePullPolicy             k8sv1.PullPolicy        `json:"imagePullPolicy,omitempty"`
	MigrationConfiguration      *MigrationConfiguration `json:"migrations,omitempty"`
	MachineType                 string                  `json:"machineType,omitempty"`
	NetworkConfiguration        *NetworkConfiguration   `json:"network,omitempty"`
	OVMFPath                    string                  `json:"ovmfPath,omitempty"`
	SELinuxLauncherType         string                  `json:"selinuxLauncherType,omitempty"`
	SMBIOSConfig                *SMBiosConfiguration    `json:"smbios,omitempty"`
	SupportedGuestAgentVersions []string                `json:"supportedGuestAgentVersions,omitempty"`
	MemBalloonStatsPeriod       int                     `json:"memBalloonStatsPeriod,omitempty"`
}

// ---
// +k8s:openapi-gen=true
type SMBiosConfiguration struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Product      string `json:"product,omitempty"`
	Version      string `json:"version,omitempty"`
	Sku          string `json:"sku,omitempty"`
	Family       string `json:"family,omitempty"`
}

// MigrationConfiguration holds migration options
// +k8s:openapi-gen=true
type MigrationConfiguration struct {
	NodeDrainTaintKey                 *string `json:"nodeDrainTaintKey,omitempty"`
	ParallelOutboundMigrationsPerNode *uint32 `json:"parallelOutboundMigrationsPerNode,string,omitempty"`
	ParallelMigrationsPerCluster      *uint32 `json:"parallelMigrationsPerCluster,string,omitempty"`
	*PerVMIMigrationConfiguration     `json:",inline"`
}

// PerVMIMigrationConfiguration holds migration options that can be configurated per vmi
// +k8s:openapi-gen=true
type PerVMIMigrationConfiguration struct {
	AllowAutoConverge       *bool              `json:"allowAutoConverge,string,omitempty"`
	BandwidthPerMigration   *resource.Quantity `json:"bandwidthPerMigration,omitempty"`
	CompletionTimeoutPerGiB *int64             `json:"completionTimeoutPerGiB,string,omitempty"`
	ProgressTimeout         *int64             `json:"progressTimeout,string,omitempty"`
	UnsafeMigrationOverride *bool              `json:"unsafeMigrationOverride,string,omitempty"`
	MigrationMode           *MigrationMode     `json:"migrationMode,omitempty"`
}

// DeveloperConfiguration holds developer options
// +k8s:openapi-gen=true
type DeveloperConfiguration struct {
	FeatureGates           []string          `json:"featureGates,omitempty"`
	LessPVCSpaceToleration int               `json:"pvcTolerateLessSpaceUpToPercent,string,omitempty"`
	MemoryOvercommit       int               `json:"memoryOvercommit,string,omitempty"`
	NodeSelectors          map[string]string `json:"nodeSelectors,omitempty"`
	UseEmulation           bool              `json:"useEmulation,string,omitempty"`
	CPUAllocationRatio     float64           `json:"cpuAllocationRatio,string,omitempty"`
}

// NetworkConfiguration holds network options
// +k8s:openapi-gen=true
type NetworkConfiguration struct {
	NetworkInterface                  string `json:"defaultNetworkInterface,omitempty"`
	PermitSlirpInterface              bool   `json:"permitSlirpInterface,string,omitempty"`
	PermitBridgeInterfaceOnPodNetwork bool   `json:"permitBridgeInterfaceOnPodNetwork,string,omitempty"`
}
