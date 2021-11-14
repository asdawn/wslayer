/**
设备集，仅用于后台数据状态缓存和差分处理，不考虑消息格式问题

ver 1.0
*/
package device

import (
	"errors"
	"sync"
)

/**
设备状态清单
*/
type DeviceSet struct {
	DeviceClass string             //设备类
	Devices     map[string]*Device //ID->设备
	//Interval    int64                    //消息时间间隔
	LastModifyTime int64 //上一次消息的时间
	//OnTime      func(*DeviceSet)         //定时函数，
	//OnChange    func(*DeviceSet)
	RWLock sync.RWMutex
}

/**
新建一个DeviceSet对象，Now和Last相同
*/
func NewDeviceSet(DeviceClass string) *DeviceSet {
	return &DeviceSet{
		DeviceClass:    DeviceClass,
		Devices:        make(map[string]*Device),
		LastModifyTime: -1,
	}
}

/**
发生修改时的回调函数
*/
type OnChangeFunction func(*DeviceSet, *Device) error

/**
设置对象取值
device: 对象取值
onSetChange: deviceSet修改的回调函数，直接在新的go routine中执行，不占用当前线程
onDeviceMove: 对象变化的回调函数
thresholds: 判断是否变化的依据
返回（是否是新建对象，是否发生错误）
*/
func (deviceSet *DeviceSet) SetDevice(device *Device) (bool, error) {
	if device == nil {
		return false, errors.New("device should not be a null pinter")
	}
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	var err error = nil
	var id string = device.ID
	_, exists := deviceSet.Devices[id]
	deviceSet.Devices[id] = device
	//更新集合修改时间
	if device.T > deviceSet.LastModifyTime {
		deviceSet.LastModifyTime = device.T
	}

	return !exists, err
}

/**
设置对象取值
devices: 对象取值
onSetChange: deviceSet修改的回调函数，直接在新的go routine中执行，不占用当前线程
onDeviceMove: 对象变化的回调函数
thresholds: 判断是否变化的依据
返回（新建对象数，修改对象数，是否发生错误）
*/
func (deviceSet *DeviceSet) SetDevices(devices []*Device) (int, int, error) {
	if len(devices) == 0 {
		return 0, 0, nil
	}
	newCnt := 0
	modifyCnt := 0
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	var err error = nil
	for _, device := range devices {
		var id string = device.ID
		_, exists := deviceSet.Devices[id]
		if exists {
			modifyCnt++
		} else {
			newCnt++
		}
		deviceSet.Devices[id] = device
		//更新集合修改时间
		if device.T > deviceSet.LastModifyTime {
			deviceSet.LastModifyTime = device.T
		}
	}
	return newCnt, modifyCnt, err
}

/**
删除对象
device: 对象id
返回：（是否找到，是否出错）
*/
func (deviceSet *DeviceSet) RemoveDevice(id string) (bool, error) {
	if len(id) == 0 {
		return false, errors.New("len(id) should not be 0")
	}
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	var err error = nil
	_, exists := deviceSet.Devices[id]
	if exists {
		delete(deviceSet.Devices, id)
	}

	return !exists, err
}

/**
删除对象
devices: 对象id清单
返回：（删除成功个数，是否出错）
*/
func (deviceSet *DeviceSet) RemoveDevices(ids []string) (int, error) {
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	n := 0
	var err error = nil
	for _, id := range ids {
		_, exists := deviceSet.Devices[id]
		if exists {
			delete(deviceSet.Devices, id)
			n++
		}
	}
	return n, err
}

/**
获取指定ID的device数据
*/
func (deviceSet *DeviceSet) GetDevice(id string) *Device {
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	var device, ok = deviceSet.Devices[id]
	if !ok {
		return nil
	} else {
		return device
	}
}

/**
获取当前的Device数组
*/
func (deviceSet *DeviceSet) GetDevices() []*Device {
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	var devices []*Device
	devices = make([]*Device, 0)
	for _, device := range deviceSet.Devices {
		devices = append(devices, device)
	}
	return devices
}

/**
清空当前的设备数组
*/
func (deviceSet *DeviceSet) Clear() {
	deviceSet.RWLock.Lock()
	defer deviceSet.RWLock.Unlock()
	deviceSet.Devices = make(map[string]*Device)
}
