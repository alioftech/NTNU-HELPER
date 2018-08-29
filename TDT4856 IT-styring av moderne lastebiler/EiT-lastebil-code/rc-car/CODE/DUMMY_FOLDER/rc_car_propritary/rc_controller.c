#include "rc_controller.h"
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <math.h>
#include <string.h>

static int32_t DESIRED_DISTANCE = 100; //mm
static double Kp = 10;
static double Kd = 7;
static double Ki = 4;
static double Kf = 1;

#define max(X, Y) (((X) > (Y)) ? (X) : (Y))
#define min(X, Y) (((X) < (Y)) ? (X) : (Y))
#define sign(x) ((x > 0) ? 1 : ((x < 0) ? -1 : 1))

double get_speed(double iteration_time, int32_t measured, int32_t feed, double *last_error, double *integral){
	double error = (double) (measured - DESIRED_DISTANCE);
	*integral += error*iteration_time;
	double derivative = (error - (*last_error))/iteration_time;
	*last_error = error;
	double output =  Kp*error + Kd*derivative + Ki*(*integral) + Kf*feed;
	if(output > 511){
		*integral -= error*iteration_time;
		output = 511;
	}else if(output < 0){
		*integral -=error*iteration_time;
		output = 0;
	}
	return output;
}