cmake_minimum_required(VERSION 3.10)
project(communication C)

set(CMAKE_C_STANDARD 99)


add_executable(communication pc.c)
target_link_libraries(communication pthread)
