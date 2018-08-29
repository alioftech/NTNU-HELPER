#pragma once

//io_init returns no-zero on sucess and 0 on failiure
int io_init();

void io_set_bit(int channel);

void io_clear_bit(int channel);

void io_write_analog(int channel, int value);

int io_read_bit(int channel);

int io_read_analog(int channel);
