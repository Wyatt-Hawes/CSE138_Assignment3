################### 
# Course: CSE138
# Date: Fall 2023
# Assignment: 3
# Authors: Reza NasiriGerdeh, Zach Gottesman, Lindsey Kuper, Patrick Redmond
# This document is the copyrighted intellectual property of the authors.
# Do not copy or distribute in any form without explicit permission.
###################

import requests
import subprocess
import time
import unittest
import collections


### initialize constants

hostname = 'localhost' # Windows and Mac users can change this to the docker vm ip
hostBaseUrl = 'http://{}'.format(hostname)

imageName = "asg3img"
subnetName = "asg3net"
subnetRange = "10.10.0.0/16"
containerPort = "8090"

class ReplicaConfig(collections.namedtuple('ReplicaConfig', ['name', 'addr', 'host_port'])):
    @property
    def socketAddress(self):
        return '{}:{}'.format(self.addr, containerPort)
    def __str__(self):
        return self.name

alice = ReplicaConfig(name='alice', addr='10.10.0.2', host_port=8082)
bob   = ReplicaConfig(name='bob',   addr='10.10.0.3', host_port=8083)
carol = ReplicaConfig(name='carol', addr='10.10.0.4', host_port=8084)
all_replicas = [alice, bob, carol]
viewStr = lambda replicas: ','.join(r.socketAddress for r in replicas)
viewSet = lambda replicas: set(r.socketAddress for r in replicas)

def sleep(n):
    multiplier = 1
    # Increase the multiplier if you need to during debugging, but make sure to
    # set it back to 1 and test your work before submitting.
    print('(sleeping {} seconds to stabilize)'.format(n*multiplier))
    time.sleep(n*multiplier)


### docker linux commands

def removeSubnet(required=True):
    command = ['docker', 'network', 'rm', subnetName]
    print('removeSubnet:', ' '.join(command))
    subprocess.run(command, stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL, check=required)

def createSubnet():
    command = ['docker', 'network', 'create',
            '--subnet={}'.format(subnetRange), subnetName]
    print('createSubnet:', ' '.join(command))
    subprocess.check_call(command, stdout=subprocess.DEVNULL)

def buildDockerImage():
    command = ['docker', 'build', '-t', imageName, '.']
    print('buildDockerImage:', ' '.join(command))
    subprocess.check_call(command)

def runReplica(instance, view_replicas):
    assert view_replicas, 'the view can\'t be empty because it must at least contain this replica'
    command = ['docker', 'run', '--rm', '--detach',
        '--publish={}:{}'.format(instance.host_port, containerPort),
        "--net={}".format(subnetName),
        "--ip={}".format(instance.addr),
        "--name={}".format(instance.name),
        "-e=SOCKET_ADDRESS={}:{}".format(instance.addr, containerPort),
        "-e=VIEW={}".format(viewStr(view_replicas)),
        imageName]
    print('runReplica:', ' '.join(command))
    subprocess.check_call(command)

def stopAndRemoveInstance(instance, required=True):
    command = ['docker', 'stop', instance.name]
    print('stopAndRemoveInstance:', ' '.join(command))
    subprocess.run(command, stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL, check=required)
    command = ['docker', 'remove', instance.name]
    print('stopAndRemoveInstance:', ' '.join(command))
    subprocess.run(command, stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL, check=required)

def killInstance(instance, required=True):
    '''Kill is sufficient when containers are run with `--rm`'''
    command = ['docker', 'kill', instance.name]
    print('killInstance:', ' '.join(command))
    subprocess.run(command, stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL, check=required)

def connectToNetwork(instance):
    command = ['docker', 'network', 'connect', subnetName, instance.name]
    print('connectToNetwork:', ' '.join(command))
    subprocess.check_call(command)

def disconnectFromNetwork(instance):
    command = ['docker', 'network', 'disconnect', subnetName, instance.name]
    print('disconnectFromNetwork:', ' '.join(command))
    subprocess.check_call(command)


### test suite

class TestHW3(unittest.TestCase):

    @classmethod
    def setUpClass(cls):
        print('= Cleaning up resources possibly left over from a previous run..')
        stopAndRemoveInstance(alice, required=False)
        stopAndRemoveInstance(bob,   required=False)
        stopAndRemoveInstance(carol, required=False)
        removeSubnet(required=False)
        sleep(1)
        print("= Creating resources required for this run..")
        createSubnet()

    def setUp(self):
        print("== Running replicas..")
        runReplica(alice, all_replicas)
        runReplica(bob,   all_replicas)
        runReplica(carol, all_replicas)
        sleep(3)

    def tearDown(self):
        print("== Destroying replicas..")
        killInstance(alice)
        killInstance(bob)
        killInstance(carol)

    @classmethod
    def tearDownClass(cls):
        print("= Cleaning up resources from this run..")
        removeSubnet()


    # def test_single_server_operations(self):
    #     print('> Lets test standard operations on a single server.')
    #     host_port = alice.host_port
    #     metadata = None

    #     print('> Insert test_key:test_value, should get 201')
    #     response = requests.put('http://{}:{}/kvs/{}'.format(hostname, host_port, 'test_key'),
    #             json={'value':'test_value', 'causal-metadata': metadata})
    #     self.assertIn('result', response.json(), response.status_code)
    #     self.assertEqual(response.status_code, 201)
    #     #print(response.json())
    #     self.assertIn('causal-metadata', response.json())
    #     self.assertEqual(response.json()['result'], 'created')
    #     metadata = response.json()['causal-metadata']
    #     print('Success')

    #     print('> Put test_key:overwrite, should get 200')
    #     response = requests.put('http://{}:{}/kvs/{}'.format(hostname, host_port, 'test_key'),
    #             json={'value':'overwrite', 'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 200)
    #     self.assertIn('result', response.json())
    #     self.assertIn('causal-metadata', response.json())
    #     self.assertEqual(response.json()['result'], 'replaced')
    #     metadata = response.json()['causal-metadata']
    #     print('Success')

    #     print('> PUT 400 error on keylen > 50')
    #     response = requests.put('http://{}:{}/kvs/{}'.format(hostname, host_port, '123456789012345678901234567890123456789012345678901234567890'),
    #             json={'value':'my key was wayyy too long', 'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 400)
    #     self.assertIn('error', response.json())
    #     self.assertEqual(response.json()['error'], 'Key is too long')
    #     print('Success')

    #     print('> GET test_key should return overwrite')
    #     response = requests.get('http://{}:{}/kvs/{}'.format(hostname, host_port, 'test_key'),
    #             json={'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 200)
    #     self.assertIn('result', response.json())
    #     self.assertIn('causal-metadata', response.json())
    #     self.assertEqual(response.json()['result'], 'found')
    #     self.assertEqual(response.json()['value'],'overwrite')
    #     metadata = response.json()['causal-metadata']
    #     print('Success')

    #     print('> GET unknown key should 404')
    #     response = requests.get('http://{}:{}/kvs/{}'.format(hostname, host_port, 'this_key_does_not_exist'),
    #             json={'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 404)
    #     self.assertIn('error', response.json())
    #     self.assertEqual(response.json()['error'], 'Key does not exist')
    #     print('Success')

    #     print('> Delete key unknown key 404')
    #     response = requests.delete('http://{}:{}/kvs/{}'.format(hostname, host_port, 'this_key_does_not_exist'),
    #             json={'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 404)
    #     self.assertIn('error', response.json())
    #     self.assertEqual(response.json()['error'], 'Key does not exist')
    #     print('Success')

    #     print('> Delete key test_key 200')
    #     response = requests.delete('http://{}:{}/kvs/{}'.format(hostname, host_port, 'test_key'),
    #             json={'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 200)
    #     self.assertIn('result', response.json())
    #     self.assertEqual(response.json()['result'], 'deleted')
    #     print('Success')
    
    def test_operation_replications(self):
        metadata = None
        prev_metadata = None
        print('Put key_test:my_val in alice')
        response = requests.put('http://{}:{}/kvs/{}'.format(hostname, alice.host_port, 'key_test'),
                json={'value':'my_val', 'causal-metadata': metadata})
        self.assertEqual(response.status_code, 201)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'created')
        metadata = response.json()['causal-metadata']
        prev_metadata = response.json()['causal-metadata']
        print('Success')

        sleep(1)

        print('> All replicas should now have key_test:my_val')
        for replica in all_replicas:
            response = requests.get('http://{}:{}/kvs/{}'.format(hostname, replica.host_port, 'key_test'),
                    json={'causal-metadata':metadata})
            self.assertEqual(response.status_code, 200, msg='at replica, {}'.format(replica))
            self.assertIn('result', response.json(), msg='at replica, {}'.format(replica))
            self.assertIn('causal-metadata', response.json(), msg='at replica, {}'.format(replica))
            self.assertEqual(response.json()['value'], 'my_val', msg='at replica, {}'.format(replica))
            metadata = response.json()['causal-metadata']
            print(replica.name + ' passed')
        print('All Successful')

        print('Disconnect Carol\n============')
        disconnectFromNetwork(carol)
        print('============\nSuccessfully disconnected')
        #print('Wait for 5 seconds to stabilize')
        #sleep(5)

        print('> Put key_test:another_key in alice. This should replicate to bob but not carol (since disconnected)')
        response = requests.put('http://{}:{}/kvs/{}'.format(hostname, alice.host_port, 'key_test'),
                json={'value':'another_key', 'causal-metadata': metadata})
        self.assertEqual(response.status_code, 200)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'replaced')
        metadata = response.json()['causal-metadata']

        sleep(1)
        print('Reconnect Carol\n============')
        connectToNetwork(carol)
        print('============\nDone connecting Carol')
        print('Alice & Bob should have new key, Carol should have old')

        print('> ALICE GET key_test should return another_key')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, alice.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 200)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'found')
        self.assertEqual(response.json()['value'],'another_key')
        print('Success')

        print('> BOB GET key_test should return another_key')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, bob.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 200)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'found')
        self.assertEqual(response.json()['value'],'another_key')
        print('Success')

        print('> CAROL GET with most recent metadata should 503')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, carol.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 503)
        print('Success')

        print('> CAROL GET with old metadata, key_test should return my_val (Since it couldnt be replicated to)')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, carol.host_port, 'key_test'),
                json={'causal-metadata': prev_metadata})
        self.assertEqual(response.status_code, 200)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'found')
        self.assertEqual(response.json()['value'],'my_val')
        print('Success')

        print('> Killing CAROL and restarting should make the key available\n============')
        killInstance(carol)
        print('> Now launch carol again')
        runReplica(carol,all_replicas)
        print('> Give 5 seconds to stabilize\n============\n')
        sleep(5)
        print('> Now carol should have the most up to date value, another_key')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, carol.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 200)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'found')
        self.assertEqual(response.json()['value'],'another_key')
        metadata = response.json()['causal-metadata']
        print('Success')
        print('> Now lets check our 503 errors if given outdated metadata')
        
        # What is the actuall condition here?
        # print('> GETS should 503 with outdated metadata')
        # for replica in all_replicas:
        #     response = requests.get('http://{}:{}/kvs/{}'.format(hostname, replica.host_port, 'key_test'),
        #             json={'causal-metadata':prev_metadata})
        #     self.assertEqual(response.status_code, 503)
        #     print(replica.name + ' passed')
        # print('All GETs passed')
        
        print('> PUTS should 503 with outdated metadata')
        for replica in all_replicas:
            response = requests.put('http://{}:{}/kvs/{}'.format(hostname, replica.host_port, 'key_test'),
                    json={'causal-metadata':prev_metadata})
            self.assertEqual(response.status_code, 503)
            print(replica.name + ' passed')
        print('All PUTs passed')
        
        print('> DELETEs should 503 with outdated metadata')
        for replica in all_replicas:
            response = requests.delete('http://{}:{}/kvs/{}'.format(hostname, replica.host_port, 'key_test'),
                    json={'causal-metadata':prev_metadata})
            self.assertEqual(response.status_code, 503)
            print(replica.name + ' passed')
        print('All DELETEs passed')
        
        print('> Now delete key_test on carol')
        response = requests.delete('http://{}:{}/kvs/{}'.format(hostname, carol.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 200)
        self.assertIn('result', response.json())
        self.assertIn('causal-metadata', response.json())
        self.assertEqual(response.json()['result'], 'deleted')
        metadata = response.json()['causal-metadata']
        print('Success')
        sleep(1)
        print('> Alice, Carol, and Bob should all 404 key_test')
        print('> Alice...')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, alice.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 404)
        self.assertIn('error', response.json())
        self.assertEqual(response.json()['error'], 'Key does not exist')
        print('Success')
        print('> Bob...')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, bob.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 404)
        self.assertIn('error', response.json())
        self.assertEqual(response.json()['error'], 'Key does not exist')
        print('Success')
        print('> Carol...')
        response = requests.get('http://{}:{}/kvs/{}'.format(hostname, carol.host_port, 'key_test'),
                json={'causal-metadata': metadata})
        self.assertEqual(response.status_code, 404)
        self.assertIn('error', response.json())
        self.assertEqual(response.json()['error'], 'Key does not exist')
        print('Success')
    
    # def test_view_replication(self):
    #     #removed = carol
    #     #remaining_replicas = sorted(set(all_replicas) - set([removed]))
    #     #return = viewSet()
    #     # THIS TEST IS UNCHANGED FROM THE GIVEN TESTS
    #     metadata = None

    #     print('=== Check replica views')
    #     for replica in all_replicas:
    #         response = requests.get('http://{}:{}/view'.format(hostname, replica.host_port))
    #         self.assertEqual(response.status_code, 200, msg='at replica, {}'.format(replica))
    #         self.assertIn('view', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertEqual(set(response.json()['view']), viewSet(all_replicas), msg='at replica, {}'.format(replica))


    #     print('>>> Put apple:strudel into the store')

    #     response = requests.put('http://{}:{}/kvs/{}'.format(hostname, alice.host_port, 'apple'),
    #             json={'value':'strudel', 'causal-metadata': metadata})
    #     self.assertEqual(response.status_code, 201)
    #     self.assertIn('result', response.json())
    #     self.assertIn('causal-metadata', response.json())
    #     self.assertEqual(response.json()['result'], 'created')
    #     metadata = response.json()['causal-metadata']

    #     print('... Wait for replication')
    #     sleep(5)

    #     print('=== Check apple:strudel at replicas {}'.format(','.join(r.name for r in all_replicas)))
    #     for replica in all_replicas:
    #         response = requests.get('http://{}:{}/kvs/{}'.format(hostname, replica.host_port, 'apple'),
    #                 json={'causal-metadata':metadata})
    #         self.assertEqual(response.status_code, 200, msg='at replica, {}'.format(replica))
    #         self.assertIn('result', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertIn('causal-metadata', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertEqual(response.json()['value'], 'strudel', msg='at replica, {}'.format(replica))
    #         metadata = response.json()['causal-metadata']


    #     removed = carol
    #     remaining_replicas = sorted(set(all_replicas) - set([removed]))
    #     print('>>> Remove replica {}'.format(removed))
    #     killInstance(removed)

    #     print('... Wait for stabilization')
    #     sleep(5)


    #     print('>>> Put chocolate:eclair into the store')
    #     response = requests.put('http://{}:{}/kvs/{}'.format(hostname, alice.host_port, 'chocolate'),
    #             json={'value':'eclair', 'causal-metadata':metadata})
    #     self.assertEqual(response.status_code, 201)
    #     self.assertIn('result', response.json())
    #     self.assertIn('causal-metadata', response.json())
    #     self.assertEqual(response.json()['result'], 'created')
    #     metadata = response.json()['causal-metadata']

    #     print('... Wait for replication, down detection, and stabilization')
    #     sleep(5)

    #     print('=== Check chocolate:eclair at replicas {}'.format(','.join(r.name for r in remaining_replicas)))
    #     for replica in remaining_replicas:
    #         response = requests.get('http://{}:{}/kvs/{}'.format(hostname, replica.host_port, 'chocolate'),
    #             json={'causal-metadata':metadata})
    #         self.assertEqual(response.status_code, 200, msg='at replica, {}'.format(replica))
    #         self.assertIn('result', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertIn('causal-metadata', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertEqual(response.json()['value'], 'eclair', msg='at replica, {}'.format(replica))
    #         metadata = response.json()['causal-metadata']

    #     print('=== Check view at replicas {}'.format(','.join(r.name for r in remaining_replicas)))
    #     for replica in remaining_replicas:
    #         response = requests.get('http://{}:{}/view'.format(hostname, replica.host_port))
    #         self.assertEqual(response.status_code, 200, msg='at replica, {}'.format(replica))
    #         self.assertIn('view', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertEqual(set(response.json()['view']), viewSet(remaining_replicas), msg='at replica, {}'.format(replica))


    #     print('>>> Start replica {}'.format(removed))
    #     runReplica(removed, all_replicas)

    #     print('... Wait for stabilization')
    #     sleep(5)


    #     print('=== Check replica views')
    #     for replica in all_replicas:
    #         response = requests.get('http://{}:{}/view'.format(hostname, replica.host_port))
    #         self.assertEqual(response.status_code, 200, msg='at replica, {}'.format(replica))
    #         self.assertIn('view', response.json(), msg='at replica, {}'.format(replica))
    #         self.assertEqual(set(response.json()['view']), viewSet(all_replicas), msg='at replica, {}'.format(replica))


    #     print('=== Check keys apple,chocolate at replicas {}'.format(removed.name))
    #     for key,val in {'apple':'strudel', 'chocolate':'eclair'}.items():
    #         response = requests.get('http://{}:{}/kvs/{}'.format(hostname, removed.host_port, key),
    #             json={'causal-metadata':metadata})
    #         self.assertEqual(response.status_code, 200, msg='for key, {}'.format(key))
    #         self.assertIn('result', response.json(), msg='for key, {}'.format(key))
    #         self.assertIn('causal-metadata', response.json(), msg='for key, {}'.format(key))
    #         self.assertEqual(response.json()['value'], val, msg='for key, {}'.format(key))
    #         metadata = response.json()['causal-metadata']
        




if __name__ == '__main__':
    try:
        buildDockerImage()
    except KeyboardInterrupt:
        TestHW3.setUpClass()
        TestHW3.tearDownClass()
    unittest.main(verbosity=0)
