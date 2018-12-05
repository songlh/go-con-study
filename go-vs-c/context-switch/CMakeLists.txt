cmake_minimum_required(VERSION 3.10)
project(contextSwitch C)

set(CMAKE_C_STANDARD 99)


add_executable(contextSwitch taskswc.c)
target_link_libraries(contextSwitch pthread)
