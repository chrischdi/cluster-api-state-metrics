// SPDX-License-Identifier: MIT

package store

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

func TestMachineSetStore(t *testing.T) {
	startTime := 1501569018
	metav1StartTime := metav1.Unix(int64(startTime), 0)

	cases := []generateMetricsTestCase{
		{
			Obj: &clusterv1.MachineSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "ms1",
					Namespace:         "ns1",
					CreationTimestamp: metav1StartTime,
					ResourceVersion:   "10596",
					UID:               types.UID("foo"),
					OwnerReferences: []metav1.OwnerReference{
						{
							Controller: pointer.Bool(true),
							Kind:       "foo",
							Name:       "bar",
						},
					},
				},
			},
			Want: `
				# HELP capi_machineset_created Unix creation timestamp
				# HELP capi_machineset_labels Kubernetes labels converted to Prometheus labels.
				# HELP capi_machineset_owner Information about the machineset's owner.
				# TYPE capi_machineset_created gauge
				# TYPE capi_machineset_labels gauge
				# TYPE capi_machineset_owner gauge
				capi_machineset_created{machineset="ms1",namespace="ns1",uid="foo"} 1.501569018e+09
				capi_machineset_labels{machineset="ms1",namespace="ns1",uid="foo"} 1
				capi_machineset_owner{machineset="ms1",namespace="ns1",owner_is_controller="true",owner_kind="foo",owner_name="bar",uid="foo"} 1
			`,
			MetricNames: []string{"capi_machineset_labels", "capi_machineset_created", "capi_machineset_owner"},
		},
		{
			Obj: &clusterv1.MachineSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "ms1",
					Namespace:         "ns1",
					CreationTimestamp: metav1StartTime,
					ResourceVersion:   "10596",
					UID:               types.UID("foo"),
				},
				Spec: clusterv1.MachineSetSpec{
					Replicas: pointer.Int32(3),
				},
				Status: clusterv1.MachineSetStatus{
					Replicas:             3,
					FullyLabeledReplicas: 2,
					ReadyReplicas:        1,
					AvailableReplicas:    1,
				},
			},
			Want: `
				# HELP capi_machineset_spec_replicas Number of desired replicas for a machineset.
				# HELP capi_machineset_status_available_replicas The number of available replicas per machineset.
				# HELP capi_machineset_status_fully_labeled_replicas The number of fully labeled replicas per machineset.
				# HELP capi_machineset_status_ready_replicas The number of ready replicas per machineset.
				# HELP capi_machineset_status_replicas The number of replicas per machineset.
				# TYPE capi_machineset_spec_replicas gauge
				# TYPE capi_machineset_status_available_replicas gauge
				# TYPE capi_machineset_status_fully_labeled_replicas gauge
				# TYPE capi_machineset_status_ready_replicas gauge
				# TYPE capi_machineset_status_replicas gauge
				capi_machineset_spec_replicas{machineset="ms1",namespace="ns1",uid="foo"} 3
				capi_machineset_status_available_replicas{machineset="ms1",namespace="ns1",uid="foo"} 1
				capi_machineset_status_fully_labeled_replicas{machineset="ms1",namespace="ns1",uid="foo"} 2
				capi_machineset_status_ready_replicas{machineset="ms1",namespace="ns1",uid="foo"} 1
				capi_machineset_status_replicas{machineset="ms1",namespace="ns1",uid="foo"} 3
			`,
			MetricNames: []string{"capi_machineset_status_replicas", "capi_machineset_status_fully_labeled_replicas", "capi_machineset_status_ready_replicas", "capi_machineset_status_available_replicas", "capi_machineset_spec_replicas"},
		},
	}
	for i, c := range cases {
		c.Func = generator.ComposeMetricGenFuncs(machineSetMetricFamilies(nil))
		c.Headers = generator.ExtractMetricFamilyHeaders(machineSetMetricFamilies(nil))
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
