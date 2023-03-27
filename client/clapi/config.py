import os.path
from io import StringIO
import importlib.resources as pkg_resources
import yaml


def load_config():
    
    # the path to the package
    clapi = pkg_resources.files('clapi')
    
    # load the config file from the package
    config = os.path.join(clapi, 'config.yaml')
    with open(config, "r") as f:
        config = yaml.safe_load(f)
        
    # find the absolute path to tls files
    
    
    path = os.path.expanduser(config['ca_cert_file'])
    if not os.path.isabs(path):
        path = os.path.join(clapi, path)
    config['ca_cert_file'] = path

    path = os.path.expanduser(config['client_cert_file'])
    if not os.path.isabs(path):
        path = os.path.join(clapi, path)
    config['client_cert_file'] = path

    path = os.path.expanduser(config['client_key_file'])
    if not os.path.isabs(path):
        path = os.path.join(clapi, path)
    config['client_key_file'] = path
    
    return config
    

