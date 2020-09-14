/*
 * Copyright 2019-2020 VMware, Inc.
 * All Rights Reserved.
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*   http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

package k8s

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/lib"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/nodes"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/objects"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/internal/status"
	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/utils"

	advl4v1alpha1pre1 "github.com/vmware-tanzu/service-apis/apis/v1alpha1pre1"
	advl4crd "github.com/vmware-tanzu/service-apis/pkg/client/clientset/versioned"
	advl4informer "github.com/vmware-tanzu/service-apis/pkg/client/informers/externalversions"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func NewAdvL4Informers(cs advl4crd.Interface) {
	var advl4InformerFactory advl4informer.SharedInformerFactory

	advl4InformerFactory = advl4informer.NewSharedInformerFactoryWithOptions(cs, time.Second*30)
	gatewayInformer := advl4InformerFactory.Networking().V1alpha1pre1().Gateways()
	gatewayClassInformer := advl4InformerFactory.Networking().V1alpha1pre1().GatewayClasses()

	lib.SetAdvL4Informers(&lib.AdvL4Informers{
		GatewayInformer:      gatewayInformer,
		GatewayClassInformer: gatewayClassInformer,
	})
}

// SetupAdvL4EventHandlers handles setting up of AdvL4 event handlers
func (c *AviController) SetupAdvL4EventHandlers(numWorkers uint32) {
	utils.AviLog.Infof("Setting up AdvL4 Event handlers")
	informer := lib.GetAdvL4Informers()

	gatewayEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if c.DisableSync {
				return
			}
			gw := obj.(*advl4v1alpha1pre1.Gateway)
			namespace, _, _ := cache.SplitMetaNamespaceKey(utils.ObjKey(gw))
			key := lib.Gateway + "/" + utils.ObjKey(gw)
			utils.AviLog.Infof("key: %s, msg: ADD", key)

			status.InitializeGatewayConditions(gw)
			validateGatewayForStatusUpdates(key, gw)
			checkGWForGatewayPortConflict(key, gw)

			bkt := utils.Bkt(namespace, numWorkers)
			c.workqueue[bkt].AddRateLimited(key)
		},
		UpdateFunc: func(old, new interface{}) {
			if c.DisableSync {
				return
			}
			oldObj := old.(*advl4v1alpha1pre1.Gateway)
			gw := new.(*advl4v1alpha1pre1.Gateway)
			if !reflect.DeepEqual(oldObj.Spec, gw.Spec) {
				namespace, _, _ := cache.SplitMetaNamespaceKey(utils.ObjKey(gw))
				key := lib.Gateway + "/" + utils.ObjKey(gw)
				utils.AviLog.Infof("key: %s, msg: UPDATE", key)

				if ipChanged := checkGWForIPUpdate(key, gw, oldObj); ipChanged {
					return
				}
				validateGatewayForStatusUpdates(key, gw)
				checkGWForGatewayPortConflict(key, gw)

				bkt := utils.Bkt(namespace, numWorkers)
				c.workqueue[bkt].AddRateLimited(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			if c.DisableSync {
				return
			}
			gw := obj.(*advl4v1alpha1pre1.Gateway)
			namespace, _, _ := cache.SplitMetaNamespaceKey(utils.ObjKey(gw))
			key := lib.Gateway + "/" + utils.ObjKey(gw)
			utils.AviLog.Infof("key: %s, msg: DELETE", key)
			bkt := utils.Bkt(namespace, numWorkers)
			c.workqueue[bkt].AddRateLimited(key)
		},
	}

	gatewayClassEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if c.DisableSync {
				return
			}
			gwclass := obj.(*advl4v1alpha1pre1.GatewayClass)
			namespace, _, _ := cache.SplitMetaNamespaceKey(utils.ObjKey(gwclass))
			key := lib.GatewayClass + "/" + utils.ObjKey(gwclass)
			utils.AviLog.Infof("key: %s, msg: ADD", key)
			bkt := utils.Bkt(namespace, numWorkers)
			c.workqueue[bkt].AddRateLimited(key)
		},
		UpdateFunc: func(old, new interface{}) {
			if c.DisableSync {
				return
			}
			oldObj := old.(*advl4v1alpha1pre1.GatewayClass)
			gwclass := new.(*advl4v1alpha1pre1.GatewayClass)
			if !reflect.DeepEqual(oldObj.Spec, gwclass.Spec) {
				namespace, _, _ := cache.SplitMetaNamespaceKey(utils.ObjKey(gwclass))
				key := lib.GatewayClass + "/" + utils.ObjKey(gwclass)
				utils.AviLog.Infof("key: %s, msg: UPDATE", key)
				bkt := utils.Bkt(namespace, numWorkers)
				c.workqueue[bkt].AddRateLimited(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			if c.DisableSync {
				return
			}
			gwclass := obj.(*advl4v1alpha1pre1.GatewayClass)
			key := lib.GatewayClass + "/" + utils.ObjKey(gwclass)
			namespace, _, _ := cache.SplitMetaNamespaceKey(utils.ObjKey(gwclass))
			utils.AviLog.Infof("key: %s, msg: DELETE", key)
			bkt := utils.Bkt(namespace, numWorkers)
			c.workqueue[bkt].AddRateLimited(key)
		},
	}

	informer.GatewayInformer.Informer().AddEventHandler(gatewayEventHandler)
	informer.GatewayClassInformer.Informer().AddEventHandler(gatewayClassEventHandler)

	return
}

func validateGatewayForStatusUpdates(key string, gateway *advl4v1alpha1pre1.Gateway) {
	defer status.UpdateGatewayStatusObject(gateway, &gateway.Status)

	gwClassObj, err := lib.GetAdvL4Informers().GatewayClassInformer.Lister().Get(gateway.Spec.Class)
	if err != nil {
		status.UpdateGatewayStatusGWCondition(gateway, &status.UpdateGWStatusConditionOptions{
			Type:    "Pending",
			Status:  corev1.ConditionTrue,
			Message: fmt.Sprintf("Corresponding networking.x-k8s.io/gatewayclass not found %s", gateway.Spec.Class),
			Reason:  "InvalidGatewayClass",
		})
		utils.AviLog.Warnf("key: %s, msg: Corresponding networking.x-k8s.io/gatewayclass not found %s %v",
			key, gateway.Spec.Class, err)
		return
	}

	for _, listener := range gateway.Spec.Listeners {
		gwName, nameOk := listener.Routes.RouteSelector.MatchLabels[lib.GatewayNameLabelKey]
		gwNamespace, nsOk := listener.Routes.RouteSelector.MatchLabels[lib.GatewayNamespaceLabelKey]
		if !nameOk || !nsOk ||
			(nameOk && gwName != gateway.Name) ||
			(nsOk && gwNamespace != gateway.Namespace) {
			status.UpdateGatewayStatusGWCondition(gateway, &status.UpdateGWStatusConditionOptions{
				Type:    "Pending",
				Status:  corev1.ConditionTrue,
				Message: "Incorrect gateway matchLabels configuration",
				Reason:  "InvalidMatchLabels",
			})
			return
		}
	}

	// Additional check to see if the gatewayclass is a valid avi gateway class or not.
	if gwClassObj.Spec.Controller != lib.AviGatewayController {
		// Return an error since this is not our object.
		status.UpdateGatewayStatusGWCondition(gateway, &status.UpdateGWStatusConditionOptions{
			Type:    "Pending",
			Status:  corev1.ConditionTrue,
			Message: fmt.Sprintf("Unable to identify controller %s", gwClassObj.Spec.Controller),
			Reason:  "UnidentifiedController",
		})
	}
}

func checkSvcForGatewayPortConflict(svc *corev1.Service, key string) {
	gateway, portProtocols := nodes.ParseL4ServiceForGateway(svc, key)
	if gateway == "" {
		return
	}
	found, gwSvcListeners := objects.ServiceGWLister().GetGwToSvcs(gateway)
	if !found {
		return
	}

	gwNSName := strings.Split(gateway, "/")
	gw, err := lib.GetAdvL4Informers().GatewayInformer.Lister().Gateways(gwNSName[0]).Get(gwNSName[1])
	if err != nil {
		utils.AviLog.Warnf("key: %s, msg: Unable to find gateway: %v", key, err)
		return
	}

	// detect port conflict
	for _, portProtocol := range portProtocols {
		if val, ok := gwSvcListeners[portProtocol]; ok {
			if !utils.HasElem(val, svc.Namespace+"/"+svc.Name) {
				val = append(val, svc.Namespace+"/"+svc.Name)
			}
			if len(val) > 1 {
				portProtocolArr := strings.Split(portProtocol, "/")
				status.UpdateGatewayStatusListenerConditions(gw, portProtocolArr[1], &status.UpdateGWStatusConditionOptions{
					Type:   "PortConflict",
					Status: corev1.ConditionTrue,
					Reason: fmt.Sprintf("conflicting port configuration provided in service %s and %s/%s", val, svc.Namespace, svc.Name),
				})
				status.UpdateGatewayStatusObject(gw, &gw.Status)
				return
			}
		}
	}

	// detect unsupported protocol
	// TODO

	return
}

func checkGWForGatewayPortConflict(key string, gw *advl4v1alpha1pre1.Gateway) {
	found, gwSvcListeners := objects.ServiceGWLister().GetGwToSvcs(gw.Namespace + "/" + gw.Name)
	if !found {
		return
	}

	var gwProtocols []string
	// port conflicts
	for _, listener := range gw.Spec.Listeners {
		portProtoGW := string(listener.Protocol) + "/" + strconv.Itoa(int(listener.Port))
		if !utils.HasElem(gwProtocols, string(listener.Protocol)) {
			gwProtocols = append(gwProtocols, string(listener.Protocol))
		}

		if val, ok := gwSvcListeners[portProtoGW]; ok && len(val) > 1 {
			status.UpdateGatewayStatusListenerConditions(gw, strconv.Itoa(int(listener.Port)), &status.UpdateGWStatusConditionOptions{
				Type:   "PortConflict",
				Status: corev1.ConditionTrue,
				Reason: fmt.Sprintf("conflicting port configuration provided in service %s and %v", val, gwSvcListeners[portProtoGW]),
			})
			status.UpdateGatewayStatusObject(gw, &gw.Status)
			return
		}
	}

	// unsupported protocol
	for portProto, svcs := range gwSvcListeners {
		svcProtocol := strings.Split(portProto, "/")[0]
		if !utils.HasElem(gwProtocols, svcProtocol) {
			status.UpdateGatewayStatusListenerConditions(gw, strings.Split(portProto, "/")[1], &status.UpdateGWStatusConditionOptions{
				Type:   "UnsupportedProtocol",
				Status: corev1.ConditionTrue,
				Reason: fmt.Sprintf("Unsupported protocol found in services %v", svcs),
			})
			status.UpdateGatewayStatusObject(gw, &gw.Status)
			return
		}
	}
}

func checkGWForIPUpdate(key string, gw, oldGw *advl4v1alpha1pre1.Gateway) bool {
	var newIPAddress, oldIPAddress string
	if len(gw.Spec.Addresses) > 0 {
		newIPAddress = gw.Spec.Addresses[0].Value
	}
	if len(oldGw.Spec.Addresses) > 0 {
		oldIPAddress = oldGw.Spec.Addresses[0].Value
	}

	// honor new IP coming in
	// old: nil, new: X
	if oldIPAddress == "" && newIPAddress != "" {
		return false
	}

	// old: X, new: nil | old: X, new: Y
	if newIPAddress != oldIPAddress {
		errString := "IPAddress Updates on gateway not supported, Please recreate gateway object with the new preferred IPAddress"
		status.UpdateGatewayStatusGWCondition(gw, &status.UpdateGWStatusConditionOptions{
			Type:    "Pending",
			Status:  corev1.ConditionTrue,
			Reason:  "InvalidAddress",
			Message: errString,
		})
		utils.AviLog.Errorf("key: %s, msg: %s Last IPAddress: %s, Current IPAddress: %s",
			key, errString, oldIPAddress, newIPAddress)
		status.UpdateGatewayStatusObject(gw, &gw.Status)
		return true
	}

	return false
}
