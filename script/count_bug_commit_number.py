#!/usr/bin/env python
# -*- coding: utf-8 -*-
import optparse
import re
import sys
import re

reCOMMIT = re.compile(r'^commit ([0-9a-f]{40})')

key_words = ['deadlock', 'deadlocks',
             'lock', 'locks', 'block', 'blocked', 'blocks', 'blocking'
             'race', 'datarace', 'dataraces',
             'synchronization', 'synchronizations', 'concurrency',
             'mutex', 'mutexes', 'atomic', 'compete', 'competes',
             'once', 'leak', 'context', 'leaked',  'leaking']


def is_useful_commit(line):
    if line.startswith('    '):
        items = line.split(' ')
        for item in items:
            if item in key_words:
                return True

    return False


def main():
    usage = 'python parseGitLog.py --inputFile log.txt --keyWord xxx'
    fmt = optparse.IndentedHelpFormatter(max_help_position=50, width=100)
    parser = optparse.OptionParser(usage=usage, formatter=fmt)

    parser.add_option('--inputFile', dest='inputFile', help='Input File')
    # parser.add_option('--keyWord', dest='keyWord', help='Key Words')

    reComment = re.compile(r'^    [^\s]')

    (opts, args) = parser.parse_args()

    numCommits = 0
    hashCommitSet = set()
    strCurrentCommit = ""

    with open(opts.inputFile, 'r', encoding='latin-1') as fin:

        while True:
            line = fin.readline()
            if not line:
                break

            match = reCOMMIT.match(line)

            if match:
                numCommits += 1
                strCurrentCommit = match.group(1)
                continue

            match = reComment.match(line)

            if match:
                if is_useful_commit(line):
                    if strCurrentCommit:
                        hashCommitSet.add(strCurrentCommit)

    print("Total number bug is %d" % len(hashCommitSet))


if __name__ == '__main__':
    main()
