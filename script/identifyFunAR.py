import os 
from os import listdir
from os.path import isfile, join
import re
import shutil
import sys
from sets import Set

reGoRAdd = re.compile(r'\+[\s]+go [0-9_a-zA-Z.\[\]]+\(')
reGoRRemove = re.compile(r'-[\s]+go [0-9_a-zA-Z.\[\]]+\(')
reAddFuncStart = re.compile(r'\+([\s]*)func[\s]+')
reRemoveFuncStart = re.compile(r'\-([\s]+)func[\s]+[0-9_a-zA-Z]+\(')
reCOMMIT = re.compile(r'^commit ([0-9a-f]{40})')
reDiff = re.compile(r'diff --git')

def isInsideFunctionCreation(lines, numStart):
	ii = numStart
	while ii >= 0:
		match = reAddFuncStart.match(lines[ii])
		if match:
			sBlank = match.group(1)

			jj = numStart + 2
			while jj < len(lines):
				if lines[jj].startswith('+'+sBlank+'}'):
					return True

				if not lines[jj].startswith('+'):
					return False

				match = reDiff.match(lines[ii])
				if match:
					return False

				jj +=1

			continue

		match = reDiff.match(lines[ii])
		if match:
			return False

		if not lines[ii].startswith('+'):
			return False

		ii -=1


	return False

def isInsideFunctionRemove(lines, numStart):
	ii = numStart
	while ii >= 0:
		match = reRemoveFuncStart.match(lines[ii])
		if match:
			sBlank = match.group(1)

			jj = numStart + 2
			while jj < len(lines):
				if lines[jj].startswith('-'+sBlank+'}'):
					return True

				if not lines[jj].startswith('-'):
					return False

				match = reDiff.match(lines[ii])
				if match:
					return False

				jj +=1

			continue

		match = reDiff.match(lines[ii])
		if match:
			return False

		if not lines[ii].startswith('-'):
			return False

		ii -=1


	return False



if __name__ == '__main__':
	sDirectory = sys.argv[1]
	files = [join(sDirectory, f) for f in listdir(sDirectory) if isfile(join(sDirectory, f)) ]

	setFunctionCreation = Set([])
	setAdd = Set([])

	#files = ['./etcd/go/create/new-function/cce88a85049b73898490ae0ee516537ca5e371b5']

	for f in files:		
		flag = False
		with open(f, 'r') as ff:
			lines = ff.readlines()

			for index in range(len(lines)):
				match = reGoRAdd.match(lines[index])

				if match:

					ii = index - 1

					if isInsideFunctionCreation(lines, ii):
						setFunctionCreation.add(f)
						flag = True
					else:
						#print lines[index]
						setAdd.add(f)

					
					continue

		if not flag:
			print f

	#for f in setAdd:
	#	print f