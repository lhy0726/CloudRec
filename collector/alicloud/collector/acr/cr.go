// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package acr

import (
	"context"
	cr20181201 "github.com/alibabacloud-go/cr-20181201/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/cloudrec/alicloud/collector"
	"github.com/core-sdk/constant"
	"github.com/core-sdk/log"
	"github.com/core-sdk/schema"
	"go.uber.org/zap"
)

func GetCRResource() schema.Resource {
	return schema.Resource{
		ResourceType:       collector.ACR,
		ResourceTypeName:   collector.ACR,
		ResourceGroupType:  constant.CONTAINER,
		Desc:               `https://api.aliyun.com/product/cr`,
		ResourceDetailFunc: GetInstanceDetail,
		RowField: schema.RowField{
			ResourceId:   "$.Instance.InstanceId",
			ResourceName: "$.Instance.InstanceName",
		},
		Dimension: schema.Regional,
	}
}

func GetInstanceDetail(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
	cli := service.(*collector.Services).ACR

	var page int32 = 1
	var pageSize int32 = 30

	request := &cr20181201.ListInstanceRequest{}
	request.PageNo = tea.Int32(page)
	request.PageSize = tea.Int32(pageSize)

	for {
		resp, err := cli.ListInstance(request)
		if err != nil {
			log.CtxLogger(ctx).Error("ListInstance error", zap.Error(err))
			return err
		}
		for i := 0; i < len(resp.Body.Instances); i++ {
			res <- Detail{
				Instance:                 resp.Body.Instances[i],
				InstanceInternetEndpoint: getInstanceEndpoint(ctx, cli, resp.Body.Instances[i].InstanceId),
				Repository:               getRepository(ctx, cli, resp.Body.Instances[i].InstanceId),
			}
		}

		if resp.Body.Instances == nil || int32(len(resp.Body.Instances)) < pageSize {
			break
		}

		page += 1
		request.PageNo = tea.Int32(page)
	}
	return nil
}

type Detail struct {
	Instance                 *cr20181201.ListInstanceResponseBodyInstances
	InstanceInternetEndpoint *cr20181201.GetInstanceEndpointResponseBody
	Repository               *cr20181201.GetRepositoryResponseBody
}

func getInstanceEndpoint(ctx context.Context, cli *cr20181201.Client, instanceId *string) (res *cr20181201.GetInstanceEndpointResponseBody) {
	getInstanceEndpointRequest := &cr20181201.GetInstanceEndpointRequest{
		// EndpointType only support "internet"
		EndpointType: tea.String("internet"),
		InstanceId:   instanceId,
	}
	resp, err := cli.GetInstanceEndpoint(getInstanceEndpointRequest)

	if err != nil {
		log.CtxLogger(ctx).Error("GetInstanceEndpoint error", zap.Error(err))
		return nil
	}
	return resp.Body
}

// https://api.aliyun.com/api/cr/2018-12-01/GetRepository?tab=DEMO&lang=GO
func getRepository(ctx context.Context, cli *cr20181201.Client, instanceId *string) (res *cr20181201.GetRepositoryResponseBody) {
	request := &cr20181201.GetRepositoryRequest{
		InstanceId: instanceId,
	}
	resp, err := cli.GetRepository(request)

	if err != nil {
		log.CtxLogger(ctx).Error("getRepository error", zap.Error(err))
		return nil
	}
	return resp.Body
}
