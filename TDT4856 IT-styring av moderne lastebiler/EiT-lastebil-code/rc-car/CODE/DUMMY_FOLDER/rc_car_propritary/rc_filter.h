#ifndef RC_FILTER__
#define RC_FILTER__

typedef struct{
    double q; // Process noise covariance
    double r; // Measurement noise covariance
    double x; // Value
    double p; // Estimation error covariance
    double k; // Kalman gain
} kalman_state;

kalman_state kalman_init(double q, double r, double p, double initial_value);
void kalman_update(kalman_state* state, double measurement);

#endif // RC_FILTER__
/** @} */
