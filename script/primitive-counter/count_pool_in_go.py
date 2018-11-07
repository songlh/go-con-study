#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os

GO_EDN = '/*.go'
GO_TEST_END = '_test.go'

pool_usage_count = 0


def valid_end(fname):
    if 'vendor' in fname or 'third_library' in fname or 'third_party' in fname:
        return False

    if '.go' in fname and '_test.go' not in fname:
        return True


def list_go_file(path):
    global pool_usage_count
    for root, subdirs, files in os.walk(path):
        for file in files:
            file_name = os.path.join(root, file)
            if valid_end(file_name):
                # print(file_name)
                with open(file_name, 'r') as f:
                    lines = f.readlines()
                    for line in lines:
                        items = line.split()
                        for item in items:
                            if 'pool.Get()' in item or 'pool.Put(' in item:
                                print(item)
                                pool_usage_count += 1


def main():
    global pool_usage_count
    path = "/home/kevin/GoStudy/src/github.com/coreos/etcd"
    list_go_file(path)
    print("pool usage:", pool_usage_count)


if __name__ == '__main__':
    main()
