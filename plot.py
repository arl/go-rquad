#!/usr/bin/env python2
import sys
import re
from collections import namedtuple, defaultdict
from itertools import groupby, cycle

import numpy as np
import matplotlib.pyplot as plt

Point = namedtuple('Point', 'type name nodes time')
dots = cycle(['r', 'b', 'g'])
prog = re.compile(r'Benchmark(Basic|CNTree)(.*)Res(\d*)-\d*\s*\d*\s*(\d*)')
cmpbench = defaultdict(list)
sizes = {'Creation': 4096, 'PointLocation': 1024, 'Neighbours': 4096}

with open(sys.argv[1]) as f:
    for line in f.readlines():
        res = prog.match(line.strip())
        if res:
            grp = res.groups()
            type_ = grp[0]
            name = grp[1]
            res = int(grp[2])
            time = float(grp[3])
            dim = sizes[name]
            nodes = (dim*dim)/res
            pt = Point(type_, name, nodes, time)
            cmpbench[pt.name].append(pt)

for name, data in cmpbench.iteritems():
    # new plot starts here
    plt.title(name)
    handles = []
    for k, g in groupby(data, lambda x:x.type):
        pts_x, pts_y = [], []
        for pt in sorted(list(g), key=lambda x:x.nodes):
            pts_x.append(pt.nodes)
            pts_y.append(pt.time)
        plt.xscale('linear')
        plt.yscale('linear')
        plt.xlabel('Number of nodes', fontsize=22)
        plt.ylabel('Time per operation (nanoseconds)', fontsize=22)

        hnd, = plt.plot(pts_x, pts_y, next(dots), label=k)
        handles.append(hnd)
        plt.legend(handles=handles)
        # print 'x', pts_x
        # print 'y', pts_y

    plt.gcf().set_size_inches(18, 9)
    # plt.show()
    name = name + '.png'
    print 'output plot: ', name
    plt.savefig(name)
    plt.clf()
