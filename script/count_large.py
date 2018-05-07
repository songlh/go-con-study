import sys
import os
from os import listdir
from os.path import isfile, join
import shutil


if __name__=='__main__':
	sDirectory = sys.argv[1]
	#sTarget = sys.argv[2]

	onlyfiles = [join(sDirectory, f) for f in listdir(sDirectory) if isfile(join(sDirectory, f))]

	for f in onlyfiles:
		#if os.path.getsize(f)/1024/1024 > 0:
			#shutil.move(f, sTarget)
		ff = open(f)
		while True:
			line = ff.readline()
			if not line:
				break

			if line.startswith('+') and line.find('go ') != -1:
				if len(line.strip()) > 30:
					print line

		ff.close()	