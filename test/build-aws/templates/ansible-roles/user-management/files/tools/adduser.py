#!/usr/bin/env python3
import argparse
import os
import getpass
import string
import secrets
import random, hashlib
from base64 import b64encode

from ldappy import load_config, connect
from ldappy import create_group, create_user
from ldappy import check_gid_number, check_uid_number


pw_chars = string.ascii_letters + string.digits + string.punctuation
def generate_pw():
    user_pw = ""
    for _ in range(16):
        user_pw += secrets.choice(pw_chars)
    return user_pw


def encode_pw(pw):
    enc_pw = pw.encode('utf-8')
    salt = os.urandom(16)
    hash_pw = hashlib.sha1(enc_pw + salt).digest()
    
    encoded = "{SSHA}" + b64encode(hash_pw + salt).decode("utf-8")

    return encoded


def add_user(conn, domain_dn, domain_dns, user_name, user_password, given_name, family_name):
    
    # find an id_number that can be used for uid and gid
    id_number = None
    while id_number is None:
        id_number = random.randint(5000, 50000)
        
        if check_gid_number(conn, domain_dn, id_number) == False or check_uid_number(conn, domain_dn, id_number) == False:
            id_number = None
        
    # greate the group and user
    create_group(conn, domain_dn, user_name, id_number)
    create_user(conn, domain_dn, domain_dns,
                        user_name, user_password, 
                        given_name, family_name, 
                        id_number, id_number)
    
    return id_number


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--config", help="path to config file", type=str, default=None)
    parser.add_argument("-r", help="generate random password", action="store_true")
    parser.add_argument("user_name", help="user name", type=str)
    parser.add_argument("given_name", help="given name", type=str)
    parser.add_argument("family_name", help="family name", type=str)
    
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
    
    # get the user password
    if args.r:
        user_pw = generate_pw()
        print(f"Generated Password: {user_pw}")
    else:
        user_pw = getpass.getpass()
        pw2 = getpass.getpass("Repeat Password: ")
        if user_pw != pw2:
            print(f"passwords differ")
            sys.exit(1)
    
    encoded_pw = encode_pw(user_pw)
    
    id_number = add_user(conn, root_dn, root_dns,
                    args.user_name, encoded_pw, args.given_name, args.family_name)

    print(f"Created user and group:")
    print(f"       name: {args.user_name}")
    print(f"  id number: {id_number}")
    if args.r:
        print(f"   password: {user_pw}")


if __name__ == "__main__":
    main()

# given_name
# family_name
# user_name
# uid
# gid
# password

