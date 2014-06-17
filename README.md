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
        return fmt.Sprintf("Portable Battery:\n\t\t"+
                "Location: %s\n\t\t"+
                "Manufacturer: %s\n\t\t"+
                "Manufacturer Date: %s\n\t\t"+
                "Serial Number: %s\n\t\t"+
                "Device Name: %s\n\t\t"+
                "Device Chemistry: %s\n\t\t"+
                "Design Capacity: %d\n\t\t"+
                "Design Voltage: %d\n\t\t"+
                "SBDS Version Number: %s\n\t\t"+
                "Maximum Error In Battery Data: %d\n\t\t"+
                "SBDS Serial Number: %d\n\t\t"+
                "SBDS Manufacture Date: %d\n\t\t"+
                "SBDS Device Chemistry: %s\n\t\t"+
                "Design Capacity Multiplier: %d\n\t\t"+
                "OEM Sepecific: %d\n",
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

func (h DMIHeader) PortableBattery() PortableBattery {
        var p PortableBattery
        data := h.data
        p.Location = h.FieldString(int(data[0x04]))
        p.Manufacturer = h.FieldString(int(data[0x05]))
        p.ManufacturerDate = h.FieldString(int(data[0x06]))
        p.SerialNumber = h.FieldString(int(data[0x07]))
        p.DeviceName = h.FieldString(int(data[0x08]))
        p.DeviceChemistry = PortableBatteryDeviceChemistry(data[0x09])
        p.DesignCapacity = U16(data[0x0A:0x0C])
        p.DesignVoltage = U16(data[0x0C:0x0E])
        p.SBDSVersionNumber = h.FieldString(int(data[0x0E]))
        p.MaximumErrorInBatteryData = data[0x0F]
        p.SBDSSerialNumber = U16(data[0x10:0x12])
        p.SBDSManufactureDate = U16(data[0x12:0x14])
        p.SBDSDeviceChemistry = h.FieldString(int(data[0x14]))
        p.DesignCapacityMultiplier = data[0x15]
        p.OEMSepecific = U32(data[0x16:0x1A])
        return p
}
```
