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
package com.alipay.application.service.risk.engine.handler;

import com.alipay.application.service.risk.engine.ConditionItem;
import com.alipay.application.service.risk.engine.Fact;
import com.alipay.application.service.risk.engine.Operator;

import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

public class OperatorHandlerComplete {

    private static final List<OperatorHandler> operatorHandlers = Arrays.asList(new EqHandler(), new AllInHandler());

    public static boolean handle(ConditionItem conditionItem, List<Fact> facts) {
        if (conditionItem.getOperator() == null) {
            throw new RuntimeException("operator is null");
        }

        List<Fact> collect = facts.stream().filter(fact -> fact.getKey().equals(conditionItem.getKey()))
                .collect(Collectors.toList());
        if (collect.isEmpty()) {
            //If there is no judgment condition, the judgment condition will be ignored.
            return true;
        }
        if (conditionItem.getOperator() == Operator.EQ) {
            return new EqHandler().handle(conditionItem, collect);
        }

        if (conditionItem.getOperator() == Operator.ALL_IN) {
            return new AllInHandler().handle(conditionItem, collect);
        }

        if (conditionItem.getOperator() == Operator.ANY_IN) {
            return new AnyInHandler().handle(conditionItem, collect);
        }

        throw new RuntimeException("operator is not support");
    }
}
