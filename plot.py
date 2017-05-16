#! /usr/bin/env python2

import sys
import re
from collections import defaultdict, namedtuple
from itertools import cycle, groupby

import matplotlib.pyplot as plt


Point = namedtuple('Point', 'type nodes time')
dots = cycle(['r', 'g'])
LINE_PROG = re.compile(r'Benchmark(Basic|CNTree)(.*)Res(\d*)-\d*\s*\d*\s*(\d*)')
sizes = {'Creation': 4096, 'PointLocation': 1024, 'Neighbours': 4096}


def rquad_line_parser(line):
    res = LINE_PROG.match(line)
    if res:
        grp = res.groups()
        type_ = grp[0]
        name = grp[1]
        res = int(grp[2])
        time = float(grp[3])
        dim = sizes[name]
        nodes = (dim * dim) / res
        return name, Point(type_, nodes, time)


def extract(filename, parser):
    """
    Extract data from benchmarks report `filename`

    `parser` is a user-defined function taking a benchmark line, and returning a
    2-tuple (plotname, DATAPOINT) where plotname is the name of the plot on
    which DATAPOINT should be plotted, and DATAPOINT is the domain-specific
    object representing the values of this line.
    """
    benchs = defaultdict(list)
    with open(filename) as f:
        for line in f.readlines():
            parsed = parser(line)
            if parsed:
                name, datapoint = parsed
                benchs[name].append(datapoint)
    return benchs


def rquad_plot(name, dps):
    print 'plotting: ', name
    plt.title(name, fontsize=22)
    handles = []
    for k, g in groupby(dps, lambda x: x.type):
        pts_x, pts_y = [], []
        for pt in sorted(list(g), key=lambda x: x.nodes):
            pts_x.append(pt.nodes)
            pts_y.append(pt.time)
        plt.xscale('linear')
        plt.yscale('linear')
        plt.xlabel('Number of nodes', fontsize=22)
        plt.ylabel('Time per operation (nanoseconds)', fontsize=22)

        hnd, = plt.plot(pts_x, pts_y, next(dots), label=k)
        handles.append(hnd)
        plt.legend(handles=handles, prop={'size': 22})
        # print 'x', pts_x
        # print 'y', pts_y

    plt.gcf().set_size_inches(18, 9)
    # plt.show()
    filename = name + '.png'
    print 'output plot: ', filename
    plt.savefig(filename)
    plt.clf()


def main():
    if len(sys.argv) < 2:
        print "usage: plot.py FILENAME"
        return
    benchs = extract(sys.argv[1], rquad_line_parser)
    for plot_name in benchs.keys():
        rquad_plot(plot_name, benchs[plot_name])


if __name__ == "__main__":
    main()
