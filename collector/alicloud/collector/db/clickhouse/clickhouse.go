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

package clickhouse

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/clickhouse"
	"github.com/cloudrec/alicloud/collector"
	"github.com/core-sdk/constant"
	"github.com/core-sdk/log"
	"github.com/core-sdk/schema"
	"go.uber.org/zap"
)

func GetClickHouseResource() schema.Resource {
	return schema.Resource{
		ResourceType:       collector.ClickHouse,
		ResourceTypeName:   collector.ClickHouse,
		ResourceGroupType:  constant.DATABASE,
		Desc:               `https://api.aliyun.com/product/clickhouse`,
		ResourceDetailFunc: GetInstanceDetail,
		RowField: schema.RowField{
			ResourceId:   "$.DBCluster.DBClusterId",
			ResourceName: "$.DBCluster.DBClusterDescription",
			Address:      "$.DBCluster.ConnectionString",
		},
		Dimension: schema.Regional,
	}
}

func GetInstanceDetail(ctx context.Context, service schema.ServiceInterface, res chan<- any) error {
	cli := service.(*collector.Services).Clickhouse

	request := clickhouse.CreateDescribeDBClustersRequest()
	request.PageSize = requests.NewInteger(constant.DefaultPageSize)
	request.PageNumber = requests.NewInteger(1)
	count := 0
	for {
		response, err := cli.DescribeDBClusters(request)
		if err != nil {
			log.CtxLogger(ctx).Error("DescribeDBClusters error", zap.Error(err))
			return err
		}
		for _, i := range response.DBClusters.DBCluster {
			res <- &Detail{
				DBCluster:   i,
				NetInfoItem: describeDBClusterNetInfoItems(ctx, cli, i.DBClusterId),
				IPArray:     describeDBClusterAccessWhiteList(ctx, cli, i.DBClusterId),
			}

		}
		count += len(response.DBClusters.DBCluster)
		if count >= response.TotalCount {
			break
		}
		request.PageNumber = requests.NewInteger(response.PageNumber + 1)
	}
	return nil
}

type Detail struct {
	endPoint string

	// cluster basic info
	DBCluster clickhouse.DBCluster

	// net info
	NetInfoItem []clickhouse.NetInfoItem

	// whitelist info
	IPArray []clickhouse.IPArray
}

// Query the network information of the specified cloud database ClickHouse cluster
func describeDBClusterNetInfoItems(ctx context.Context, cli *clickhouse.Client, DBClusterId string) (netInfoItem []clickhouse.NetInfoItem) {
	request := clickhouse.CreateDescribeDBClusterNetInfoItemsRequest()
	request.DBClusterId = DBClusterId
	resp, err := cli.DescribeDBClusterNetInfoItems(request)
	if err != nil {
		log.CtxLogger(ctx).Error("DescribeDBClusterNetInfoItems error", zap.Error(err))
		return
	}
	return resp.NetInfoItems.NetInfoItem
}

// View the IP whitelist of the specified cloud database ClickHouse cluster
func describeDBClusterAccessWhiteList(ctx context.Context, cli *clickhouse.Client, DBClusterId string) (IPArray []clickhouse.IPArray) {
	request := clickhouse.CreateDescribeDBClusterAccessWhiteListRequest()
	request.DBClusterId = DBClusterId
	resp, err := cli.DescribeDBClusterAccessWhiteList(request)
	if err != nil {
		log.CtxLogger(ctx).Error("DescribeDBClusterAccessWhiteList error", zap.Error(err))
		return
	}
	return resp.DBClusterAccessWhiteList.IPArray
}
