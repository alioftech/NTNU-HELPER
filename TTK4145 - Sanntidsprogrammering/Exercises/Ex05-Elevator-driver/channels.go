package Driver


//in port 4
const PORT_4_SUBDEVICE        = 3
const PORT_4_CHANNEL_OFFSET   = 16
const PORT_4_DIRECTION        = COMEDI_INPUT
const OBSTRUCTION             = (0x300+23)
const STOP                    = (0x300+22)
const BUTTON_COMMAND1         = (0x300+21)
const BUTTON_COMMAND2         = (0x300+20)
const BUTTON_COMMAND3         = (0x300+19)
const BUTTON_COMMAND4         = (0x300+18)
const BUTTON_UP1              = (0x300+17)
const BUTTON_UP2              = (0x300+16)

//in port 1
const PORT_1_SUBDEVICE        = 2
const PORT_1_CHANNEL_OFFSET   = 0
const PORT_1_DIRECTION        = COMEDI_INPUT
const BUTTON_DOWN2            = (0x200+0)
const BUTTON_UP3              = (0x200+1)
const BUTTON_DOWN3            = (0x200+2)
const BUTTON_DOWN4            = (0x200+3)
const SENSOR_FLOOR1           = (0x200+4)
const SENSOR_FLOOR2           = (0x200+5)
const SENSOR_FLOOR3           = (0x200+6)
const SENSOR_FLOOR4           = (0x200+7)

//out port 3
const PORT_3_SUBDEVICE        =3
const PORT_3_CHANNEL_OFFSET   =8
const PORT_3_DIRECTION        =COMEDI_OUTPUT
const MOTORDIR                = (0x300+15)
const LIGHT_STOP              = (0x300+14)
const LIGHT_COMMAND1          = (0x300+13)
const LIGHT_COMMAND2          = (0x300+12)
const LIGHT_COMMAND3          = (0x300+11)
const LIGHT_COMMAND4          = (0x300+10)
const LIGHT_UP1               = (0x300+9)
const LIGHT_UP2               = (0x300+8)

//out port 2
const PORT_2_SUBDEVICE        = 3
const PORT_2_CHANNEL_OFFSET   = 0
const PORT_2_DIRECTION        = COMEDI_OUTPUT
const LIGHT_DOWN2             = (0x300+7)
const LIGHT_UP3               = (0x300+6)
const LIGHT_DOWN3             = (0x300+5)
const LIGHT_DOWN4             = (0x300+4)
const LIGHT_DOOR_OPEN         = (0x300+3)
const LIGHT_FLOOR_IND2        = (0x300+1)
const LIGHT_FLOOR_IND1        = (0x300+0)

//out port 0
const MOTOR                   = (0x100+0)

//non-existing ports (for alignment)
const BUTTON_DOWN1            = -1
const BUTTON_UP4              = -1
const LIGHT_DOWN1             = -1
const LIGHT_UP4               = -1