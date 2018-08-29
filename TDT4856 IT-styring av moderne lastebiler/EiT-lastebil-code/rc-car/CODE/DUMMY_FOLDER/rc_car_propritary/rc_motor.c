#include <stdio.h>
#include <string.h>
#include "app_util_platform.h"
#include "app_error.h"

#include "nrf_gpio.h"
#include "nrf_drv_pwm.h"
#include "nrf_delay.h"

#include "rc_motor.h"


static nrf_drv_pwm_t m_pwm0 = NRF_DRV_PWM_INSTANCE(1);

// This array cannot be allocated on stack (hence "static") and it must
// be in RAM (hence no "const", though its content is not changed).
static nrf_pwm_values_individual_t m_demo1_seq_values;
static nrf_pwm_sequence_t const    m_demo1_seq =
{
    .values.p_individual = &m_demo1_seq_values,
    .length              = NRF_PWM_VALUES_LENGTH(m_demo1_seq_values),
    .repeats             = 0,
    .end_delay           = 0
};

void motor_init(){
    //set motor direction pins as outputs
   // nrf_gpio_pin_dir_set(motor_pwm_a, NRF_GPIO_PIN_DIR_OUTPUT);
    //nrf_gpio_pin_dir_set(motor_pwm_b, NRF_GPIO_PIN_DIR_OUTPUT); 
    nrf_gpio_pin_dir_set(motor_in_1_a, NRF_GPIO_PIN_DIR_OUTPUT); 
    nrf_gpio_pin_dir_set(motor_in_2_a, NRF_GPIO_PIN_DIR_OUTPUT);
    nrf_gpio_pin_dir_set(motor_in_3_b, NRF_GPIO_PIN_DIR_OUTPUT);
    nrf_gpio_pin_dir_set(motor_in_4_b, NRF_GPIO_PIN_DIR_OUTPUT);

    //force no movement
    //nrf_gpio_pin_clear(motor_pwm_a);
    //nrf_gpio_pin_clear(motor_pwm_b);
    nrf_gpio_pin_clear(motor_in_1_a);
    nrf_gpio_pin_clear(motor_in_2_a);
    nrf_gpio_pin_clear(motor_in_3_b);
    nrf_gpio_pin_clear(motor_in_4_b);

    //local pwm setup
    uint32_t                   err_code;
    nrf_drv_pwm_config_t const config0 =
    {
        .output_pins =
        {
            motor_pwm_a, // channel 0
            NRF_DRV_PWM_PIN_NOT_USED, // channel 2
            motor_pwm_b, // channel 1
            NRF_DRV_PWM_PIN_NOT_USED  // channel 3
        },
        .irq_priority = APP_IRQ_PRIORITY_LOWEST, ///< Interrupt priority.
        .base_clock   = NRF_PWM_CLK_125kHz,      ///< Base clock frequency.
        .count_mode   = NRF_PWM_MODE_UP,         ///< Operating mode of the pulse generator counter.
        .top_value    = PWM_TOP_VALUE,           //< Value up to which the pulse generator counter counts.
        .load_mode    = NRF_PWM_LOAD_INDIVIDUAL,     ///< Mode of loading sequence data from RAM.
        .step_mode    = NRF_PWM_STEP_AUTO        ///< Mode of advancing the active sequence.
    };

    //pwm channel values, pwm is inverted!!!!
    m_demo1_seq_values.channel_0 = PWM_TOP_VALUE;
    m_demo1_seq_values.channel_1 = PWM_TOP_VALUE;
    m_demo1_seq_values.channel_2 = PWM_TOP_VALUE;
    m_demo1_seq_values.channel_3 = PWM_TOP_VALUE;
    
    //init pwm
    err_code = nrf_drv_pwm_init(&m_pwm0, &config0, NULL);
    APP_ERROR_CHECK(err_code);
    //set defaulult direction forward
    motor_set_dir(BOTH, FORWARD);
    motor_set_speed(BOTH, 0);
}

void motor_set_dir(uint32_t side, uint32_t dir){
    //need to disable both enable pins before direction change
    //motor_stop();
    //nrf_delay_ms(1);

    uint32_t motor_1;
    uint32_t motor_2;
    uint32_t num_of_sides;
    switch (side){
        case LEFT :
            num_of_sides = 1;
            motor_1 = motor_in_4_b;
            motor_2 = motor_in_3_b;
            break;
        case RIGHT :
            num_of_sides = 1;
            motor_1 = motor_in_1_a;
            motor_2 = motor_in_2_a;
            break;
        case BOTH :
            num_of_sides = 2;
            motor_1 = motor_in_1_a;
            motor_2 = motor_in_2_a;
    }
    for (uint32_t i=0; i<num_of_sides; i++){
        switch (dir) {
                case FORWARD :
                    nrf_gpio_pin_clear(motor_2);
                    nrf_delay_ms(1);
                    nrf_gpio_pin_set(motor_1);
                    break;
                case BACKWARD :
                    nrf_gpio_pin_clear(motor_1);
                    nrf_delay_ms(1);
                    nrf_gpio_pin_set(motor_2);
                    break;
        }
        //Change motors for next iteration
        if(side == BOTH){
            motor_1 = motor_in_4_b;
            motor_2 = motor_in_3_b;
        }   //end if
    } //end for
    //motor_start();
}
void motor_set_speed(uint32_t side, uint32_t speed){
    if (side == BOTH){
        m_demo1_seq_values.channel_0 = PWM_TOP_VALUE - speed;
        m_demo1_seq_values.channel_2 = PWM_TOP_VALUE - speed;
    }

    else if (side == RIGHT)
        m_demo1_seq_values.channel_0 = PWM_TOP_VALUE - speed;
    else if (side == LEFT)
        m_demo1_seq_values.channel_2 = PWM_TOP_VALUE - speed;
    
}

void motor_start(){
   nrf_drv_pwm_simple_playback(&m_pwm0, &m_demo1_seq, 1,
                                NRF_DRV_PWM_FLAG_LOOP);
 }

void motor_stop(){
    nrf_drv_pwm_stop(&m_pwm0, true); //stop pwm and wait until it is stopped, henvce true
}


void set_motors(uint32_t left_speed, uint32_t right_speed, uint32_t left_dir, uint32_t right_dir){
    if(left_dir == right_dir){
        motor_set_dir(BOTH, left_dir);
    }else{
        motor_set_dir(LEFT, left_dir);
        motor_set_dir(RIGHT, right_dir);
    }

    if(left_speed == right_speed){
        motor_set_speed(BOTH, left_speed);
    }else{
        motor_set_speed(LEFT, left_speed);
        motor_set_speed(RIGHT, right_speed);
    }
}

/** @} */

