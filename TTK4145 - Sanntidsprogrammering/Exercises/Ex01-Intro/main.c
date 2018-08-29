#include <pthread.h>
#include <stdio.h>

int i = 0;

pthread_mutex_t lock;

void* threadFunction_1(){
	for(int k = 0; k < 1000003; k++){
		pthread_mutex_lock(&lock);
		i++;
		pthread_mutex_unlock(&lock);
	}

}

void* threadFunction_2(){
	for(int k = 0; k < 1000000; k++){
		pthread_mutex_lock(&lock);
		i--;
		pthread_mutex_unlock(&lock);
	}
}

int main(){
	pthread_t thread_1;
	pthread_t thread_2;

	pthread_create(&thread_1, NULL, threadFunction_1, NULL);
	pthread_create(&thread_2, NULL, threadFunction_2, NULL);

	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);
	printf("EX 2 After running both threads, i = %d\n", i);
	
	return 0;
}