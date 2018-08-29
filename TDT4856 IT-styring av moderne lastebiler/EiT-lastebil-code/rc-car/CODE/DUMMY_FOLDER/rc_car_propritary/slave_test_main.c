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
#define NRF_LOG_MODULE_NAME "CAR"
#include "nrf_log.h"
#include "nrf_log_ctrl.h"

// From car
#include "rc_motor.h"
#include "rc_ultrasound.h"
#include "rc_steering.h"
#include "rc_controller.h"
#include "rc_filter.h"
#include "rc_unit_test.h"

// From framework
#include "rc_messages_and_defines.h"
#include "rc_utilities.h"

#define max(X, Y) (((X) > (Y)) ? (X) : (Y))
#define min(X, Y) (((X) < (Y)) ? (X) : (Y))


#define RTT_PRINTF(...) \
do { \
     char str[64];\
     sprintf(str, __VA_ARGS__);\
     SEGGER_RTT_WriteString(0, str);\
 } while(0)

#define printf RTT_PRINTF

#define TIMEOUT 100
#define RETIRES 3

static nrf_esb_payload_t        tx_payload = NRF_ESB_CREATE_PAYLOAD(0, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00);

static nrf_esb_payload_t        rx_payload;

/*lint -save -esym(40, BUTTON_1) -esym(40, BUTTON_2) -esym(40, BUTTON_3) -esym(40, BUTTON_4) -esym(40, LED_1) -esym(40, LED_2) -esym(40, LED_3) -esym(40, LED_4) */


static remote_packet_t remote_msg;
static car_packet_t car_msg;
static master_packet_t master_msg;
static uint8_t STATE;
static uint8_t NEXT_STATE;
static uint8_t receive_ack;



void nrf_esb_event_handler(nrf_esb_evt_t const * p_event)
{
    return;
    switch (p_event->evt_id)
    {
        case NRF_ESB_EVENT_TX_SUCCESS:
            SEGGER_RTT_WriteString(0, "CAR::::: TX SUCCESS EVENT\r\n");
            break;
        case NRF_ESB_EVENT_TX_FAILED:
            SEGGER_RTT_WriteString(0, "CAR::::: TX FAILED EVENT\r\n");
            (void) nrf_esb_flush_tx();
            (void) nrf_esb_start_tx();
            break;
        case NRF_ESB_EVENT_RX_RECEIVED:
						printf("%s\n","RX Event");
            nrf_delay_ms(2);
            if (nrf_esb_read_rx_payload(&rx_payload) == NRF_SUCCESS)
            {
                if (rx_payload.length > 0)
                {
                    nrf_delay_ms(2);   
										
                    printf("RX with payload\n");
                    nrf_delay_ms(2);                                   

                    if(extract_sender_id_from_payload(&rx_payload) == car_msg.senderID){
                        convert_payload_to_remote_message(&remote_msg, &rx_payload);
                        switch(remote_msg.type){
                            case MSG_REMOTE_TYPE_JOYSTICK:
                                break;

                            case MSG_REMOTE_TYPE_REQUEST_POOLING:
                                if(STATE == STATE_CAR_SINGLE_MODE){
                                    NEXT_STATE = STATE_CAR_TRUCK_POOLING_PENDING;
                                }
                                break;

                            case MSG_REMOTE_TYPE_STOP_POOLING:
                                if(STATE == (STATE_CAR_TRUCK_POOLING_SLAVE || STATE_CAR_TRUCK_POOLING_PENDING)){
                                    NEXT_STATE = STATE_CAR_SINGLE_MODE;
                                }
                                break;
                        }
                    }else{
                        switch(extract_type_from_payload(&rx_payload)){
                            case MSG_REMOTE_TYPE_ADVERTISE_AVAILABLE:
																nrf_delay_ms(2);
																printf("%s\n","MSG_REMOTE_TYPE_ADVERTISE_AVAILABLE");
                                if(STATE == STATE_CAR_WAIT_FOR_REMOTE){
																		nrf_delay_ms(2);
																		printf("%s\n","NExt -> STATE_CAR_SINGLE_MODEs");	
                                    convert_payload_to_remote_message(&remote_msg, &rx_payload);
                                    NEXT_STATE = STATE_CAR_SINGLE_MODE;
                                }
                                break;
                            case MSG_CAR_TYPE_SPEED_INFO:
                                convert_payload_to_master_message(&master_msg, &rx_payload);

                            case MSG_CAR_TYPE_REQUEST_POOLING:
                                if(STATE == STATE_CAR_SINGLE_MODE){
                                    NEXT_STATE = STATE_CAR_TRUCK_POOLING_MASTER;
                                }
                                break;

                            case MSG_CAR_TYPE_ACCEPT_POOLING:
                                if(STATE == STATE_CAR_TRUCK_POOLING_PENDING){
                                    NEXT_STATE = STATE_CAR_TRUCK_POOLING_SLAVE;
                                }
                                break;
                            case MSG_CAR_TYPE_ACKNOWLEDGE_MASTER:
                                if(STATE == STATE_CAR_TRUCK_POOLING_MASTER){
                                    receive_ack = 1;
                                }

                        }
                    }
                } // payload >0
            } // read(&rx_payload) == sucsess
            break; // break case RX_EVENT
    } // Switch
    NRF_GPIO->OUTCLR = 0xFUL << 12;
    NRF_GPIO->OUTSET = (p_event->tx_attempts & 0x0F) << 12;
} // Function


void clocks_start( void )
{
    NRF_CLOCK->EVENTS_HFCLKSTARTED = 0;
    NRF_CLOCK->TASKS_HFCLKSTART = 1;

    while (NRF_CLOCK->EVENTS_HFCLKSTARTED == 0);
}



// Changes made:
/*
* added payload.length = 3 instead of default config 32
* changed to PROTOCOL_ESB, not PROTOCOL_ESB_DPL (dynamig payload length)
* added typedef remote_packet struct remote_msg

*/
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

void radio_send_ack(){
    convert_car_message_to_payload(&car_msg, &tx_payload);
    radio_transmit_mode();
    while(nrf_esb_write_payload(&tx_payload) != NRF_SUCCESS){
        //NOP
    }
    nrf_delay_ms(1); 
    if (NRF_LOG_PROCESS() == false)
    {
     __WFE();
    }
}

void radio_wait_for_new_message(){
    radio_receive_mode();
    if (NRF_LOG_PROCESS() == false)
    {
     __WFE();
    }
}

uint32_t radio_send_and_ack_master_message(uint32_t timeout){
    convert_master_message_to_payload(&master_msg, &tx_payload);
    printf("Wait for master ack\n");   
    receive_ack = 0;
    for(uint32_t i=0; i<RETIRES; i++){
        uint32_t wait = timeout;
        radio_transmit_mode();

		nrf_delay_ms(2);
        printf("SenderID: %d \r\n", remote_msg.senderID);   
        convert_master_message_to_payload(&master_msg, &tx_payload);
        receive_ack = 0;

        while(nrf_esb_write_payload(&tx_payload) != NRF_SUCCESS){
            //NOP
        }
    
        radio_receive_mode();
        while(--wait > 0){
            if(receive_ack){
                break;
            }
            nrf_delay_ms(1);
        }
        if(receive_ack){
            printf("After timeout loop\n");        
            break;
        }
        else
            printf("timeout occred, retry nunber: %d\n", i);
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
    motor_init();
    ultrasound_init();

    STATE = STATE_CAR_TRUCK_POOLING_SLAVE;
    car_msg.type = MSG_CAR_TYPE_CONNECTED_TO;

    nrf_esb_start_rx();
    int8_t i = 0;

    int my_id = 0;

    uint32_t left_speed = 0;
    uint32_t right_speed = 0;

    uint32_t left_dir =  0;
    uint32_t right_dir = 0;

    double integral = 0;
    double last_error = 0;
    int32_t dist = 0;

    kalman_state kalman = kalman_init(0.3,3,0,0);
	
    printf("%s\n","STARTING CAR");
    motor_start();
    //unit_test_motor();
    //printf("unit test, done\n");
    
    while (true)
    {
        // Wait for message from remote or lead car.
        //radio_wait_for_new_message();
        switch(STATE){
            case STATE_CAR_WAIT_FOR_REMOTE:
                // Wait until a remote is ready to connect.
								nrf_delay_ms(2);
                printf("%s\n","STATE_CAR_WAIT_FOR_REMOTE" );
                //while(NEXT_STATE != STATE_CAR_SINGLE_MODE);
                while(NEXT_STATE != STATE_CAR_SINGLE_MODE){
                    nrf_delay_ms(1);
                }
                car_msg.type = MSG_CAR_TYPE_CONNECTED_TO;
                my_id = remote_msg.senderID;
                car_msg.senderID = my_id;
                radio_send_ack();
                set_led(my_id);
                break;
            case STATE_CAR_SINGLE_MODE:
                // Get joystick info from remote_msg and set side speeds accordingly
                printf("%s\n","STATE_CAR_SINGLE_MODE" );
                motorSpeeds(remote_msg.y, remote_msg.x, &left_speed, &right_speed);

                left_dir =  0;
                right_dir = 0;

                motorDirections(&left_speed, &right_speed, &left_dir, &right_dir);
                set_motors(left_speed, right_speed, left_dir, right_dir);

                car_msg.type = MSG_CAR_TYPE_ACKNOWLEDGE;
                radio_send_ack();
                break;

            case STATE_CAR_TRUCK_POOLING_PENDING:
                // Get joystick info from remote_msg and set side speeds accordingly
                printf("%s\n", "STATE_CAR_TRUCK_POOLING_PENDING");
                motorSpeeds(remote_msg.y, remote_msg.x, &left_speed, &right_speed);

                left_dir =  0;
                right_dir = 0;         
                motorDirections(&left_speed, &right_speed, &left_dir, &right_dir);
                set_motors(left_speed, right_speed, left_dir, right_dir);

                car_msg.type = MSG_CAR_TYPE_ACKNOWLEDGE;
                radio_send_ack();
                break;

            case STATE_CAR_TRUCK_POOLING_SLAVE:
                //printf("%s\n", "STATE_CAR_TRUCK_POOLING_SLAVE" );
                dist = ultrasound_get_distance();
                kalman_update(&kalman, (double) dist);
                dist = (int32_t) kalman.x;
                master_msg.speed_info = 0;
                double speed = get_speed(0.1, dist, master_msg.speed_info, &last_error, &integral);
                printf("Dist: %d\t Error: %f\t Speed: %f\n", dist, last_error, speed );
                uint32_t tspeed = (uint32_t) max(0,min(speed, 511));
                set_motors(tspeed, tspeed, 1, 1);
                //car_msg.type = MSG_CAR_TYPE_ACKNOWLEDGE_MASTER;
                //radio_send_ack();
                break;

            case STATE_CAR_TRUCK_POOLING_MASTER:
                printf("%s\n", "STATE_CAR_TRUCK_POOLING_MASTER" );
                motorSpeeds(remote_msg.y, remote_msg.x, &left_speed, &right_speed);

                uint32_t left_dir =  0;
                uint32_t right_dir = 0;

                motorDirections(&left_speed, &right_speed, &left_dir, &right_dir);
                set_motors(left_speed, right_speed, left_dir, right_dir);

                car_msg.type = MSG_CAR_TYPE_ACKNOWLEDGE;
                radio_send_ack();

                master_msg.speed_info = remote_msg.x;
                radio_send_and_ack_master_message(TIMEOUT);
                break;  
        }
		//STATE = NEXT_STATE;
        nrf_delay_ms(100);
    }
}


/*lint -restore */
/** @} */
