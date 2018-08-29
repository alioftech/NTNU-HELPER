#ifndef RC_MOTOR__
#define RC_MOTOR__


// RIGHT
#define motor_pwm_a     16 // NB NB 16
#define motor_in_1_a    15
#define motor_in_2_a    14
// LEFT
#define motor_pwm_b     11 // NB NB 11
#define motor_in_3_b    13 // 13
#define motor_in_4_b    12 // 12

#define PWM_TOP_VALUE 511

#define BACKWARD 0
#define FORWARD 1
#define LEFT 0
#define RIGHT 1
#define BOTH 2

void motor_init(void);
void motor_set_dir(uint32_t side, uint32_t dir);
void motor_set_speed(uint32_t side, uint32_t speed);
void motor_stop(void);
void motor_start(void);
void set_motors(uint32_t left_speed, uint32_t right_speed, uint32_t left_dir, uint32_t right_dir);

#endif // RC_MOTOR__

/*
11 HIGH IS ENABLE
12 LOW, 13 HIGH IS BACKWARD
12 HIGH, 13 LOW IS FOWRARD
*/



/** @} */
