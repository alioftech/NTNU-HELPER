#ifndef RC_CONTROLLER_
#define RC_CONTROLLER_

#include <stdint.h>

double get_speed(double iteration_time, int32_t measured, int32_t feed, double *last_error, double *integral);

#endif
/** @} */
