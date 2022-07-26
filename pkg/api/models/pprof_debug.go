// +build !release

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

package models

import (
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/utils"
)

var Debug *Debugging
var once sync.Once

type Debugging struct{}

func (p *Debugging) InitModel() {
	utils.AviLog.Debugf("Debugging APIs available")
	once.Do(func() {
		// TODO
	})
}

func (p *Debugging) ApiOperationMap() []OperationMap {
	var operationMapList []OperationMap

	heapOperation := OperationMap{
		Route:  "/api/pprof/heap",
		Method: "GET",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			pprof.Handler("heap").ServeHTTP(w, r)
		},
	}

	operationMapList = append(operationMapList, heapOperation)
	return operationMapList
}
