#!/usr/bin/env python3
import argparse
import os
import random

from ldappy import load_config, connect
from ldappy import create_group, check_gid_number


def add_group(conn, domain_dn, domain_dns, group_name):
    
    # find a gid_number
    gid_number = None
    while gid_number is None:
        gid_number = random.randint(5000, 50000)
        
        if check_gid_number(conn, domain_dn, gid_number) == False:
            gid_number = None
        
    # greate the group
    create_group(conn, domain_dn, group_name, gid_number)
    
    return gid_number
    


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config file", type=str, default=None)
    parser.add_argument("group_name", help="group name", type=str)
    
    args = parser.parse_args()
    
    if args.config is not None and not os.path.exists(args.config):
        raise ValueError("configuration file not found")

    config = load_config(args.config)
    
    _, conn = connect(config['ldap_uri'], config['ldap_bind_dn'], config['ldap_bind_pw'],
                    start_tls=config['ldap_start_tls'], 
                    ca_cert_file=config['ldap_ca_cert_file']
                    )
    root_dn  = config['ldap_root_dn']
    root_dns = config['ldap_root_dns']
    
    gid_number = add_group(conn, root_dn, root_dns, args.group_name)
    print(f"Created group:")
    print(f"       name: {args.group_name}")
    print(f"  id number: {gid_number}")

if __name__ == "__main__":
    main()

