import sys
import os
from os import listdir
from os.path import isfile, join
import shutil


if __name__=='__main__':
	sDirectory = sys.argv[1]
	sTarget = sys.argv[2]

	onlyfiles = [join(sDirectory, f) for f in listdir(sDirectory) if isfile(join(sDirectory, f))]

	for f in onlyfiles:
		if os.path.getsize(f)/1024/1024 > 0:
			shutil.move(f, sTarget)	