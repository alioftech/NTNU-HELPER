#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "nrf.h"
#include "nrf_delay.h"
#include "nrf_drv_timer.h"
#include "nrf_gpio.h"
#include "bsp.h"

#include "rc_ultrasound.h"

#define read_from_reg(addr) (*(addr))

void ultrasound_init(){
	nrf_gpio_pin_dir_set(ultrasound_trig, NRF_GPIO_PIN_DIR_OUTPUT);
	nrf_gpio_pin_dir_set(ultrasound_echo, NRF_GPIO_PIN_DIR_INPUT);
	nrf_gpio_pin_clear(ultrasound_trig);
	// Define counter top value
    *SYSTICK_RVR = SYSTICK_TOP-1;
}

uint32_t ultrasound_get_distance(){
    // Send trigger to ultrasound device
    nrf_gpio_pin_set(ultrasound_trig);
    nrf_delay_us(10);
    nrf_gpio_pin_clear(ultrasound_trig);
    // Wait until pulse on echo pin
    // Clear current value register, any write value clear register
    *SYSTICK_CVR = 0;
    // Start counter, timer counts from TOP and downwards
    *SYSTICK_CSR |= (1 << 0);
    while (nrf_gpio_pin_read(ultrasound_echo) == 0);
    // Get start counter value after entire puls is transmitted
    uint32_t pulse_start_time =  read_from_reg(SYSTICK_CVR);
    while (nrf_gpio_pin_read(ultrasound_echo) == 1);
    // Get echo conter value after the entire puls is received
    uint32_t pulse_stop_time = read_from_reg(SYSTICK_CVR);
    // Disable counter
    *SYSTICK_CSR &= ~(1 << 0);
    //printf("Start: %d\tStop: %d\tDiff: %d\t \r\n", pulse_start_time, pulse_stop_time, pulse_start_time-pulse_stop_time);
    // Calculate distance in mm
    uint32_t pulse_in_us = (pulse_start_time-pulse_stop_time)/64;
    // Mult by 10 to get mm, divide by 58 gives cm
    return pulse_in_us * 10 / (58);
}

/** @} */
