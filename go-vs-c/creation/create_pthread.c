#include <stdio.h>
#include <pthread.h>
#include <time.h>
#include <sched.h>

int NUM_THREADS = 10000;

void *consumer() {
    sched_yield();
}

unsigned int TimespecToNanosec(struct timespec tspec) {
    unsigned int nanosec = tspec.tv_sec * 1000000000 + tspec.tv_nsec;
    return nanosec;
}

void PrintElapsed(struct timespec start_time, struct timespec end_time, int num_threads) {
    unsigned int start_time_nanosec = TimespecToNanosec(start_time);
    unsigned int end_time_nanosec = TimespecToNanosec(end_time);
    unsigned int elapsed_nanosec = end_time_nanosec - start_time_nanosec;
    printf("The number of C thread is %d, creation time is %u ns\n", num_threads, elapsed_nanosec);
}

int main() {
    int i;
//    int err, j = 0;
    struct timespec start_time, end_time;
    pthread_t threads[NUM_THREADS];

    // timing start
    clock_gettime(CLOCK_REALTIME, &start_time);

    for (i = 0; i < NUM_THREADS; i++) {
        pthread_create(&threads[i], NULL, consumer, NULL);
//        err = pthread_create(&threads[i], NULL, consumer, NULL);
//        if (err == 0) {
//            j += 1;
//        }
    }

    // timing end
    clock_gettime(CLOCK_REALTIME, &end_time);
    PrintElapsed(start_time, end_time, NUM_THREADS);
    return (0);
}
