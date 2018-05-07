import os 
from os import listdir
from os.path import isfile, join
import re
import shutil
import sys


if __name__ == '__main__':
	sDirectory = sys.argv[1]

	files = [join(sDirectory, f) for f in listdir(sDirectory) if isfile(join(sDirectory, f)) ]

	reGoRAdd = re.compile(r'\+[\s]+go [0-9_a-zA-Z.\[\]]+\(')
	reGoRRemove = re.compile(r'-[\s]+go [0-9_a-zA-Z.\[\]]+\(')

	for f in files:
		ff = open(f)
		fAdd = False
		fRemove = False
		while True:
			line = ff.readline()

			if not line:
				break

			match = reGoRAdd.match(line)

			if match:
				fAdd = True

			match = reGoRRemove.match(line)

			if match:
				fRemove = True

		if fAdd or fRemove:
			print f,
			if fAdd:
				print 'add',

			if fRemove:
				print 'remove',

			print

		ff.close()
