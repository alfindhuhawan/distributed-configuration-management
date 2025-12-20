package logger

import (
	"context"
	"distributed-configuration-management/shared/driver"
	"fmt"
	"time"
)

func NewSimpleJSONLogger(appData driver.ApplicationData) Logger {
	return &simpleJSONLoggerImpl{AppData: appData}
}

type JSONLogModel struct {
	AppName   string    `json:"appName" bson:"app_name"`
	AppInstID string    `json:"appInstID" bson:"app_inst_id"`
	Start     string    `json:"start" bson:"start"`
	Severity  string    `json:"severity" bson:"severity"`
	Message   string    `json:"message" bson:"message"`
	Location  string    `json:"location" bson:"location"`
	Time      time.Time `json:"time" bson:"time"`
	TraceID   string    `json:"traceID" bson:"trace_id"`
}

func newJSONLogModel(lg *simpleJSONLoggerImpl, flag, loc string, msg any, trid string) JSONLogModel {

	message := fmt.Sprintf("%v %v", trid, msg)

	if flag == "ERROR" {
		message = fmt.Sprintf("%v %v %v", trid, loc, msg)
	}

	return JSONLogModel{
		AppName:   lg.AppData.AppName,
		AppInstID: lg.AppData.AppInstanceID,
		Start:     lg.AppData.StartTime,
		Severity:  flag,
		Message:   message,
		Location:  loc,
		Time:      time.Now(),
		TraceID:   trid,
	}
}

type simpleJSONLoggerImpl struct {
	AppData driver.ApplicationData
}

func (l simpleJSONLoggerImpl) Info(ctx context.Context, message string, args ...any) {
	messageWithArgs := fmt.Sprintf(message, args...)
	l.printLog(ctx, "INFO", messageWithArgs)
}

func (l simpleJSONLoggerImpl) Error(ctx context.Context, message string, args ...any) {
	messageWithArgs := fmt.Sprintf(message, args...)
	l.printLog(ctx, "ERROR", messageWithArgs)
}

func (l simpleJSONLoggerImpl) printLog(ctx context.Context, flag string, data any) {
	traceID := GetTraceID(ctx)
	fmt.Printf("%-5s %s %-60v %s\n", flag, traceID, data, getFileLocationInfo(3))

	// err := l.Pub.Publish(event.AllLogEvent, 3, payload.Payload{
	// 	Data:      newJSONLogModel(&l, flag, getFileLocationInfo(3), data, traceID),
	// 	Publisher: l.AppData,
	// 	TraceID:   traceID,
	// })
	// if err != nil {
	// 	fmt.Printf("ERROR >>>> %v\n", err.Error())
	// 	return
	// }

	//fmt.Println(newJSONLogModel(&l, flag, getFileLocationInfo(3), data, traceID))
}
