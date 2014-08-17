godmi-utils
===========

#简介
 在实现godmi的过程中，发现大部分的工作都是：定义一个类型的结构体，然后定义结构体构造方法，String方法。其实后面的两项工作通过解析结构体的定义都可以自动化。
 
#例子：
 [root@tod godmi-utils]# cat template.go
```go
type PortableBatteryDeviceChemistry byte

type PortableBattery struct {
        InfoCommon
        Location                  string
        Manufacturer              string
        ManufacturerDate          string
        SerialNumber              string
        DeviceName                string
        DeviceChemistry           PortableBatteryDeviceChemistry
        DesignCapacity            uint16
        DesignVoltage             uint16
        SBDSVersionNumber         string
        MaximumErrorInBatteryData byte
        SBDSSerialNumber          uint16
        SBDSManufactureDate       uint16
        SBDSDeviceChemistry       string
        DesignCapacityMultiplier  byte
        OEMSepecific              uint32
}
```
如上所示，定义了一个结构体，其中InfoCommon的定义放在godmi-gentype.go里面了，上面的定义放在另外一个文件template.go
那么，通过godmi-gentype.go, 生成以下函数。

 [root@tod godmi-utils]# go run godmi-gentype.go -template ./template.go -typename="PortableBattery"
```go
func (p PortableBattery) String() string {
        return fmt.Sprintf("Portable Battery\n"+
                "\t\tLocation: %s\n"+
                "\t\tManufacturer: %s\n"+
                "\t\tManufacturer Date: %s\n"+
                "\t\tSerial Number: %s\n"+
                "\t\tDevice Name: %s\n"+
                "\t\tDevice Chemistry: %s\n"+
                "\t\tDesign Capacity: %d\n"+
                "\t\tDesign Voltage: %d\n"+
                "\t\tSBDS Version Number: %s\n"+
                "\t\tMaximum Error In Battery Data: %d\n"+
                "\t\tSBDS Serial Number: %d\n"+
                "\t\tSBDS Manufacture Date: %d\n"+
                "\t\tSBDS Device Chemistry: %s\n"+
                "\t\tDesign Capacity Multiplier: %d\n"+
                "\t\tOEM Sepecific: %d",
                p.Location,
                p.Manufacturer,
                p.ManufacturerDate,
                p.SerialNumber,
                p.DeviceName,
                p.DeviceChemistry,
                p.DesignCapacity,
                p.DesignVoltage,
                p.SBDSVersionNumber,
                p.MaximumErrorInBatteryData,
                p.SBDSSerialNumber,
                p.SBDSManufactureDate,
                p.SBDSDeviceChemistry,
                p.DesignCapacityMultiplier,
                p.OEMSepecific,
        )
}

func (h DMIHeader) PortableBattery() *PortableBattery {
	data := h.data
	return &PortableBattery{
		Location:                  h.FieldString(int(data[0x04])),
		Manufacturer:              h.FieldString(int(data[0x05])),
		ManufacturerDate:          h.FieldString(int(data[0x06])),
		SerialNumber:              h.FieldString(int(data[0x07])),
		DeviceName:                h.FieldString(int(data[0x08])),
		DeviceChemistry:           PortableBatteryDeviceChemistry(data[0x09]),
		DesignCapacity:            u16(data[0x0A:0x0C]),
		DesignVoltage:             u16(data[0x0C:0x0E]),
		SBDSVersionNumber:         h.FieldString(int(data[0x0E])),
		MaximumErrorInBatteryData: data[0x0F],
		SBDSSerialNumber:          u16(data[0x10:0x12]),
		SBDSManufactureDate:       u16(data[0x12:0x14]),
		SBDSDeviceChemistry:       h.FieldString(int(data[0x14])),
		DesignCapacityMultiplier:  data[0x15],
		OEMSepecific:              u32(data[0x16:0x1A]),
	}
}
```
