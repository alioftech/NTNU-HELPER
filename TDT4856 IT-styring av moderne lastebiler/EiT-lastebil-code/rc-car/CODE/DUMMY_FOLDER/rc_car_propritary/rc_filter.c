#include "rc_filter.h"

kalman_state kalman_init(double q, double r, double p, double initial_value){
    kalman_state result;
    result.q = q;
    result.p = p;
    result.r = r;
    result.x = initial_value;

    return result;
}

void kalman_update(kalman_state* state, double measurement){
    // Update prediction
    state->p = state->p + state->q;

    // Update measurement
    state->k = state->p / (state->p + state->r);
    state->x = state->x + state->k*(measurement - state->x);
    state->p = (1 - state->k)*state->p;
}