#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os

GO_EDN = '/*.go'
GO_TEST_END = '_test.go'

anony_usage_count = 0
non_anon_usage_count = 0

def valid_end(fname):
    if 'vendor' in fname or 'third_library' in fname or 'third_party' in fname:
        return False

    if '.go' in fname and '_test.go' not in fname:
        return True


def list_go_file(path):
    global non_anon_usage_count, anony_usage_count
    for root, subdirs, files in os.walk(path):
        for file in files:
            file_name = os.path.join(root, file)
            if valid_end(file_name):
                # print(file_name)
                with open(file_name, 'r') as f:
                    lines = f.readlines()
                    for line_num, line in enumerate(lines):
                        items = line.split()
                        for index, item in enumerate(items):
                            if 'go' == item:
                                if len(items) == (index+1):
                                    # print(items)
                                    break
                                if 'func(' in items[index+1]:
                                    anony_usage_count += 1
                                    print(file_name, line_num)
                                    print(items[index], items[index+1])
                                elif '(' in items[index+1]:
                                    non_anon_usage_count += 1
                                    print(file_name, line_num)
                                    print(items[index], items[index+1])
                                break


def main():
    global channel_usage_count
    path = "/home/kevin/GoStudy/src/github.com/coreos/etcd"
    list_go_file(path)
    print("non anonymouse call", non_anon_usage_count)
    print("anonymous call", anony_usage_count)


if __name__ == '__main__':
    main()
