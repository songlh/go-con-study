import os 
from os import listdir
from os.path import isfile, join
import re
import shutil
import sys
from sets import Set


if __name__ == '__main__':
	sLogFile = sys.argv[1]

	reGoRAdd = re.compile(r'\+[\s]+go [0-9_a-zA-Z.\[\]]+\(')
	reGoRRemove = re.compile(r'-[\s]+go [0-9_a-zA-Z.\[\]]+\(')
	reCOMMIT = re.compile(r'^commit ([0-9a-f]{40})')

	strCurrentCommit = ''

	setGoAdd = Set([])
	setGoRemove = Set([])

	f = open(sLogFile, 'r')

	while True:
		line = f.readline()
		if not line:
			break

		match = reCOMMIT.match(line)

		if match:
			strCurrentCommit = match.group(1)
			continue

		match = reGoRAdd.match(line)
		if match:
			setGoAdd.add(strCurrentCommit)
			continue

		match = reGoRRemove.match(line)
		if match:
			setGoRemove.add(strCurrentCommit)
			continue

	f.close()

	goAddFile = open('go_create.txt', 'w')
	for key in setGoAdd:
		goAddFile.write(key)
		goAddFile.write('\n')

	goAddFile.close()

	goRemoveFile = open('go_remove.txt', 'w')
	for key in setGoAdd:
		goRemoveFile.write(key)
		goRemoveFile.write('\n')

	goRemoveFile.close()