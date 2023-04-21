#!/usr/bin/env python3
import argparse
import os.path
from ldappy import load_config, connect


def get_users(conn, search_base):

    attributes = [
        "cn", "givenName", "sn", 
        "uid", "uidNumber", "gidNumber", "userPassword",
        "mail",
        "homeDirectory",
        "loginShell"
    ]

    entries = conn.extend.standard.paged_search(search_base, '(objectClass=posixAccount)', attributes=attributes, paged_size=10)
    for entry in entries:
        print(f"dn: {entry['dn']}")
        for k in attributes:
            vs = entry['attributes'][k]
            if not isinstance(vs, list):
                vs = [vs]
            for v in vs:
                print(f"  {k}: {v}")
    

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('-c', '--config', help='path to config file', type=str, default=None)
    
    args = parser.parse_args()
    
    if args.config is not None and not os.path.exists(args.config):
        raise ValueError("configuration file not found")

    config = load_config(args.config)
    
    _, conn = connect(config['ldap_uri'], config['ldap_bind_dn'], config['ldap_bind_pw'],
                    start_tls=config['ldap_start_tls'], 
                    ca_cert_file=config['ldap_ca_cert_file']
                    )
    root_dn = config['ldap_root_dn']
    
    users = get_users(conn, root_dn)


if __name__ == "__main__":
    main()
