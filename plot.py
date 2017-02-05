#!/usr/bin/env python2
import sys
import re
from collections import namedtuple, defaultdict
from itertools import groupby, cycle

import numpy as np
import matplotlib.pyplot as plt

Point = namedtuple('Point', 'type name resolution time')

# dots = cycle(['r--', 'bs'])
dots = cycle(['r', 'b', 'g'])

prog = re.compile(r'Benchmark(Basic|CNTree)(.*)Res(\d*)-\d*\s*\d*\s*(\d*)')

cmpbench = defaultdict(list)

with open(sys.argv[1]) as f:
    for line in f.readlines():
        res = prog.match(line.strip())
        if res:
            grp = res.groups()
            pt = Point(grp[0], grp[1], int(grp[2]), float(grp[3]))
            cmpbench[pt.name].append(pt)


for name, data in cmpbench.iteritems():
    # new plot starts here
    plt.title(name)
    for k, g in groupby(data, lambda x:x.type):
        pts_x, pts_y = [], []
        for pt in sorted(list(g), key=lambda x:x.resolution, reverse=True):
            pts_x.append(pt.resolution)
            pts_y.append(pt.time)
        plt.xlabel('tree resolution (pixels)', fontsize=18)
        plt.xscale('linear')
        plt.yscale('linear')
        plt.ylabel('time per operation (nanoseconds)', fontsize=18)

        ax = plt.gca()
        ax.set_xlim(max(pts_x), min(pts_x))

        plt.plot(pts_x, pts_y, next(dots))
        # print 'x', pts_x
        # print 'y', pts_y

    # plt.show()
    plt.gcf().set_size_inches(18, 9)
    name = name + '.png'
    print 'output plot: ', name
    plt.savefig(name)
    plt.clf()