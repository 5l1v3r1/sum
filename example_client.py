#!/usr/bin/python
import sys
import os
import random
import datetime
import grpc

path = os.path.join(os.path.abspath(os.path.dirname(__file__)), "proto")
sys.path.append(path)

import sum_pb2
import sum_pb2_grpc 

start = 0
end = 0
num_columns = 475
num_rows = 3000
index = {}
client = None
oracle_name = 'dotAll'
oracle_id = None

def gen_record(columns):
    return sum_pb2.Record( \
        data=[random.uniform(0,100) for i in range(0, columns)],
        meta=[sum_pb2.NamedValue( \
            name="example_metadata",
            value="Random number is %f" % random.random()
        )]
    )

def check(resp):
    if resp.success == False:
        print "ERROR: %s" % resp.msg
        quit()

def timer_start():
    global start
    sys.stdout.flush()
    start = datetime.datetime.now()

def timer_stop(with_avg=True):
    global start, end, index, client
    end = datetime.datetime.now()
    diff = end - start
    elapsed_ms = (diff.days * 86400000) + (diff.seconds * 1000) + (diff.microseconds / 1000)
    if with_avg:
        print "%d ms / %.2fms avg" % ( elapsed_ms, float(elapsed_ms) / float(len(index)) )
    else:
        print "%d ms" % elapsed_ms

def define_oracle(filename, name):
    global client

    resp = client.FindOracle(sum_pb2.ByName(name=name))
    check(resp)

    if len(resp.oracles) == 0:
        print "Defining oracle %s ..." % name
        with open( filename, 'r') as fp:
            oracle = sum_pb2.Oracle(name=name, code=fp.read())
            resp = client.CreateOracle(oracle)
            check(resp)
            print "  -> id:%s" % resp.msg
            return resp.msg

    else:
        o = resp.oracles[0]
        print "Oracle %s -> id:%s" % ( o.name, o.id )
        return o.id

    return None

if __name__ == '__main__':
    client = sum_pb2_grpc.SumServiceStub(grpc.insecure_channel('127.0.0.1:50051'))

    oracle_id = define_oracle('example_oracle.js', oracle_name)
    print

    print "CREATE (%dx%d) : " % ( num_rows, num_columns ),
    timer_start()
    for row in range (0, num_rows):
        record = gen_record(num_columns)
        resp = client.CreateRecord(record)
        check(resp)
        # msg contains the identifier
        index[resp.msg] = record
    timer_stop()

    print "CALL %s x%d : " % (oracle_name, len(index)),
    timer_start()
    for ident, record in index.iteritems():
        resp = client.Run(sum_pb2.Call(oracle_id=oracle_id, args=("\"%s\"" % ident, "0.1",)))
        check(resp)
        break
        # print resp.json
    timer_stop(False)

    print "DEL x%d : " % len(index),
    timer_start()
    for ident, record in index.iteritems():
        check( client.DeleteRecord(sum_pb2.ById(id=ident)) )
    timer_stop()