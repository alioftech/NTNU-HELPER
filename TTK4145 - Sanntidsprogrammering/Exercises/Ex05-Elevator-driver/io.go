
package Driver

import "C"

func io_init() bool{ 
	C.io_init()
}

func io_set_bit(channel int) {
	C.io_set_bit(C.int(channel)
}

func io_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel)
}

func io_read_bit(channel int) int{
	return int( C.io_read_bit(C.int(channel)))
}

func io_read_analog(channel int) int{
	return int(C.io_read_analog(C.int(channel)))
}

func io_write_analog(channel, value int){
	C.io_write_analog(C.int(channel), C.int(value))
}