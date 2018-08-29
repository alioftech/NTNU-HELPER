#ifndef __RC_UTILITIES_H
#define __RC_UTILITIES_H

#include <stdbool.h>
#include <stdint.h>
#include "nrf_esb.h"

#include "rc_messages_and_defines.h"

void leds_init(void);
void set_led(uint32_t number);
void clear_led(uint32_t number);

void buttons_init(void);
uint8_t get_pressed_button(void);


void convert_car_message_to_payload(car_packet_t const * car_msg, nrf_esb_payload_t * p_payload);
void convert_payload_to_car_message(car_packet_t  * car_msg, nrf_esb_payload_t const * p_payload);

void convert_remote_message_to_payload(remote_packet_t const * remote_msg, nrf_esb_payload_t * p_payload);
void convert_payload_to_remote_message(remote_packet_t  * remote_msg, nrf_esb_payload_t const * p_payload);

void convert_master_message_to_payload(master_packet_t const * master_msg, nrf_esb_payload_t * p_payload);
void convert_payload_to_master_message(master_packet_t *master_msg, nrf_esb_payload_t const *p_payload);

uint32_t extract_sender_id_from_payload(nrf_esb_payload_t const *p_payload);
uint32_t extract_type_from_payload(nrf_esb_payload_t const *p_payload);

#endif // UTILITIES
/** @} */
