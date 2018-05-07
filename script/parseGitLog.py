import optparse
import re
import sys

from sets import Set

reCOMMIT = re.compile(r'^commit ([0-9a-f]{40})')


def main():

	usage = 'python parseGitLog.py --inputFIle log.txt --keyWord xxx'
	fmt = optparse.IndentedHelpFormatter(max_help_position=50, width=100)
	parser = optparse.OptionParser(usage=usage, formatter=fmt)

	parser.add_option('--inputFile', dest='inputFile', help='Input File')
	#parser.add_option('--keyWord', dest='keyWord', help='Key Words')

	(opts, args) = parser.parse_args()

	##opts.keyWord = opts.keyWord.lower()

	numCommits = 0
	strCurrentCommit = ''
	setChanCommits = Set([])
	setLockCommits = Set([])

	with open(opts.inputFile, 'r') as fin:
		while True:
			line = fin.readline()
			if not line:
				break

			match = reCOMMIT.match(line)

			if match:
				numCommits += 1
				strCurrentCommit = match.group(1)
				continue

			if line.find('make(chan ') != -1 and (line.startswith('+') or line.startswith('-') ) :
				#if line.startswith('+') or line.startswith('-'):
				print line[:-1]
				if strCurrentCommit not in setChanCommits:
					print strCurrentCommit
					setChanCommits.add(strCurrentCommit)





if __name__ == '__main__':
	main()