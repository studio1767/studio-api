import os.path
import requests
from gql import Client
from gql.transport.requests import RequestsHTTPTransport


def new_client(cfg):

    api_server, api_port = cfg['api_server'], cfg['api_port']
    ca_cert_file = cfg['ca_cert_file']
    client_cert_file = cfg['client_cert_file']
    client_key_file = cfg['client_key_file']
    
    if not os.path.exists(ca_cert_file):
        raise ValueError(f"certificate file not found: {ca_cert_file}")
    if not os.path.exists(client_cert_file):
        raise ValueError(f"client certificate file not found: {client_cert_file}")
    if not os.path.exists(client_key_file):
        raise ValueError(f"client key file not found: {client_key_file}")
    
    # authenticate with the api... and get the cookie
    with requests.Session() as s:
        r = s.get(f"https://{api_server}:{api_port}/hello", verify=ca_cert_file, cert=(client_cert_file, client_key_file))
        if r.status_code == 200:
            print(f"auth passed: {r.status_code} {r.text.strip()}")
        else:
            print(f"auth failed: {r.status_code}")
    
        print(type(s.cookies))
        for k, v in s.cookies.items():
            print(f"{k}: {type(v)}")
    
        transport = RequestsHTTPTransport(
            url=f"https://{api_server}:{api_port}/",
            verify=ca_cert_file,
            cookies=s.cookies,
            retries=3,
        )

    # finally, create the gql client
    client = Client(transport=transport, fetch_schema_from_transport=True)
    return client

