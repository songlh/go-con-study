#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function

import argparse

import requests

from bs4 import BeautifulSoup

BASE_URL = ''
BASE_GOROUTINE = 'Goroutine analysis'
MAIN_GOROUTINE = 'runtime.main'


def analysis_total_result(results):
    if len(results) == 0:
        assert 'result should not empty!' == 1

    new_results = sorted(results.iteritems(), key=lambda (k, v): (v, k), reverse=True)
    # MAIN_GOROUTINE time is max
    max_exec_time = new_results[0][1]
    percentages = {0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0}

    for _result in new_results:
        _exec_time = _result[1]
        _percentage = (_exec_time * 1.0) / max_exec_time
        if _percentage < 0.1:
            percentages[0] += 1
        elif _percentage < 0.2:
            percentages[1] += 1
        elif _percentage < 0.3:
            percentages[2] += 1
        elif _percentage < 0.4:
            percentages[3] += 1
        elif _percentage < 0.5:
            percentages[4] += 1
        elif _percentage < 0.6:
            percentages[5] += 1
        elif _percentage < 0.7:
            percentages[6] += 1
        elif _percentage < 0.8:
            percentages[7] += 1
        elif _percentage < 0.9:
            percentages[8] += 1
        elif 0.9 <= _percentage:
            percentages[9] += 1

    print("go goroutine execution time percentage report:")
    print("This go trace log have totally {0} goroutines,"
          " max execution time is {1} ns.".format(len(results), max_exec_time))

    for key in percentages:
        print("{0}% --------- {1}%, {2}".format(key*10, (key+1)*10, percentages[key]))


def parse_each_page(url):
    results = dict()
    r = requests.get(url)

    if r.status_code != 200:
        assert 'Could not open a url at %s' % url == 1

    c = r.content
    soup = BeautifulSoup(c, "html.parser")

    trs = soup.find_all('tr')
    for tr in trs[1:]:
        tds = tr.find_all('td')
        goid = tds[0].get_text()
        results[goid] = int(tds[1].get_text())

    return results


def main():
    parser = argparse.ArgumentParser()

    parser.add_argument("-u", "--url", help="url",
                        type=str)

    args = parser.parse_args()

    if not args.url:
        assert "Must input a url" == 1

    BASE_URL = args.url
    r = requests.get(BASE_URL)

    if r.status_code != 200:
        assert 'Could not open a url at %s' % BASE_URL == 1

    c = r.content
    soup = BeautifulSoup(c, "html.parser")
    links = soup.find_all('a')
    goroutine_analysis_url = None

    for link in links:
        if link.get_text() == BASE_GOROUTINE:
            goroutine_analysis_url = BASE_URL + link.get('href')

    goroutine_page = requests.get(goroutine_analysis_url)
    if goroutine_page.status_code != 200:
        assert 'Could not open a url at %s' % goroutine_analysis_url == 1

    goroutine_page_content = goroutine_page.content
    soup = BeautifulSoup(goroutine_page_content, "html.parser")
    goroutine_links = soup.find_all('a')

    total_results = dict()

    for _link in goroutine_links:
        goroutine_link = BASE_URL + _link.get('href')
        results = parse_each_page(goroutine_link)
        total_results.update(results)

    analysis_total_result(total_results)


if __name__ == '__main__':
    main()
