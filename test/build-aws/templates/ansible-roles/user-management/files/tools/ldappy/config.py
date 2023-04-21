import os.path
import re
from io import StringIO
import importlib.resources as pkg_resources
import yaml


root_dn_re = re.compile("dc=(.*)")

def load_config(path=None):

    # the path to the package
    ldappy = pkg_resources.files('ldappy')

    # load the default config file from the package
    config = os.path.join(ldappy, 'config.yaml')
    with open(config, 'r') as f:
        config = yaml.safe_load(f)

    # look for override config files
    cfg_files = ["~/.config/ldappy.yaml"]
    if path is not None:
        cfg_files.append(path)
    
    for cfg_file in cfg_files:
        cfg_file = os.path.expanduser(cfg_file)
        if not os.path.exists(cfg_file):
            continue
        with open(cfg_file) as fd:
            next_config = yaml.safe_load(fd)
    
        config = {**config, **next_config}
    
    # find the absolute path to ca certificate file
    cert_path = os.path.expanduser(config['ldap_ca_cert_file'])
    if not os.path.isabs(cert_path):
        cert_path = os.path.join(ldappy, cert_path)
    config['ldap_ca_cert_file'] = cert_path
        
    # calculate the root dns
    if config.get('ldap_root_dns', None) is None:
        root_dn = config['ldap_root_dn']
        root_dcs = root_dn.split(',')
        
        dcs = []
        for root_dc in root_dcs:
            mo = root_dn_re.match(root_dc)
            if mo is None:
                raise ValueError(f"couldn't match {root_dc}")
            dcs.append(mo.group(1))
        
        root_dns = ".".join(dcs)
        config['ldap_root_dns'] = root_dns
            
    return config
    
