#!/usr/bin/env python3
# -*- coding: utf-8 *-*

import pprint
from pymongo import MongoClient
import pymongo.errors


# settings = {
#     'host': 'mongo1.c1.lxch.lan',
#     'port': 27017,
#     'user': 'pcdummy',
#     'password': 'gvbsEZ0UhvdAT07sI0sI',
#     'authdatabase': 'admin',
#     'db': 'saltstack',
#     'collection': 'salt_tops'
# }

settings = {
    'host': 'mongo1.pcdummy.lan',
    'port': 27017,
    'user': 'pcdummy',
    'password': 'ahK9eereirei3fo',
    'authdatabase': 'admin',
    'db': 'lxdmyadmin',
    'collection': 'salt_tops'
}


def update_roles(s):
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

    mgo_coll = mgo_client[s['db']][s['collection']]
    for doc in mgo_coll.find():
        doc['environment'] = 'base'
        mgo_coll.update_one({'_id': doc['_id']}, {'$set': doc}, upsert=True)


        # includes = []
        # for include in doc['include']:
        #     includes.append(include['file'])

        # doc['include'] = includes

        # mgo_coll.update_one({'_id': doc['_id']}, {'$set': doc}, upsert=True)
        #pprint.pprint(includes)


if __name__ == '__main__':
    update_roles(settings)
