#!/usr/bin/env python3
import argparse
import logging
import time

import grpc
import v1.project_pb2 as project_pb2
import v1.project_pb2_grpc as project_pb2_grpc

import clapi



def run(server, ca_cert, user_key, user_cert, expected):
    
    credentials = grpc.ssl_channel_credentials(ca_cert, user_key, user_cert)
    with grpc.secure_channel(server, credentials) as channel:
        stub = project_pb2_grpc.StudioStub(channel)

        try:
            print("pinging server:")
            response = stub.Ping(project_pb2.PingRequest(name='world'))
        except:
            if expected[0] == False:
                print("-> PASS: ping failed")
            else:
                print("-> FAIL: ping failed")
        else:
            print(f"-> PASS: {response.message}")
        
        try:
            print("create project:")
            stamp = int(time.time())
            response = stub.CreateProject(project_pb2.ProjectRequest(name=f"project-{stamp}", code=f"p{stamp}"))
        except:
            if expected[1] == False:
                print("-> PASS: create failed")
            else:
                print("-> FAIL: create failed")
        else:
            print(f"-> PASS: {response.id}")

        try:
            print("projects:")
            responses = stub.Projects(project_pb2.ProjectFilter(regex="*"))
        except:
            if expected[2] == False:
                print("-> PASS: projects failed")
            else:
                print("-> FAIL: projects failed")
        else:
            for response in responses:
                print(f"-> PASS: {response.id} {response.name} {response.code}")
    

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('config_file', help='configuration file to load', type=str, default=None)
    args = parser.parse_args()

    # load the config file
    config = clapi.load_config(args.config_file)
    
    # the server address and port to connect to
    server = f"{config['api_server']}:{config['api_port']}"
    print(f"server: {server}")

    # read the ca certificate
    with open(config['ca_cert_file'], 'rb') as f:
        ca_cert = f.read()
        
    # run admin tests
    print("*** admin tests ***")
    with open(config['admin_key_file'], 'rb') as f:
        key = f.read()
    with open(config['admin_cert_file'], 'rb') as f:
        cert = f.read()
    run(server, ca_cert, key, cert, (True, True, True))

    # run operator tests
    print("*** operator tests ***")
    with open(config['operator_key_file'], 'rb') as f:
        key = f.read()
    with open(config['operator_cert_file'], 'rb') as f:
        cert = f.read()
    run(server, ca_cert, key, cert, (True, True, True))

    # run user tests
    print("*** user tests ***")
    with open(config['user_key_file'], 'rb') as f:
        key = f.read()
    with open(config['user_cert_file'], 'rb') as f:
        cert = f.read()
    run(server, ca_cert, key, cert, (True, False, True))




if __name__ == '__main__':
    logging.basicConfig()
    main()

