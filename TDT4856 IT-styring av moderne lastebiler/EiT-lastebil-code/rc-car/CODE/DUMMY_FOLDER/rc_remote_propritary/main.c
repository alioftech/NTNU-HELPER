/* Copyright (c) 2014 Nordic Semiconductor. All Rights Reserved.
 *
 * The information contained herein is property of Nordic Semiconductor ASA.
 * Terms and conditions of usage are described in detail in NORDIC
 * SEMICONDUCTOR STANDARD SOFTWARE LICENSE AGREEMENT.
 *
 * Licensees are granted free, non-transferable use of the information. NO
 * WARRANTY of ANY KIND is provided. This heading must NOT be removed from
 * the file.
 *
 */
#include "SEGGER_RTT.h"
#include <stdbool.h>
#include <stdint.h>
#include <string.h>
#include <stdio.h>

#include "sdk_common.h"
#include "nrf.h"
#include "nrf_esb.h"
#include "nrf_error.h"
#include "nrf_esb_error_codes.h"
#include "nrf_delay.h"
#include "nrf_gpio.h"
#include "boards.h"
#include "nrf_delay.h"
#include "app_util.h"
#define NRF_LOG_MODULE_NAME "REMOTE"
#include "nrf_log.h"
#include "nrf_log_ctrl.h"

#include "rc_saadc.h"
#include "rc_messages_and_defines.h"
#include "rc_utilities.h"

#define RTT_PRINTF(...) \
do { \
     char str[64];\
     sprintf(str, __VA_ARGS__);\
     SEGGER_RTT_WriteString(0, str);\
 } while(0)

#define printf RTT_PRINTF
 
#define TIMEOUT 500
#define RETIRES 6
#define DELAYTIME 300
 
static nrf_esb_payload_t        tx_payload = NRF_ESB_CREATE_PAYLOAD(0, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00);

static nrf_esb_payload_t        rx_payload;

/*lint -save -esym(40, BUTTON_1) -esym(40, BUTTON_2) -esym(40, BUTTON_3) -esym(40, BUTTON_4) -esym(40, LED_1) -esym(40, LED_2) -esym(40, LED_3) -esym(40, LED_4) */


static remote_packet_t remote_msg;
static car_packet_t car_msg;
static uint8_t STATE;
static uint8_t NEXT_STATE;
static uint8_t receive_ack;



void nrf_esb_event_handler(nrf_esb_evt_t const * p_event)
{
    switch (p_event->evt_id)
    {
        case NRF_ESB_EVENT_TX_SUCCESS:
            SEGGER_RTT_WriteString(0, "TX SUCCESS EVENT\n");
            break;
        case NRF_ESB_EVENT_TX_FAILED:
            SEGGER_RTT_WriteString(0, "TX FAILED EVENT\n");
            (void) nrf_esb_flush_tx();
            (void) nrf_esb_start_tx();
            break;
        case NRF_ESB_EVENT_RX_RECEIVED:
            SEGGER_RTT_WriteString(0, "RX RECEIVED EVENT\n");
            while (nrf_esb_read_rx_payload(&rx_payload) == NRF_SUCCESS)
            {
                if (rx_payload.length > 0 && extract_sender_id_from_payload(&rx_payload) == remote_msg.senderID)
                {
                    convert_payload_to_car_message(&car_msg, &rx_payload);
                    switch(STATE) {
                        case STATE_REMOTE_ADVERTISE_AVAILABLE :
                            if(car_msg.senderID == remote_msg.senderID && car_msg.type == MSG_CAR_TYPE_CONNECTED_TO){
                                receive_ack = 1;
                            }
                            break;
                        case STATE_REMOTE_SINGLE_MODE : 
                            if(car_msg.senderID == remote_msg.senderID && car_msg.type == MSG_CAR_TYPE_ACKNOWLEDGE){
                                receive_ack = 1;
                            }
                            break;
                        case STATE_REMOTE_TRUCK_POOLING_PENDING :
                            if(car_msg.senderID == remote_msg.senderID && car_msg.type == MSG_CAR_TYPE_ACKNOWLEDGE){
                                receive_ack = 1;
                            }
                            else if(car_msg.senderID == remote_msg.senderID && car_msg.type == MSG_CAR_TYPE_POOLING_SLAVE){
                                receive_ack = 1;
                            }
                            break;
                         case STATE_REMOTE_TRUCK_POOLING_ENABLED :
                            if(car_msg.senderID == remote_msg.senderID && car_msg.type == MSG_CAR_TYPE_POOLING_SLAVE){
                                receive_ack = 1;
                            }
                    }
                }
            }
            break;
    }
}


void clocks_start( void )
{
    NRF_CLOCK->EVENTS_HFCLKSTARTED = 0;
    NRF_CLOCK->TASKS_HFCLKSTART = 1;

    while (NRF_CLOCK->EVENTS_HFCLKSTARTED == 0);
}



uint32_t esb_init( void )
{
    uint32_t err_code;
    uint8_t base_addr_0[4] = {0xE7, 0xE7, 0xE7, 0xE7};
    uint8_t base_addr_1[4] = {0xC2, 0xC2, 0xC2, 0xC2};
    uint8_t addr_prefix[8] = {0xE7, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8 };

    nrf_esb_config_t nrf_esb_config         = NRF_ESB_DEFAULT_CONFIG;
    nrf_esb_config.protocol                 = NRF_ESB_PROTOCOL_ESB_DPL;
    nrf_esb_config.retransmit_delay         = 600;
    nrf_esb_config.bitrate                  = NRF_ESB_BITRATE_2MBPS;
    nrf_esb_config.event_handler            = nrf_esb_event_handler;
    nrf_esb_config.mode                     = NRF_ESB_MODE_PTX;
    nrf_esb_config.selective_auto_ack       = true;
    nrf_esb_config.payload_length           = 11;
   
    tx_payload.noack = true;

    err_code = nrf_esb_init(&nrf_esb_config);

    VERIFY_SUCCESS(err_code);

    err_code = nrf_esb_set_base_address_0(base_addr_0);
    VERIFY_SUCCESS(err_code);

    err_code = nrf_esb_set_base_address_1(base_addr_1);
    VERIFY_SUCCESS(err_code);

    err_code = nrf_esb_set_prefixes(addr_prefix, 8);
    VERIFY_SUCCESS(err_code);

    return err_code;
}


void radio_receive_mode(){
    nrf_esb_disable();
    nrf_delay_ms(1);
    esb_init();
    nrf_delay_ms(1);
    nrf_esb_start_rx();
    nrf_delay_ms(1);
}

void radio_transmit_mode(){
    nrf_delay_ms(1);
    nrf_esb_stop_rx();
    nrf_delay_ms(1);
    nrf_esb_disable();
    nrf_delay_ms(1);
    esb_init();
    nrf_delay_ms(1);
    nrf_esb_start_tx();
    nrf_delay_ms(1);
}

uint32_t radio_send_and_ack_message(uint32_t timeout){
    convert_remote_message_to_payload(&remote_msg, &tx_payload);
    printf("Wait for ack:\n");   
    receive_ack = 0;
    for(uint32_t i=0; i<RETIRES; i++){
        uint32_t wait = timeout;
        //printf("before transmit mode\n");
        //printf("after transmit mode\n");
        radio_transmit_mode();
        while(nrf_esb_write_payload(&tx_payload) != NRF_SUCCESS){
            //NOP
        }
        if (NRF_LOG_PROCESS() == false)
        {
         __WFE();
        }
        //printf("after write payload\n");
        radio_receive_mode();
        //printf("after receive mode\n");
        while(--wait > 0){
            if(receive_ack){
                break;
            }
            nrf_delay_ms(1);
        }
        if(receive_ack){
            //printf("Received ack\n");        
            break;
        }
        else{
            //printf("timeout occred, retry number: %d\n", i);
        }
    }
    return receive_ack;
}

int main(void)
{
    ret_code_t err_code;

    err_code = NRF_LOG_INIT(NULL);
    APP_ERROR_CHECK(err_code);

    clocks_start();

    err_code = esb_init();
    APP_ERROR_CHECK(err_code);

    leds_init();
    buttons_init();
    joystick_init();

    SEGGER_RTT_WriteString(0, "EiT Remote\r\n");
    STATE = STATE_REMOTE_WAIT_FOR_CHANNEL_SELECT;
    NEXT_STATE = STATE;
    remote_msg.x = 0;
    remote_msg.y = 0;
    nrf_esb_start_tx();
    while (true)
    {
        
        switch (STATE){
            case STATE_REMOTE_WAIT_FOR_CHANNEL_SELECT :
                printf("State: WAIT FOR CHANNEL\n"); 
                while(!get_pressed_button()){
                    // Wait until user presses button and selects sender ID
                }
                remote_msg.senderID = get_pressed_button();
								printf("SenderID is selected: %d\n", remote_msg.senderID);
                NEXT_STATE = STATE_REMOTE_ADVERTISE_AVAILABLE;
                break;
                
            case STATE_REMOTE_ADVERTISE_AVAILABLE :
                printf("State: ADVERTISE AVAILABLE\n");
								clear_led(remote_msg.senderID);
                printf("Advertising with SenderID: %d\n", remote_msg.senderID);
                remote_msg.type = MSG_REMOTE_TYPE_ADVERTISE_AVAILABLE;
                // Wait until remote receives car available
                uint32_t connection_establihed = 0;
                printf("Advertise");
                while(NEXT_STATE == STATE_REMOTE_ADVERTISE_AVAILABLE){
                    connection_establihed = radio_send_and_ack_message(100); // 100 ms between each resend
                    if(connection_establihed){
                        NEXT_STATE = STATE_REMOTE_SINGLE_MODE;
                        break;
                    }
                }
                set_led(remote_msg.senderID);
                printf("Connected to car\n");
                break;

            case STATE_REMOTE_SINGLE_MODE :
                //printf("State: SINGLE MODE\n");
                remote_msg.type    = MSG_REMOTE_TYPE_SINGLE_MODE_STEERING;
                remote_msg.x       = joystick_read(x_dir);
                remote_msg.y       = joystick_read(y_dir);
                remote_msg.button  = joystick_button_read();
                printf("X: %d\t Y: %d \t Button: %d\n", remote_msg.x, remote_msg.y, remote_msg.button);
                if(radio_send_and_ack_message(TIMEOUT)){
                    if(remote_msg.button){
                        NEXT_STATE = STATE_REMOTE_TRUCK_POOLING_PENDING;
                        break;
                    }
                }
                else // Connection lost, try to connect again first
                    NEXT_STATE = STATE_REMOTE_ADVERTISE_AVAILABLE;
                break;

            case STATE_REMOTE_TRUCK_POOLING_PENDING :
                printf("State: TRUCK POOLING PENDING\n");
                remote_msg.type    = MSG_REMOTE_TYPE_TRUCK_POOLING_REQUEST;
                remote_msg.x       = joystick_read(x_dir);
                remote_msg.y       = joystick_read(y_dir);
                remote_msg.button  = joystick_button_read();
                if(radio_send_and_ack_message(TIMEOUT)){
                    if(car_msg.type == MSG_CAR_TYPE_POOLING_SLAVE)
                        NEXT_STATE = STATE_REMOTE_TRUCK_POOLING_ENABLED;
                }
                else //Connection lost, try to connect again first
                    NEXT_STATE = STATE_REMOTE_ADVERTISE_AVAILABLE; 
                break;

            case STATE_REMOTE_TRUCK_POOLING_ENABLED :
                printf("State: TRUCK POOLING ENABLED\n");


                remote_msg.type    = MSG_REMOTE_TYPE_TRUCK_POOLING_SLAVE;
                remote_msg.x       = 0;
                remote_msg.y       = joystick_read(y_dir);
                remote_msg.button  = joystick_button_read();

                if(remote_msg.button){
                    remote_msg.type = MSG_REMOTE_TYPE_TRUCK_POOLING_STOP;
                    NEXT_STATE = STATE_REMOTE_SINGLE_MODE;
                }

                if(!radio_send_and_ack_message(TIMEOUT)){
                    NEXT_STATE = STATE_REMOTE_ADVERTISE_AVAILABLE;
                }
                break;
        }
        STATE = NEXT_STATE;
        nrf_delay_ms(DELAYTIME);
    }

}
/*lint -restore */
/** @} */
