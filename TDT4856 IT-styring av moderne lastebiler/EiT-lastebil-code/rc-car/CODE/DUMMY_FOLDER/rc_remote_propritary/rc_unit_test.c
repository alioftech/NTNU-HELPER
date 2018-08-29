#include <stdint.h>
#include <stdio.h>
#include <string.h>

#include "nrf_delay.h"


#include "rc_saadc.h"

#include "nrf_log.h"
#include "nrf_log_ctrl.h"

void unit_test_saadc(){
    //static nrf_saadc_value_t     adc_val[2][1];
    saadc_enable(); 
 
    for(uint32_t i=0; i<100; i++){
        uint32_t x_val = joystick_read(x_dir);
        uint32_t y_val = joystick_read(y_dir);
        NRF_LOG_INFO("X: %d\t Y: %d\t \n\r", x_val, y_val);
        NRF_LOG_FLUSH();
        nrf_delay_ms(100);
    }
}



/** @} */
