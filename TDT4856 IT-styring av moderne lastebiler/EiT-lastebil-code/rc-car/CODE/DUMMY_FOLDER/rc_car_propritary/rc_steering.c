#include "rc_steering.h"
#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <string.h>


#define maxVal 512


// Scale output to 9 bit and set direction.
void steering_speeds(uint32_t inturn, uint32_t inthrottle, uint32_t *left_speed, uint32_t *right_speed, uint32_t *left_dir, uint32_t *right_dir){
	// Scale from 0..1023 to -512...511
	int32_t turn = inturn - 512;
	int32_t throttle = inthrottle - 512;
	
	if(abs(throttle) < 20){
		throttle = 0;
	}

	if(abs(turn) < 10){
		turn = 0;
	}

	if(throttle > 0){
		throttle = throttle*0.7 + 150;
    }else if(throttle < 0){
    	throttle = throttle*0.7 - 150;
    }


	int32_t left_side = throttle;
	int32_t right_side = throttle;
	

	if(throttle >= 0){
		left_side = throttle + turn;
		right_side = throttle - turn;	
	}else{
		left_side = throttle - turn;
		right_side = throttle + turn;
	}

	if(left_side >= 0){
        *left_dir = 1;
    }else {
    	*left_dir = 0;
    	left_side = abs(left_side);
    }

    if(right_side >= 0){
        *right_dir = 1;
    }else{
    	*right_dir = 0;
        right_side = abs(right_side);
    }


    double norm = 1;
	if(left_side > right_side && left_side > maxVal){
		norm = (double) left_side/maxVal;
	}else if(right_side > left_side && right_side > maxVal){
		norm = (double) right_side/maxVal;
	}

	double t_left = left_side/norm;
	double t_right = right_side/norm;
   
    *left_speed = (uint32_t) round(t_left);
    *right_speed = (uint32_t) round(t_right);
}


