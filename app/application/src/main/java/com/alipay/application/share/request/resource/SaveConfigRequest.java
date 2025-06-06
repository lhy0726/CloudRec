/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.alipay.application.share.request.resource;

import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import lombok.Getter;
import lombok.Setter;

import java.util.List;

/*
 *@title SaveConfigRequest
 *@description
 *@author jietian
 *@version 1.0
 *@create 2024/9/10 19:27
 */
@Getter
@Setter
public class SaveConfigRequest {

    /**
     * id 保存时需要
     */
    private Long id;

    /**
     * 订阅名称
     */
    @NotEmpty(message = "订阅名称不能为空")
    private String name;

    /**
     * action
     */
    @NotEmpty(message = "actionList不能为空")
    private List<Action> actionList;

    /**
     * 条件
     */
    private String condition;

    /**
     * 配置的风险规则
     */
    @NotNull(message = "配置的风险规则不能为空")
    private List<Config> ruleConfigList;

    @Getter
    @Setter
    public static class Config {
        private Integer id;
        private String key;
        private String keyName;
        private Operator operator;
        private Object value;
    }

    public enum Operator {

        ANY, ALL, EQ, ALL_IN, ANY_IN;
    }

    @Getter
    @Setter
    public static class Action {

        /**
         * 告警类型：定时告警、即时告警
         */
        @NotEmpty(message = "告警类型不能为空")
        private String actionType;

        /**
         * 钉钉群、企业微信、接口回调
         */
        @NotEmpty(message = "告警方式不能为空")
        private String action;

        /**
         * 群名称
         */
        private String name;

        /**
         * url
         */
        @NotEmpty(message = "url不能为空")
        private String url;

        /**
         * 周期：周1、周2、周3...如果是每天都需要告警在不需要填此参数
         */
        private String period;

        /**
         * 告警时间列表10点、12点、14点
         */
        private List<String> timeList;
    }
}
