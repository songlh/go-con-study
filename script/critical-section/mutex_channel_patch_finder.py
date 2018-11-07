#!/usr/bin/python
# -*- coding: utf-8 -*-

import argparse
import os
import re
import shlex
import subprocess

COMMAND_GIT_DIFF = "git diff %s^!"
RE_GO_FILE = re.compile("[a-zA-Z0-9_\/]+\.(go)$")

# constants
CHANGE_MUTEX_TO_CHANNEL = 1
CHANGE_CHANNEL_TO_MUTEX = 2

def subproc_output(cmd, **kwargs):
    """
    Function to execute a command in subprocess
    :param cmd:
    :param kwargs:
    :return:
    """
    kwargs['shell'] = True
    kwargs['stdout'] = subprocess.PIPE
    kwargs['stderr'] = subprocess.STDOUT
    if 'shell' in kwargs:
        return subprocess.Popen(cmd, **kwargs).stdout
    else:
        return subprocess.Popen(shlex.split(cmd), **kwargs).stdout


def get_commit_content(commit_hash):
    """
    Get Commit Content by commit hash
    :param commit_hash: sas111111111aaa
    :return: the content of this commit
    """
    diff_content = (subproc_output(
        COMMAND_GIT_DIFF % commit_hash).read())
    return diff_content


def split_code_blocks(lines):
    """
    Split line to code blocks
    :param lines: all lines changed in one file.
    :return: code blocks, which are start with @@
    """
    cur_block = []
    blocks = []
    _record = False

    for index, line in enumerate(lines):
        line = line.replace("\t", " ")
        # a new changed file
        if line.startswith("diff --git "):
            # if this diff is not related with go file, ignore it
            items = line.split(" ")
            if not RE_GO_FILE.match(items[2]):
                return []

            _record = False
            if cur_block:
                blocks.append(cur_block)
                cur_block = []

        # a new code block
        if line.startswith("@@"):
            if cur_block:
                blocks.append(cur_block)
                cur_block = []
            _record = True

        else:
            # if current line is a comment, we should ignore it.
            if (line.startswith("//") or line.startswith("+  //")
                    or line.startswith("-  //") or line.startswith("+//") or
                    line.startswith("-//") or line.startswith("+ //") or
                    line.startswith("- //")):
                continue

            if _record:
                cur_block.append(line)

    return blocks


def _parse_block(block):
    """
    The parser for parsing a block
    :param block: a code block contains some added,
     minus and unchanged code.
    :return: add_lines and minus_lines
    """
    add_lines = []
    minus_lines = []

    for line in block:
        line = line.replace("\t", " ")
        # added line
        if line.startswith("+") and line != "+":
            add_lines.append(line)
        # minus line
        if line.startswith("-") and line != "-":
            minus_lines.append(line)

    return add_lines, minus_lines


def _is_lock_op(item):
    if '.Lock()' in item or '.Unlock()' in item:
        return True
    if '.RLock()' in item or '.RLocker()' in item or '.RUnlock()' in item:
        return True

    return False


def _is_channel_op(item):
    if item == "<-" or item.startswith("<-"):
        return True

    return False


def judge_channel_or_mutex(add_lines, minus_lines):
    minus_channel = 0
    added_channel = 0
    minus_mutex = 0
    added_mutex = 0

    for line in minus_lines:
        items = line.split(" ")
        for item in items:
            if _is_channel_op(item):
                minus_channel += 1

            if _is_lock_op(item):
                minus_mutex += 1

    for line in add_lines:
        items = line.split(" ")
        for item in items:
            if _is_lock_op(item):
                added_mutex += 1
            if _is_channel_op(item):
                added_channel += 1

    # TODO: this is maybe change mutex to channel
    if (minus_mutex > added_mutex) and (added_channel > minus_channel):
        return CHANGE_MUTEX_TO_CHANNEL

    # TODO: this is maybe change channel to mutex
    if (minus_channel > added_channel) and (added_mutex > minus_mutex):
        return CHANGE_CHANNEL_TO_MUTEX


def parse_blocks(blocks):
    """
    The parser used to parse all blocks
    :param blocks: code blocks
    :return: [(added_code, minus_code)]
    """
    changed_blocks = []
    for block in blocks:
        add_lines, minus_lines = _parse_block(block)
        # TODO: here we only consider the code block has add and minus
        if add_lines and minus_lines:
            changed_blocks.append((add_lines, minus_lines))
            if judge_channel_or_mutex(add_lines, minus_lines) \
                    == CHANGE_CHANNEL_TO_MUTEX:
                print("change channel to mutex: ", block)

            if judge_channel_or_mutex(add_lines, minus_lines) \
                    == CHANGE_MUTEX_TO_CHANNEL:
                print("change mutex to channel: ", block)

    return changed_blocks


def parse_diff(content, commit_hash=None):
    """
    Parse diff content
    :param content:
    :param commit_hash:
    :return:
    """
    change_channel_to_mutex_content = set()
    change_mutex_to_channel_content = set()
    lines = content.decode().split("\n")
    blocks = split_code_blocks(lines)
    changed_blocks = parse_blocks(blocks)


def main():
    parser = argparse.ArgumentParser(description='Git commits diff parser for mining')
    parser.add_argument('--input', type=str, default='commit-hash-log.txt',
                        help='The input file of commit hash log file used to parse.'
                             '"git log --pretty="%H" > commit-hash-log.txt"')

    args = parser.parse_args()

    if not os.path.exists(args.input):
        raise ValueError('Could not find input file!')

    with open(args.input, 'r') as commit_file:
        commits = commit_file.readlines()
        for commit_hash in commits:
            commit_hash = commit_hash.strip("\n")
            diff_content = get_commit_content(commit_hash)
            parse_diff(diff_content, commit_hash)


if __name__ == '__main__':
    main()
