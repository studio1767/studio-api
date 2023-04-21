#!/usr/bin/env python3
import argparse
import logging
import time

import grpc
import v1.project_pb2 as project_pb2
import v1.project_pb2_grpc as project_pb2_grpc

import clapi


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('config_file', help='configuration file to load', type=str, default=None)
    args = parser.parse_args()

    # load the config file
    config = clapi.load_config(args.config_file)
    
    # the server address and port to connect to
    server = f"{config['api_server']}:{config['api_port']}"
    print(f"server: {server}")

    # read the certificate files
    with open(config['ca_cert_file'], 'rb') as f:
        ca_cert = f.read()
    with open(config['user_cert_file'], 'rb') as f:
        user_cert = f.read()
    with open(config['user_key_file'], 'rb') as f:
        user_key = f.read()

    credentials = grpc.ssl_channel_credentials(ca_cert, user_key, user_cert)
    with grpc.secure_channel(server, credentials) as channel:
        stub = project_pb2_grpc.StudioStub(channel)

        response = stub.Hello(project_pb2.HelloRequest(name='you'))
        print(response.message)
        
        stamp = int(time.time())
        response = stub.CreateProject(project_pb2.ProjectRequest(name=f"Project-{stamp}", code=f"p{stamp}"))
        print(f"created project with id {response.id}")

        responses = stub.Projects(project_pb2.ProjectFilter(regex="*"))
        for response in responses:
            print(response.id, response.name, response.code)
        
        

if __name__ == '__main__':
    logging.basicConfig()
    main()

