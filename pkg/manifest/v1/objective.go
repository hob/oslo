/*
Package v1 contains API Schema definitions for the slo v1 API group.

# Copyright © 2022 OpenSLO Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

// Calendar struct represents calendar time window.
type Calendar struct {
	StartTime string `yaml:"startTime" validate:"required,dateWithTime" example:"2020-01-21 12:30:00"`
	TimeZone  string `yaml:"timeZone" validate:"required,timeZone" example:"America/New_York"`
}

// TimeWindow represents content of time window.
type TimeWindow struct {
	Duration  string    `yaml:"duration" validate:"required,validDuration" example:"1h"`
	IsRolling bool      `yaml:"isRolling" example:"true"`
	Calendar  *Calendar `yaml:"calendar,omitempty" validate:"required_if=IsRolling false"`
}

type Target float64

func (t Target) MarshalYAML() (interface{}, error) {
	//Max 5 decimal places
	value := fmt.Sprintf("%.5F", t)
	value = strings.TrimRight(value, "0")
	return yaml.Node{
		Kind:        8,
		Style:       32,
		Tag:         "",
		Value:       value,
		Anchor:      "",
		Alias:       nil,
		Content:     nil,
		HeadComment: "",
		LineComment: "",
		FootComment: "",
		Line:        0,
		Column:      0,
	}, nil
}

// Objective represents single threshold for SLO, for internal usage.
type Objective struct {
	DisplayName     string `yaml:"displayName,omitempty"`
	Op              string `yaml:"op,omitempty" example:"lte"`
	Value           Target `yaml:"value,omitempty" validate:"numeric,omitempty"`
	Target          Target `yaml:"target" validate:"required,numeric,gte=0,lt=1" example:"0.9"`
	TimeSliceTarget Target `yaml:"timeSliceTarget,omitempty" validate:"gte=0,lte=1,omitempty" example:"0.9"`
	TimeSliceWindow string `yaml:"timeSliceWindow,omitempty" example:"5m"`
}

// SLOSpec struct which mapped one to one with kind: slo yaml definition, internal use.
type SLOSpec struct {
	Description     string        `yaml:"description,omitempty" validate:"max=1050,omitempty"`
	Service         string        `yaml:"service" validate:"required" example:"webapp-service"`
	Indicator       *SLIInline    `yaml:"indicator,omitempty" validate:"required_without=IndicatorRef"`
	IndicatorRef    *string       `yaml:"indicatorRef,omitempty"`
	BudgetingMethod string        `yaml:"budgetingMethod" validate:"required,oneof=Occurrences Timeslices" example:"Occurrences"` //nolint:lll
	TimeWindow      []TimeWindow  `yaml:"timeWindow" validate:"required,len=1,dive"`
	Objectives      []Objective   `yaml:"objectives" validate:"required,dive"`
	AlertPolicies   []AlertPolicy `yaml:"alertPolicies" validate:"dive"`
}

// SLO struct which mapped one to one with kind: slo yaml definition, external usage.
type SLO struct {
	ObjectHeader `yaml:",inline"`
	Spec         SLOSpec `yaml:"spec" validate:"required"`
}

// Kind returns the name of this type.
func (SLO) Kind() string {
	return "SLO"
}
