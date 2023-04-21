import os.path
from ldap3 import Server, Connection, ALL
from ldap3 import Tls
import ssl


def connect(uri, bind_dn, bind_pw, *, start_tls=False, ca_cert_file=None):

    if start_tls:
        if ca_cert_file is None:
            raise ValueError("require a ca certificate file to start tls")
        ca_cert_file = os.path.expanduser(ca_cert_file)
        if not os.path.exists(ca_cert_file):
            raise ValueError(f"certificate file not found: {ca_cert_file}")
        
        tls_configuration = Tls(
            validate=ssl.CERT_REQUIRED, 
            version=ssl.PROTOCOL_TLSv1_2,
            ca_certs_file=ca_cert_file
        )
        server = Server(
            uri, 
            tls=tls_configuration, 
            get_info=ALL
        )
    
    else:
        server = Server(
            uri, 
            get_info=ALL
        )
    
    conn = Connection(
        server, 
        user=bind_dn, password=bind_pw
    )
    if start_tls:
        conn.start_tls()
    conn.bind()
    
    return server, conn

