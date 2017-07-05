#!/usr/bin/env python3
# -*- coding: utf-8 *-*


from pymongo import MongoClient
import pymongo.errors

settings = {
    'host': 'mongo1.c1.lxch.lan',
    'port': 27017,
    'user': 'pcdummy',
    'password': 'gvbsEZ0UhvdAT07sI0sI',
    'authdatabase': 'admin',
    'db': 'saltstack',
    'collection': 'salt_pillar'
}

# settings = {
#     'host': 'mongo1.pcdummy.lan',
#     'port': 27017,
#     'user': 'pcdummy',
#     'password': 'ahK9eereirei3fo',
#     'authdatabase': 'admin',
#     'db': 'lxdmyadmin',
#     'collection': 'salt_pillar'
# }


def get_ips(s):
    mgo_client = MongoClient(s['host'], s['port'])
    mgo_auth_db = mgo_client[s['authdatabase']]

    try:
        mgo_auth_db.authenticate(
            s['user'],
            s['password'],
            mechanism='SCRAM-SHA-1'
        )

    except pymongo.errors.OperationFailure:
        print("Error while authenticating with user {0!r}".format(s['user']))
        return {}

    hosts = {}
    mgo_coll = mgo_client[s['db']][s['collection']]
    for doc in mgo_coll.find():
        if '_data' in doc and 'network' in doc['_data']:
            doc_net = doc['_data']['network']

            hosts[doc['_id']] = {}
            if 'interfaces' not in doc_net:
                continue

            for interface, doc_net_int in doc_net['interfaces'].items():
                if ('ipv4address' in doc_net_int or
                        'ipv6address' in doc_net_int or
                        'pubipv6address' in doc_net_int):
                    ips = {}

                    if 'ipv4address' in doc_net_int:
                        ips['ipv4'] = doc_net_int['ipv4address']

                    if 'ipv6address' in doc_net_int:
                        ips['ipv6'] = doc_net_int['ipv6address']

                    if 'pubipv6address' in doc_net_int:
                        ips['pubipv6'] = doc_net_int['pubipv6address']

                    hosts[doc['_id']][interface] = ips

    return hosts


if __name__ == '__main__':
    hosts = get_ips(settings)

    for hostname, interfaces in hosts.items():
        if not interfaces:
            print("{0}: No interfaces.".format(hostname))
            continue

        print("{0}:".format(hostname))
        for interface, int_data in interfaces.items():
            print("\t{0}\t{1}".format(
                interface,
                '\t'.join(int_data.values())
            ))
