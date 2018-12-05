#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>


#define BUF_SIZE 64
#define WORK_COUNT 1000000
#define NP 32

typedef struct {
    int buf[BUF_SIZE]; // the buffer
    size_t len; // number of items in the buffer
    pthread_mutex_t mutex; // needed to add/remove data from the buffer
    pthread_cond_t can_produce; // signaled when items are removed
    pthread_cond_t can_consume; // signaled when items are added
} buffer_t;


unsigned long long TimespecToNanosec(struct timespec tspec) {
    unsigned long long nanosec = tspec.tv_sec * 1000000000 + tspec.tv_nsec;
    return nanosec;
}

void PrintElapsed(struct timespec start_time, struct timespec end_time, int num_threads) {
    unsigned long long start_time_nanosec = TimespecToNanosec(start_time);
    unsigned long long end_time_nanosec = TimespecToNanosec(end_time);
    unsigned long long elapsed_nanosec = end_time_nanosec - start_time_nanosec;
    printf("The number of C thread is %d, communication time is %llu ns\n", num_threads, elapsed_nanosec);
}

void *producer(void *arg) {
    buffer_t *buffer = (buffer_t *) arg;
    int i;

    for (i = 0; i < WORK_COUNT; i++) {
//    while(1) {
        pthread_mutex_lock(&buffer->mutex);

        while (buffer->len == BUF_SIZE) { // full
            pthread_cond_wait(&buffer->can_produce, &buffer->mutex);
        }
        int t = i;
//        printf("Produced: %d\n", t);

        buffer->buf[buffer->len] = t;
        ++buffer->len;
        pthread_cond_signal(&buffer->can_consume);
        pthread_mutex_unlock(&buffer->mutex);
    }

    return NULL;
}

void *consumer(void *arg) {
    buffer_t *buffer = (buffer_t *) arg;
    int i;

    for (i = 0; i < WORK_COUNT; i++) {
//    while(1) {
        pthread_mutex_lock(&buffer->mutex);

        while (buffer->len == 0) { // empty
            pthread_cond_wait(&buffer->can_consume, &buffer->mutex);
        }
        --buffer->len;
//        printf("Consumed: %d\n", buffer->buf[buffer->len]);
        pthread_cond_signal(&buffer->can_produce);

        pthread_mutex_unlock(&buffer->mutex);
    }

    return NULL;
}

int main(int argc, char *argv[]) {

    struct timespec start_time, end_time;

    buffer_t buffer = {
            .len = 0,
            .mutex = PTHREAD_MUTEX_INITIALIZER,
            .can_produce = PTHREAD_COND_INITIALIZER,
            .can_consume = PTHREAD_COND_INITIALIZER
    };

    pthread_t threads[NP * 2];
    int thread_id, i;
    // timing start
    clock_gettime(CLOCK_REALTIME, &start_time);

    for (thread_id = 0; thread_id < NP; thread_id++) {
        pthread_create(&threads[thread_id], NULL, producer, (void *) &buffer);
    }

    for (; thread_id < NP * 2; thread_id++) {
        pthread_create(&threads[thread_id], NULL, consumer, (void *) &buffer);
    }

    for (i = 0; i < NP * 2; i++) {
        pthread_join(threads[i], NULL);
    }

    // timing end
    clock_gettime(CLOCK_REALTIME, &end_time);
    PrintElapsed(start_time, end_time, NP * 2);

    return 0;
}
