#!/usr/bin/env python
## -*- coding: utf-8 -*-


def read_log_test():
    count = 0
    with open('log.txt', "r", encoding='latin-1') as logf:
        lines = logf.readlines()
        for line in lines:
            if "=== RUN" in line:
                count += 1
    print("Total %d" % count)


def main():
    read_log_test()


if __name__ == '__main__':
    main()
