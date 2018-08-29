#include <stdio.h>
#include <pthread.h>
#include <semaphore.h>

int x = 0;
sem_t mutex;

void* add()
{
    for( int i = 0; i < 1000000; i++ )
    {
        sem_wait(&mutex);
        x++;
        sem_post(&mutex);
    }
    return NULL;
}

void* sub()
{
    for( int i = 0; i < 1000000; i++ )
    {
        sem_wait(&mutex);
        x--;
        sem_post(&mutex);
    }
    return NULL;
}

int main()
{
    pthread_t add_thread;
    pthread_t sub_thread;

    sem_init(&mutex, 0, 1);

    pthread_create( &add_thread, NULL, add, NULL );
    pthread_create( &sub_thread, NULL, sub, NULL );

    for(int i = 0; i < 100; i++ )
    {
        printf("Current value: %i\n", x);
    }

    pthread_join( add_thread, NULL );
    pthread_join( sub_thread, NULL );

    sem_destroy(&mutex);

    printf("Done, result is: %i\n", x);

    return 0;
}