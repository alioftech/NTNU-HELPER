#ifndef RC_SAADC__
#define RC_SAADC__


// NB X and Y pins must be either 3 or 4, else change the c file!!
#include <stdint.h>
#include "nrf_drv_saadc.h"

#define x_pin NRF_SAADC_INPUT_AIN1
#define y_pin NRF_SAADC_INPUT_AIN2
#define x_dir 0
#define y_dir 1
#define joystick_button_pin 28

void saadc_init(void);
void saadc_enable(void);

void joystick_button_init(void);
void joystick_init(void);

uint32_t joystick_read(uint32_t direction);
uint32_t joystick_button_read(void);

#endif // RC_SAADC__




/** @} */
