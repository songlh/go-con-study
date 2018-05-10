import os 
from os import listdir
from os.path import isfile, join
import re
import shutil
import sys


if __name__ == '__main__':
	sDirectory = sys.argv[1]

	files = [f for f in listdir(sDirectory) if isfile(join(sDirectory, f)) ]

	files = sorted(files)

	for f in files:
		print f[0:40]