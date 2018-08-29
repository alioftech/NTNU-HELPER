#include "SEGGER_RTT.h"
#include <stdbool.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>

#include "nrf_delay.h"
#include "nrf.h"

#include "bsp.h"

#include "nrf_esb.h"
#include "rc_messages_and_defines.h"
#include "rc_utilities.h"


#define RTT_PRINTF(...) \
do { \
     char str[64];\
     sprintf(str, __VA_ARGS__);\
     SEGGER_RTT_WriteString(0, str);\
 } while(0)

#define printf RTT_PRINTF


void leds_init(){
    nrf_gpio_pin_dir_set(17, NRF_GPIO_PIN_DIR_OUTPUT); 
    nrf_gpio_pin_dir_set(18, NRF_GPIO_PIN_DIR_OUTPUT); 
    nrf_gpio_pin_dir_set(19, NRF_GPIO_PIN_DIR_OUTPUT); 
    nrf_gpio_pin_dir_set(20, NRF_GPIO_PIN_DIR_OUTPUT); 
    nrf_gpio_pin_set(17);
    nrf_gpio_pin_set(18);
    nrf_gpio_pin_set(19);
    nrf_gpio_pin_set(20);
}

void buttons_init(){
    nrf_gpio_cfg_input(13, NRF_GPIO_PIN_PULLUP);
    nrf_gpio_cfg_input(14, NRF_GPIO_PIN_PULLUP);
    nrf_gpio_cfg_input(15, NRF_GPIO_PIN_PULLUP);
    nrf_gpio_cfg_input(16, NRF_GPIO_PIN_PULLUP);
}

void set_led(uint32_t number){
    if(number >=1 && number <= 4){
        number = 16 + number; // Convert from led number to pin    
        nrf_gpio_pin_clear(number);
    }
}

void clear_led(uint32_t number){
    if(number >=1 && number <= 4){
        number = 16 + number; // Convert from led number to pin    
        nrf_gpio_pin_set(number);
    }
}

uint8_t get_pressed_button(){
    // Priority button 1 first, button 4 last
    if(!nrf_gpio_pin_read(13))
        return 1;
    if(!nrf_gpio_pin_read(14))
        return 2;
    if(!nrf_gpio_pin_read(15))
        return 3;
    if(!nrf_gpio_pin_read(16))
        return 4;
// Retrun 0 if no buttons are pressed
    return 0;
}

void convert_car_message_to_payload(car_packet_t const * car_msg, nrf_esb_payload_t * p_payload){
    p_payload->data[0]   = car_msg->senderID;
    p_payload->data[1]   = car_msg->type;   
}
void convert_payload_to_car_message(car_packet_t * car_msg, nrf_esb_payload_t const * p_payload){

    // Make message struct
    car_msg->senderID   = p_payload->data[0];
    car_msg->type       = p_payload->data[1];
}


void convert_remote_message_to_payload(remote_packet_t const * remote_msg, nrf_esb_payload_t * p_payload){
    p_payload->data[0]   = remote_msg->senderID;
    p_payload->data[1]   = remote_msg->type;
    p_payload->data[2]   = remote_msg->x;       // LSB
    p_payload->data[3]   = remote_msg->x >> 8; 
    p_payload->data[4]   = remote_msg->x >> 16;
    p_payload->data[5]   = remote_msg->x >> 24; // MSB
    p_payload->data[6]   = remote_msg->y;       // LSB 
    p_payload->data[7]   = remote_msg->y >> 8;  
    p_payload->data[8]   = remote_msg->y >> 16; 
    p_payload->data[9]   = remote_msg->y >> 24; // MSB
    p_payload->data[10]  = remote_msg->button;    
}
void convert_payload_to_remote_message(remote_packet_t * remote_msg, nrf_esb_payload_t const * p_payload){
    
    
    // Convert 8 bit array values to 32 bit single value
    uint8_t joy_val_8_x[4] = {0};
    uint8_t joy_val_8_y[4] = {0};
    for(uint32_t i = 0; i < 4; i++){
        joy_val_8_x[i] = p_payload->data[i+2];
        joy_val_8_y[i] = p_payload->data[i+6];
    }
    // Make message struct
    remote_msg->senderID = p_payload->data[0];
    remote_msg->type     = p_payload->data[1];
    remote_msg->x        = *((uint32_t*) joy_val_8_x);
    remote_msg->y        = *((uint32_t*) joy_val_8_y);
    remote_msg->button   = p_payload->data[10];
}

void convert_master_message_to_payload(master_packet_t const * master_msg, nrf_esb_payload_t * p_payload){
    uint8_t speed_info_converted;
    if (master_msg->speed_info > 100){
        speed_info_converted = 200;
    }
    else if(master_msg->speed_info > 0){
        speed_info_converted = master_msg->speed_info + 100;
    }
    else if(master_msg->speed_info < -100){
        speed_info_converted = 0;
    }
    else if(master_msg->speed_info < 0){
        speed_info_converted = master_msg->speed_info + 100;
    }

    p_payload->data[0]   = master_msg->senderID;
    p_payload->data[1]   = master_msg->type;
    p_payload->data[2]   = speed_info_converted;   
}
void convert_payload_to_master_message(master_packet_t * master_msg, nrf_esb_payload_t const * p_payload){

    // Make message struct
    master_msg->senderID   = p_payload->data[0];
    master_msg->type       = p_payload->data[1];
    master_msg->speed_info = p_payload->data[2] - 100;
}


uint32_t extract_sender_id_from_payload(nrf_esb_payload_t const *p_payload){
    // Convert 8 bit array values to 32 bit single value
    uint8_t joy_val_8_x[4] = {0};
    uint8_t joy_val_8_y[4] = {0};
    for(uint32_t i = 0; i < 4; i++){
        joy_val_8_x[i] = p_payload->data[i+2];
        joy_val_8_y[i] = p_payload->data[i+6];
    }
    return p_payload->data[0];
}
uint32_t extract_type_from_payload(nrf_esb_payload_t const *p_payload){
    // Convert 8 bit array values to 32 bit single value
    uint8_t joy_val_8_x[4] = {0};
    uint8_t joy_val_8_y[4] = {0};
    for(uint32_t i = 0; i < 4; i++){
        joy_val_8_x[i] = p_payload->data[i+2];
        joy_val_8_y[i] = p_payload->data[i+6];
    }
    return p_payload->data[1];
}





/** @} */
