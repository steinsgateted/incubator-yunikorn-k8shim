/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package cache

import (
	"github.com/apache/yunikorn-k8shim/pkg/common/events"
	"github.com/apache/yunikorn-k8shim/pkg/log"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
	"github.com/looplab/fsm"
	"go.uber.org/zap"
)

//----------------------------------------------
// Application events
//----------------------------------------------
type ApplicationEventType int

const (
	SubmitApplication ApplicationEventType = iota
	RecoverApplication
	AcceptApplication
	TryReserve
	UpdateReservation
	RunApplication
	RejectApplication
	CompleteApplication
	FailApplication
	KillApplication
	KilledApplication
	ReleaseAppAllocation
	ReleaseAppAllocationAsk
	AppStateChange
	ResumingApplication
	AppTaskCompleted
)

func (ae ApplicationEventType) String() string {
	return [...]string{"SubmitApplication", "RecoverApplication", "AcceptApplication", "TryReserve", "UpdateReservation", "RunApplication", "RejectApplication", "CompleteApplication", "FailApplication", "KillApplication", "KilledApplication", "ReleaseAppAllocation", "ReleaseAppAllocationAsk", "AppStateChange", "ResumingApplication", "AppTaskCompleted"}[ae]
}

type ApplicationEvent interface {
	// an application event is associated with an application Id,
	// dispatcher finds out actual application based on this id
	// to handle this event
	GetApplicationID() string

	// the type of this event
	GetEvent() ApplicationEventType

	// an event can have multiple arguments, these arguments will be passed to
	// state machines' callbacks when doing state transition
	GetArgs() []interface{}
}

// SimpleApplicationEvent simply moves application states
type SimpleApplicationEvent struct {
	applicationID string
	event         ApplicationEventType
}

func NewSimpleApplicationEvent(appID string, eventType ApplicationEventType) SimpleApplicationEvent {
	return SimpleApplicationEvent{
		applicationID: appID,
		event:         eventType,
	}
}

func (st SimpleApplicationEvent) GetEvent() ApplicationEventType {
	return st.event
}

func (st SimpleApplicationEvent) GetArgs() []interface{} {
	return nil
}

func (st SimpleApplicationEvent) GetApplicationID() string {
	return st.applicationID
}

// ------------------------
// GeneralApplicationEvent is used for testing and rejected app events
// ------------------------
type GeneralApplicationEvent struct {
	applicationID string
	event         ApplicationEventType
	message       string
}

func NewGeneralApplicationEvent(appID string, eventType ApplicationEventType, msg string) GeneralApplicationEvent {
	return GeneralApplicationEvent{
		applicationID: appID,
		event:         eventType,
		message:       msg,
	}
}

func (st GeneralApplicationEvent) GetEvent() ApplicationEventType {
	return st.event
}

func (st GeneralApplicationEvent) GetArgs() []interface{} {
	args := make([]interface{}, 1)
	args[0] = st.message
	return args
}

func (st GeneralApplicationEvent) GetApplicationID() string {
	return st.applicationID
}

// ------------------------
// ApplicationStatusChangeEvent updates the status in the application CRD
// ------------------------
type ApplicationStatusChangeEvent struct {
	applicationID string
	event         ApplicationEventType
	state         string
}

func NewApplicationStatusChangeEvent(appID string, eventType ApplicationEventType, state string) ApplicationStatusChangeEvent {
	return ApplicationStatusChangeEvent{
		applicationID: appID,
		event:         eventType,
		state:         state,
	}
}

func (st ApplicationStatusChangeEvent) GetEvent() ApplicationEventType {
	return st.event
}

func (st ApplicationStatusChangeEvent) GetArgs() []interface{} {
	return nil
}

func (st ApplicationStatusChangeEvent) GetApplicationID() string {
	return st.applicationID
}

func (st ApplicationStatusChangeEvent) GetState() string {
	return st.state
}

// ------------------------
// SubmitTask application
// ------------------------
type SubmitApplicationEvent struct {
	applicationID string
	event         ApplicationEventType
}

func NewSubmitApplicationEvent(appID string) SubmitApplicationEvent {
	return SubmitApplicationEvent{
		applicationID: appID,
		event:         SubmitApplication,
	}
}

func (se SubmitApplicationEvent) GetEvent() ApplicationEventType {
	return se.event
}

func (se SubmitApplicationEvent) GetArgs() []interface{} {
	return nil
}

func (se SubmitApplicationEvent) GetApplicationID() string {
	return se.applicationID
}

// ------------------------
// Run application
// ------------------------
type RunApplicationEvent struct {
	applicationID string
	event         ApplicationEventType
}

func NewRunApplicationEvent(appID string) RunApplicationEvent {
	return RunApplicationEvent{
		applicationID: appID,
		event:         RunApplication,
	}
}

func (re RunApplicationEvent) GetEvent() ApplicationEventType {
	return re.event
}

func (re RunApplicationEvent) GetArgs() []interface{} {
	return nil
}

func (re RunApplicationEvent) GetApplicationID() string {
	return re.applicationID
}

// ------------------------
// Fail application
// ------------------------
type FailApplicationEvent struct {
	applicationID string
	event         ApplicationEventType
	errorMessage  string
}

func NewFailApplicationEvent(appID, message string) FailApplicationEvent {
	return FailApplicationEvent{
		applicationID: appID,
		event:         FailApplication,
		errorMessage:  message,
	}
}

func (fe FailApplicationEvent) GetEvent() ApplicationEventType {
	return fe.event
}

func (fe FailApplicationEvent) GetArgs() []interface{} {
	args := make([]interface{}, 1)
	args[0] = fe.errorMessage
	return args
}

func (fe FailApplicationEvent) GetApplicationID() string {
	return fe.applicationID
}

// ------------------------
// Reservation Update Event
// ------------------------
type UpdateApplicationReservationEvent struct {
	applicationID string
	event         ApplicationEventType
}

func NewUpdateApplicationReservationEvent(appID string) UpdateApplicationReservationEvent {
	return UpdateApplicationReservationEvent{
		applicationID: appID,
		event:         UpdateReservation,
	}
}

func (ue UpdateApplicationReservationEvent) GetEvent() ApplicationEventType {
	return ue.event
}

func (ue UpdateApplicationReservationEvent) GetArgs() []interface{} {
	return nil
}

func (ue UpdateApplicationReservationEvent) GetApplicationID() string {
	return ue.applicationID
}

// ------------------------
// Release application allocations
// ------------------------
type ReleaseAppAllocationEvent struct {
	applicationID   string
	allocationUUID  string
	terminationType string
	event           ApplicationEventType
}

func NewReleaseAppAllocationEvent(appID string, allocTermination si.TerminationType, uuid string) ReleaseAppAllocationEvent {
	return ReleaseAppAllocationEvent{
		applicationID:   appID,
		allocationUUID:  uuid,
		terminationType: si.TerminationType_name[int32(allocTermination)],
		event:           ReleaseAppAllocation,
	}
}

func (re ReleaseAppAllocationEvent) GetApplicationID() string {
	return re.applicationID
}

func (re ReleaseAppAllocationEvent) GetArgs() []interface{} {
	args := make([]interface{}, 2)
	args[0] = re.allocationUUID
	args[1] = re.terminationType
	return args
}

func (re ReleaseAppAllocationEvent) GetEvent() ApplicationEventType {
	return re.event
}

type ReleaseAppAllocationAskEvent struct {
	applicationID   string
	taskID          string
	terminationType string
	event           ApplicationEventType
}

func NewReleaseAppAllocationAskEvent(appID string, allocTermination si.TerminationType, taskID string) ReleaseAppAllocationAskEvent {
	return ReleaseAppAllocationAskEvent{
		applicationID:   appID,
		taskID:          taskID,
		terminationType: si.TerminationType_name[int32(allocTermination)],
		event:           ReleaseAppAllocationAsk,
	}
}

func (re ReleaseAppAllocationAskEvent) GetApplicationID() string {
	return re.applicationID
}

func (re ReleaseAppAllocationAskEvent) GetArgs() []interface{} {
	args := make([]interface{}, 2)
	args[0] = re.taskID
	args[1] = re.terminationType
	return args
}

func (re ReleaseAppAllocationAskEvent) GetEvent() ApplicationEventType {
	return re.event
}

// ------------------------
// Resuming application
// ------------------------
type ResumingApplicationEvent struct {
	applicationID string
	event         ApplicationEventType
}

func NewResumingApplicationEvent(appID string) ResumingApplicationEvent {
	return ResumingApplicationEvent{
		applicationID: appID,
		event:         ResumingApplication,
	}
}

func (re ResumingApplicationEvent) GetEvent() ApplicationEventType {
	return re.event
}

func (re ResumingApplicationEvent) GetArgs() []interface{} {
	return nil
}

func (re ResumingApplicationEvent) GetApplicationID() string {
	return re.applicationID
}

// ----------------------------------
// Application states
// ----------------------------------
var storeApplicationStates *applicationStates

type applicationStates struct {
	New        string
	Recovering string
	Submitted  string
	Accepted   string
	Reserving  string
	Running    string
	Rejected   string
	Completed  string
	Killing    string
	Killed     string
	Failing    string
	Failed     string
	Resuming   string
}

func ApplicationStates() *applicationStates{
	if storeApplicationStates == nil {
		storeApplicationStates = &applicationStates{
			New:        "New",
			Recovering: "Recovering",
			Submitted:  "Submitted",
			Accepted:   "Accepted",
			Reserving:  "Reserving",
			Running:    "Running",
			Rejected:   "Rejected",
			Completed:  "Completed",
			Killing:    "Killing",
			Killed:     "Killed",
			Failed:     "Failed",
			Failing:    "Failing",
			Resuming:   "Resuming",
		}
	}
	return storeApplicationStates
}

func NewAppState() *fsm.FSM {
	states := ApplicationStates()
	return fsm.NewFSM(
		states.New, fsm.Events{
			{
				Name: SubmitApplication.String(),
				Src:  []string{states.New},
				Dst:  states.Submitted,
			},
			{
				Name: RecoverApplication.String(),
				Src:  []string{states.New},
				Dst:  states.Recovering,
			},
			{
				Name: AcceptApplication.String(),
				Src:  []string{states.Submitted, states.Recovering},
				Dst:  states.Accepted,
			},
			{
				Name: TryReserve.String(),
				Src:  []string{states.Accepted},
				Dst:  states.Reserving,
			},
			{
				Name: UpdateReservation.String(),
				Src:  []string{states.Reserving},
				Dst:  states.Reserving,
			},
			{
				Name: ResumingApplication.String(),
				Src:  []string{states.Reserving},
				Dst:  states.Resuming,
			},
			{
				Name: AppTaskCompleted.String(),
				Src:  []string{states.Resuming},
				Dst:  states.Resuming,
			},
			{
				Name: RunApplication.String(),
				Src:  []string{states.Accepted, states.Reserving, states.Resuming, states.Running},
				Dst:  states.Running,
			},
			{
				Name: ReleaseAppAllocation.String(),
				Src:  []string{states.Running},
				Dst:  states.Running,
			},
			{
				Name: ReleaseAppAllocation.String(),
				Src:  []string{states.Failing},
				Dst:  states.Failing,
			},
			{
				Name: ReleaseAppAllocation.String(),
				Src:  []string{states.Resuming},
				Dst:  states.Resuming,
			},
			{
				Name: ReleaseAppAllocationAsk.String(),
				Src:  []string{states.Running, states.Accepted, states.Reserving},
				Dst:  states.Running,
			},
			{
				Name: ReleaseAppAllocationAsk.String(),
				Src:  []string{states.Failing},
				Dst:  states.Failing,
			},
			{
				Name: ReleaseAppAllocationAsk.String(),
				Src:  []string{states.Resuming},
				Dst:  states.Resuming,
			},
			{
				Name: CompleteApplication.String(),
				Src:  []string{states.Running},
				Dst:  states.Completed,
			},
			{
				Name: RejectApplication.String(),
				Src:  []string{states.Submitted},
				Dst:  states.Rejected,
			},
			{
				Name: FailApplication.String(),
				Src:  []string{states.Submitted, states.Accepted, states.Running, states.Reserving},
				Dst:  states.Failing,
			},
			{
				Name: FailApplication.String(),
				Src:  []string{states.Failing, states.Rejected},
				Dst:  states.Failed,
			},
			{
				Name: KillApplication.String(),
				Src:  []string{states.Accepted, states.Running, states.Reserving},
				Dst:  states.Killing,
			},
			{
				Name: KilledApplication.String(),
				Src:  []string{states.Killing},
				Dst:  states.Killed,
			},
		},
		fsm.Callbacks{
			"enter_state": func(event *fsm.Event) {
				go func() {
					app := event.Args[0].(*Application) //nolint:errcheck
					log.Logger().Debug("shim app state transition",
						zap.String("app", app.GetApplicationID()),
						zap.String("source", event.Src),
						zap.String("destination", event.Dst),
						zap.String("event", event.Event))
				}()
			},
			states.Reserving: func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				app.OnReserving()
			},
			SubmitApplication.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				app.HandleSubmitApplicationEvent()
			},
			RecoverApplication.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				app.HandleRecoverApplicationEvent()
			},
			RejectApplication.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				eventInfo := event.Args[1].([]interface{})
				eventArgs := make([]string, 1)
				if err := events.GetEventArgsAsStrings(eventArgs, eventInfo); err != nil {
					log.Logger().Error("fail to parse event arg", zap.Error(err))
					return
				}
				reason := eventArgs[0]
				app.HandleRejectApplicationEvent(reason)
			},
			CompleteApplication.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				app.HandleCompleteApplicationEvent()
			},
			FailApplication.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				eventInfo := event.Args[1].([]interface{})
				eventArgs := make([]string, 1)
				if err := events.GetEventArgsAsStrings(eventArgs, eventInfo); err != nil {
					log.Logger().Error("fail to parse event arg", zap.Error(err))
					return
				}
				errMsg := eventArgs[0]
				app.HandleFailApplicationEvent(errMsg)
			},
			UpdateReservation.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				app.OnReservationStateChange()
			},
			ReleaseAppAllocation.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				eventInfo := event.Args[1].([]interface{})
				eventArgs := make([]string, 2)
				if err := events.GetEventArgsAsStrings(eventArgs, eventInfo); err != nil {
					log.Logger().Error("fail to parse event arg", zap.Error(err))
					return
				}
				allocUUID := eventArgs[0]
				terminationTypeStr := eventArgs[1]
				app.HandleReleaseAppAllocationEvent(allocUUID, terminationTypeStr)
			},
			ReleaseAppAllocationAsk.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				eventInfo := event.Args[1].([]interface{})
				eventArgs := make([]string, 2)
				if err := events.GetEventArgsAsStrings(eventArgs, eventInfo); err != nil {
					log.Logger().Error("fail to parse event arg", zap.Error(err))
					return
				}
				taskID := eventArgs[0]
				terminationTypeStr := eventArgs[1]
				app.HandleReleaseAppAllocationAskEvent(taskID, terminationTypeStr)
			},
			AppTaskCompleted.String(): func(event *fsm.Event) {
				app := event.Args[0].(*Application) //nolint:errcheck
				app.HandleAppTaskCompletedEvent()
			},
		},
	)
}
