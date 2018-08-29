#include <stdio.h>
#include <pthread.h>

int x = 0;

void* add()
{
    for( int i = 0; i < 1000000; i++ )
    {
        x++;
    }
    return NULL;
}

void* sub()
{
    for( int i = 0; i < 1000000; i++ )
    {
        x--;
    }
    return NULL;
}

int main()
{
    pthread_t add_thread;
    pthread_t sub_thread;
    pthread_create( &add_thread, NULL, add, NULL );
    pthread_create( &sub_thread, NULL, sub, NULL );

    for(int i = 0; i < 100; i++ )
    {
        printf("Current value: %i\n", x);
    }

    pthread_join( add_thread, NULL );
    pthread_join( sub_thread, NULL );

    printf("Done, result is: %i\n", x);

    return 0;
}