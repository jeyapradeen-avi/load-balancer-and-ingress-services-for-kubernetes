/*
* [2013] - [2019] Avi Networks Incorporated
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

package lib

const (
	DISABLE_STATIC_ROUTE_SYNC = "DISABLE_STATIC_ROUTE_SYNC"
	CNI_PLUGIN                = "CNI_PLUGIN"
	CALICO_CNI                = "calico"
	INGRESSCOREV1             = "IngressCoreV1"
	INGRESSEXTV1              = "IngressExtV1"
	INGRESS_API               = "INGRESS_API"
	AviConfigMap              = "avi-k8s-config"
	AviNS                     = "avi-system"
	INGRESS_CLASS_ANNOT       = "kubernetes.io/ingress.class"
	AVI_INGRESS_CLASS         = "avi"
	SUBNET_IP                 = "SUBNET_IP"
	SUBNET_PREFIX             = "SUBNET_PREFIX"
	NETWORK_NAME              = "NETWORK_NAME"
	L7_SHARD_SCHEME           = "L7_SHARD_SCHEME"
	DEFAULT_SHARD_SCHEME      = "hostname"
	HOSTNAME_SHARD_SCHEME     = "hostname"
	NAMESPACE_SHARD_SCHEME    = "namespace"
	SLOW_RETRY_LAYER          = "SlowRetryLayer"
	FAST_RETRY_LAYER          = "FastRetryLayer"
	NOT_FOUND                 = "HTTP code: 404"
	STATUS_REDIRECT           = "HTTP_REDIRECT_STATUS_CODE_302"
	SLOW_SYNC_TIME            = 120
)
