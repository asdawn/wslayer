/**
用于生成要发送的地图对象控制指令

ver 1.0
*/

package message

import (
	"encoding/json"

	"../device"
)

/**
新建/更新对象
*/
type Object_management_message struct {
	CMDType      int              `json:"type"`
	ObjectClass  string           `json:"class"`
	ValuesSet    []*device.Device `json:"set"`
	ValuesDelete []string         `json:"del"`
	ValuesClear  bool             `json:"clear"`
}

/**
创建新建/更新对象消息
objectClass: 对象类名
values: 对象状态值数组，根据ID进行更新/创建
*/
func ObjectUpsertMessage(objectClass string, values []*device.Device) ([]byte, error) {
	var message = &Object_management_message{
		CMDType:     2,
		ObjectClass: objectClass,
		ValuesSet:   values,
	}
	return json.Marshal(message)
}

/**
创建对象操作空指令消息
*/
func ObjectNoopMessage(objectClass string) ([]byte, error) {
	var message = &Object_management_message{
		CMDType:     2,
		ObjectClass: objectClass,
	}
	return json.Marshal(message)
}

/**
创建清空对象消息
*/
func ObjectClearMessage(objectClass string) ([]byte, error) {
	var message = &Object_management_message{
		CMDType:     2,
		ObjectClass: objectClass,
		ValuesClear: true,
	}
	return json.Marshal(message)
}

/**
创建新建/更新对象消息
objectClass: 对象类名
ids: 要删除的对象ID数组
*/
func ObjectDeleteMessage(objectClass string, ids []string) ([]byte, error) {
	var message = &Object_management_message{
		CMDType:      2,
		ObjectClass:  objectClass,
		ValuesDelete: ids,
	}
	return json.Marshal(message)
}

/**
创建一类对象的全状态消息
deviceSet: 对象状体集
返回（消息，错误），如果deviceSet中无对象则消息返回nil而非空指令消息
*/
func ObjectFullStatusMessage(deviceSet *device.DeviceSet) ([]byte, error) {
	/*
		if deviceSet == nil {
			return nil, errors.New("deviceset should not be null")
		}
	*/
	objectClass := deviceSet.DeviceClass
	if len(deviceSet.Devices) == 0 {
		return ObjectNoopMessage(objectClass)
	} else {
		devices := deviceSet.GetDevices()
		return ObjectUpsertMessage(objectClass, devices)
	}
}
