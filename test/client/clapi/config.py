import os.path
from io import StringIO
import importlib.resources as pkg_resources
import yaml


def load_config(cfg_file):
    
    cfg_dir = os.path.dirname(cfg_file)
    
    # load the config file from the package
    with open(cfg_file, "r") as f:
        config = yaml.safe_load(f)
        
    # find the absolute path to tls files
    path = os.path.expanduser(config['ca_cert_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['ca_cert_file'] = path

    path = os.path.expanduser(config['admin_cert_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['admin_cert_file'] = path

    path = os.path.expanduser(config['operator_key_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['operator_key_file'] = path
    
    path = os.path.expanduser(config['operator_cert_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['operator_cert_file'] = path

    path = os.path.expanduser(config['admin_key_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['admin_key_file'] = path
    
    path = os.path.expanduser(config['user_cert_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['user_cert_file'] = path

    path = os.path.expanduser(config['user_key_file'])
    if not os.path.isabs(path):
        path = os.path.join(cfg_dir, path)
    config['user_key_file'] = path

    return config
    

