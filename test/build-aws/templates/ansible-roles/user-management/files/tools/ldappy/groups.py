import io, re
import random
from collections import defaultdict
from jinja2 import PackageLoader, Environment


_attribute_re = re.compile(f'([^:]+):\s+(\S+)\s*')


def check_gid_number(conn, search_base, gid_number):
    entries = conn.extend.standard.paged_search(
        search_base,
        f'(&(objectClass=posixGroup)(gidNumber={gid_number}))',
        attributes=['cn', 'gidNumber'],
        paged_size=1
    )
    for entry in entries:
        return False
    return True


def create_group(conn, domain_dn, group_name, gid_number):

    # check if the group already exists
    entries = conn.extend.standard.paged_search(
        domain_dn,
        f'(&(objectClass=posixGroup)(cn={group_name}))',
        attributes=['cn', 'gidNumber'],
        paged_size=1
    )
    for entry in entries:
        raise ValueError(f"group {group_name} already exists")

    # check if the gid already exists
    entries = conn.extend.standard.paged_search(
        domain_dn,
        f'(&(objectClass=posixGroup)(gidNumber={gid_number}))',
        attributes=['cn', 'gidNumber'],
        paged_size=1
    )
    for entry in entries:
        raise ValueError(f"group {gid_number} already exists")

    # expand the group ldif template
    env = Environment(loader=PackageLoader("ldappy", "templates"))
    template = env.get_template("group.ldif")
    doc = template.render(
        domain_dn=domain_dn,
        group_name=group_name,
        gid_number=gid_number
    )
    
    # parse the ldif doc
    dn = None
    oclass = []
    attrs = defaultdict(list)
    
    for line in io.StringIO(doc):
        mo = _attribute_re.match(line)
        if mo is None:
            raise ValueError(f"couldn't match {line}")

        key, value = mo.group(1).strip().lower(), mo.group(2).strip()

        if key == "dn":
            dn = value
            continue
        if key == "objectclass":
            oclass.append(value)
            continue
        attrs[key].append(value)
        
    rv = conn.add(dn, oclass, attrs)
    if rv == False:
        raise ValueError(f"failed to create group: {conn.result}")
