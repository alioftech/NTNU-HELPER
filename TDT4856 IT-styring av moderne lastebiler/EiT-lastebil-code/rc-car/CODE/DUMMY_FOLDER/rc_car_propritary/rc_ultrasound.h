#ifndef RC_ULTRASOUND__
#define RC_ULTRASOUND__

#include <stdint.h>

// NB Systick counts from RVR value down to zero before reload
#define SYSTICK_BASE 0xE000E000
#define SYSTICK_CSR     ((volatile uint32_t*)(SYSTICK_BASE + 0x10)) // Control and Status Register
#define SYSTICK_RVR     ((volatile uint32_t*)(SYSTICK_BASE + 0x14)) // Reload Value Register
#define SYSTICK_CVR     ((volatile uint32_t*)(SYSTICK_BASE + 0x18)) // Current Value Register
#define SYSTICK_TOP 1000000 // Define reload value 1 million ticks

// PIN SETUP
#define ultrasound_echo 3 
#define ultrasound_trig 4

void 	 ultrasound_init(void);
uint32_t ultrasound_get_distance(void);

#endif // RC_ULTRASOUND__




/** @} */
