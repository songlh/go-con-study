#include <stdio.h>
#include <pthread.h>
#include <time.h>
#include <sched.h>

#define NUM_THREADS 2
#define WORK_COUNT 1000000

struct timespec start_time, end_time;
unsigned long loop_counter = 0;
unsigned long min_latency =0; 
unsigned long max_latency=0;
unsigned long total_nanosecs=0;
double avg_latency = 0;

pthread_t threads[NUM_THREADS];
pthread_mutex_t mtx = PTHREAD_MUTEX_INITIALIZER;
void* task1() {
	while(loop_counter < 1000000){
		clock_gettime(CLOCK_REALTIME, &end_time);
		pthread_mutex_lock(&mtx);
		ReportTaskSwitchTime();
		loop_counter++;
		clock_gettime(CLOCK_REALTIME, &start_time);
		pthread_mutex_unlock(&mtx);
		sched_yield();
	}
}

unsigned int TimespecToNanosec(struct timespec tspec) {
	unsigned int nanosec = tspec.tv_sec * 1000000000 + tspec.tv_nsec;
	return(nanosec);
}

void ReportTaskSwitchTime() {
	if(loop_counter != 0 ){
		unsigned int start_time_nanosec = TimespecToNanosec(start_time);
		unsigned int end_time_nanosec = TimespecToNanosec(end_time);
		unsigned int task_switch_time = (end_time_nanosec - start_time_nanosec);
		if (min_latency == 0 && max_latency ==0) {
			min_latency = task_switch_time;
			max_latency = task_switch_time;
		}
		if (task_switch_time < min_latency) min_latency = task_switch_time; 
		if (task_switch_time > max_latency) max_latency = task_switch_time;

		total_nanosecs += task_switch_time;	
		avg_latency = total_nanosecs / loop_counter;
		printf(" Total: %d", total_nanosecs);
		printf(" Counter: %d", loop_counter);
		printf(" Current: %d", task_switch_time);
		printf(" Min: %d", min_latency);
		printf(" Avg: %d", (int)avg_latency);
		printf(" Max: %d \n", max_latency);
	}
}

int main() {
	int i;
	for(i=0; i<NUM_THREADS; i++) {
		pthread_create(&threads[i], NULL, task1, NULL);
		pthread_join(threads[i], NULL);
	}
	return(0);
}
