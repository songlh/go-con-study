import optparse
import re
import sys

import os
import subprocess


def main():
	usage = 'python dumpPatches.py --inputFile race.result'
	fmt = optparse.IndentedHelpFormatter(max_help_position=50, width=100)
	parser = optparse.OptionParser(usage=usage, formatter=fmt)

	parser.add_option('--inputFile', dest='inputFile', help='Input File')
	
	(opts, args) = parser.parse_args()

	with open(opts.inputFile, 'r') as fin:

		sDirectoryName = opts.inputFile.split('.')[0]

		if not os.path.exists(sDirectoryName):
			os.makedirs(sDirectoryName)

		while True:
			line = fin.readline()
			if not line:
				break

			sFileName = os.path.join(sDirectoryName, line[:-1]) 

			if not os.path.exists(sFileName):
				process = subprocess.Popen(['git', 'show', line[:-1]], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
				
				with open(sFileName, 'w') as fout: 
					fout.write(process.communicate()[0])


if __name__ == '__main__':
	main()