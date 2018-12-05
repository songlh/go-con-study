cmake_minimum_required(VERSION 3.10)
project(creation C)

set(CMAKE_C_STANDARD 99)


add_executable(creation create_pthread.c)
target_link_libraries(creation pthread)
